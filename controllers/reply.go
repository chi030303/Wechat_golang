package controllers

import (
	"log"
	"io"
	"fmt"
	"time"
	"encoding/xml"
	"github.com/gin-gonic/gin"
	"strings"
	"Wechat-project/config"
	"regexp"
	"strconv"
	"gorm.io/gorm"
)

// 微信消息结构体
type XMLData struct {
	ToUserName   string `xml:"ToUserName"`
	FromUserName string `xml:"FromUserName"`
	CreateTime   int64  `xml:"CreateTime"`
	MsgType      string `xml:"MsgType"`
	Content      string `xml:"Content"`
}

// 解析 XML 函数
func ParseXML(body []byte) (XMLData, error) {
    var xmlData XMLData
    err := xml.Unmarshal(body, &xmlData)
    if err != nil {
        return XMLData{}, err
    }
    return xmlData, nil
}

// 将回复消息转成微信需要的形式
func (r XMLData) send() string {
	XmlForm := `
            <xml>
                <ToUserName><![CDATA[%s]]></ToUserName>
                <FromUserName><![CDATA[%s]]></FromUserName>
                <CreateTime>%d</CreateTime>
                <MsgType><![CDATA[%s]]></MsgType>
                <Content><![CDATA[%s]]></Content>
            </xml>
            `
	return fmt.Sprintf(XmlForm, r.ToUserName, r.FromUserName, r.CreateTime, r.MsgType,r.Content)
}

// 回复查询
type ReplyQuery struct {
	Signature    string `form:"signature" binding:"required"`
	Timestamp    string `form:"timestamp" binding:"required"`
	Nonce        string `form:"nonce" binding:"required"`
	Openid       string `form:"openid" binding:"required"`
	EncryptType  string `form:"encrypt_type" binding:"required"`
	MsgSignature string `form:"msg_signature" binding:"required"`
}

// 用户结构体
type User struct {
	UserID   uint `gorm:"primaryKey;autoIncrement"`
	UserName string	`gorm:"user_name"`
	Access   int `gorm:"access"`
}

func (User) TableName() string {
	return "User"
}

// CreateUserIfNeeded 创建新用户，如果用户不存在的话
func CreateUser(userName string, db *gorm.DB) error {
    // 检查数据库中是否存在该用户
    var existingUser User
    if err := db.Where("user_name = ?", userName).First(&existingUser).Error; err != nil {
        if err == gorm.ErrRecordNotFound { // 如果用户不存在
            // 创建新用户
            newUser := User{
                UserName: userName, // 假设将 FromUserName 作为用户名
                Access:   0,        // 默认为0
            }
            if err := db.Create(&newUser).Error; err != nil {
                return err
            }
        } else {
            return err
        }
    }
    return nil
}

// 回复检测
func Reply(c *gin.Context){
	replyQuery := ReplyQuery{}
	c.ShouldBindQuery(&replyQuery)
	token := "htq"
	log.Printf("Openid:%s\n",replyQuery.Openid)
	if !WXCheckSignature(replyQuery.Signature, replyQuery.Timestamp, replyQuery.Nonce,token) {
		c.String(403, "Invalid signature")
		return
	}
	xmlData, _ := io.ReadAll(c.Request.Body)
	log.Printf("Receive: %s\n", c.Request.Body)
	recXml, err := ParseXML(xmlData)
	log.Printf("FromUserName:%s\n",recXml.FromUserName)
	
	db , _ := config.InitDB()
	// 新建用户，以便后续的收藏管理
	if err := CreateUser(recXml.FromUserName,db); err != nil {
        c.String(500, "Failed to create new user")
        return
    }

	if err != nil {
		c.String(500, "Failed to parse XML")
		return  // 返回一个空的 XMLData 对象
	}
	replyMsg := messageHandle(recXml) 
	c.String(200,replyMsg.send())
}

