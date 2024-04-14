package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
    "crypto/sha1"
    "encoding/hex"
    "fmt"
    "sort"
    "strings"
)

// 微信开发者接口校验
func WxVerify(c *gin.Context){
	// 填写开发者填写的token
	token := "htq"

	// 接收请求参数
	signature := c.Query("signature")
	timestamp := c.Query("timestamp")
	nonce := c.Query("nonce")
	echostr := c.Query("echostr")

	// 校验signature
	if WXCheckSignature(signature, timestamp, nonce, token) {
		fmt.Println("微信公众号接入校验成功！")
		c.String(http.StatusOK, echostr)
	} else {
		fmt.Println("微信公众号接入校验失败！")
		c.String(http.StatusOK, "校验失败")
	}
}

// 微信用户交互校验
func WXCheckSignature(signature, timestamp, nonce, token string) bool {

    params := []string{timestamp, nonce, token}
    sort.Strings(params)

    // 拼接参数字符串
    paramStr := strings.Join(params, "")

    // 进行sha1加密
    sha1Str := sha1.Sum([]byte(paramStr))

    // 与signature进行比较
    signatureStr := hex.EncodeToString(sha1Str[:])
    return signatureStr == signature
}