// 公众号介绍消息
func introductionMessage() string {
    // 构造介绍消息
	content := `您好，这里是小贺！欢迎来到我的公众号。在这里您可以查看[收录企业]，搜索[招聘信息]，搜索[薪资信息]，[收藏、取消收藏企业]。
	例如：1.输入"收录企业"，可以看到我已收录的企业，但是由于微信回复消息长度限制，我做了分页处理，您可以输入页数来进行跳转，您可以输入1-68，比如“收录企业 3”，如果不输入则默认页数为第一页。
	2.[招聘信息]，输入“招聘 腾讯（任一公司名都可）”，若存在该公司，则展示相关信息。
	3.[薪资信息]，输入“薪资 腾讯（任一公司名都可）”，若存在该公司，则展示相关信息。
	4.[收藏、取消收藏企业]，输入“收藏/取消收藏 字节跳动（任一公司名都可），可以将对应的企业在收藏列表中加入或删除。然后，您也可以通过输入“查看收藏”来查看已收藏的企业。`

    return content
}

// 回复处理函数
func messageHandle(recXML XMLData) XMLData {

	content := ""
    
	replyMsg := XMLData{
		ToUserName:   recXML.FromUserName,
		FromUserName: recXML.ToUserName,
		CreateTime:   time.Now().Unix(),
		MsgType:      "text",
		Content:      content,
	}

    switch recXML.MsgType {
		case "text":
			// 进一步解析text
			replyMsg.Content = parseText(recXML.FromUserName,recXML.Content)
			return replyMsg
		default:
			// 提示必须发送text类型消息
			content = "您必须发送text（文字）类型消息！"
			return replyMsg
    }
    
    // 发送回复消息
    fmt.Println("回复消息:", content)
	return replyMsg
}

// 检查用户输入是否包含数字
func containsNumber(input string) bool {
	// 定义包含数字的字符集合
	numbers := "0123456789"

	// 遍历输入字符串的每个字符
	for _, char := range input {
		// 将字符转换为字符串，并在数字集合中查找
		if strings.ContainsAny(string(char), numbers) {
			return true // 如果找到数字，返回true
		}
	}

	return false // 如果未找到数字，返回false
}

// 提取回复中的数字
func extractNumber(input string) int {
	// 定义匹配数字的正则表达式
	re := regexp.MustCompile("[0-9]+")

	// 在输入字符串中查找所有匹配的数字字符串
	matches := re.FindAllString(input, -1)

	// 将匹配到的数字字符串连接成一个整体字符串
	combined := strings.Join(matches, "")

	// 将整体字符串转换为整数类型
	num, err := strconv.Atoi(combined)
	if err != nil {
		return 0
	}

	return num
}

// 解析用户回复
func parseText(userName,rcvMsg string) string{

	// 初始化数据库连接
	db , err := config.InitDB()
	if err != nil {
		panic("failed to connect database")
	}

	// 创建控制器实例并注入数据库连接
	offerController := NewOfferController(db)
	salaryController := NewSalaryController(db)

    // 解析text消息的逻辑
    if strings.Contains(rcvMsg, "收录企业") {
		page := 1
		if containsNumber(rcvMsg){
			page  = extractNumber(rcvMsg)
		}
        return offerController.GetCompaniesList(page)
    } else if strings.Contains(rcvMsg, "招聘") {
		cleanMsg := strings.TrimSpace(strings.Replace(rcvMsg, "招聘", "", 1))
        return offerController.GetOfferInfo(cleanMsg)
    } else if strings.Contains(rcvMsg, "薪资") {
		cleanMsg := strings.TrimSpace(strings.Replace(rcvMsg, "薪资", "", 1))
        return salaryController.GetSalaryInfo(cleanMsg)
    } else if strings.Contains(rcvMsg, "取消收藏") {
		cleanMsg := strings.TrimSpace(strings.Replace(rcvMsg, "取消收藏", "", 1))
        return offerController.UnCollectCompany(userName,cleanMsg)
    } else if strings.Contains(rcvMsg, "查看收藏") {
		reply , _ := offerController.GetFavoriteCompanies(userName)
        return reply
    } else if strings.Contains(rcvMsg, "收藏") {
		cleanMsg := strings.TrimSpace(strings.Replace(rcvMsg, "收藏", "", 1))
        return offerController.CollectCompany(userName,cleanMsg)
    } else {
        return introductionMessage()
    }
}
