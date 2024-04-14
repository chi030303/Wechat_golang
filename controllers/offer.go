package controllers

import (
	"Wechat-project/utils"
	"gorm.io/gorm"
	"strings"
    "errors"
)

// 靠结构体实现重载
type OfferController struct {
    DB *gorm.DB
}

func NewOfferController(db *gorm.DB) *OfferController {
    return &OfferController{DB: db}
}

// 查询具体的招聘信息
func (oc* OfferController) GetOfferInfo(companyName string) string{
	var company utils.RecruitmentCompanies
    var result strings.Builder

    // 查询数据库，获取符合条件的公司信息
    if err := oc.DB.Where("company_name LIKE ?", "%"+companyName+"%").First(&company).Error; err != nil {
        // 处理查询错误
        if err == gorm.ErrRecordNotFound {
            return "未找到名称包含 \"" + companyName + "\" 的公司信息"
        }
        // 其他错误处理
        return "查询公司信息时出现错误"
    }

    // 组织返回信息
    result.WriteString("公司名称：" + company.CompanyName + "\n")
    if company.RecruitmentTarget != "" {
        result.WriteString("招聘目标：" + company.RecruitmentTarget + "\n")
    }
    if company.JobPosition != "" {
        result.WriteString("职位：" + company.JobPosition + "\n")
    }
    if company.RecruitmentMajors != "" {
        result.WriteString("招聘专业：" + company.RecruitmentMajors + "\n")
    }
    if company.WorkLocation != "" {
        result.WriteString("工作地点：" + company.WorkLocation + "\n")
    }
    if company.RecruitmentSession != "" {
        result.WriteString("招聘届数：" + company.RecruitmentSession + "\n")
    }
    if company.DeliveryTime != "" {
        result.WriteString("投递时间：" + company.DeliveryTime + "\n")
    }

    return result.String()
}

// 查询企业列表
func (oc *OfferController) GetCompaniesList(page int) string {
    var companies []string

	pageSize := 15

	// 定义最大页数
	maxPage := 68

	// 如果页数超过最大页数，则将页数设置为最大页数
	if page > maxPage {
		page = maxPage
	}
		
    // 计算偏移量
    offset := (page - 1) * pageSize

    // 查询数据库，获取指定页数的公司名称
    if err := oc.DB.Model(&utils.RecruitmentCompanies{}).Offset(offset).Limit(pageSize).Pluck("company_name", &companies).Error; err != nil {
        // 处理错误
        return ""
    }

    // 将公司名称拼接为字符串
    content := strings.Join(companies[:], "\n")

    return content
}

// FavoriteCompany 表示收藏公司模型
type FavoriteCompany struct {
    FavoriteID uint `gorm:"primaryKey;autoIncrement"`
    UserID     uint `gorm:"user_id"`
    CompanyID  uint `gorm:"company_id"`
}

func (FavoriteCompany) TableName() string {
	return "Favorite_Company"
}

// GetUserByName 根据用户名字查找用户ID
func (oc *OfferController) GetUserByName(userName string) (uint, error) {
    var user User
    if err := oc.DB.Where("user_name = ?", userName).First(&user).Error; err != nil {
        return 0, err
    }
    return user.UserID, nil
}

// CollectCompany 收藏公司
func (oc *OfferController) CollectCompany(userName, companyName string) string {
    // 根据用户名字查找用户ID
    userID, err := oc.GetUserByName(userName)
    if err != nil {
        return "收藏失败：" + err.Error()
    }

    // 查询招聘公司表，根据公司名称找到对应的公司ID
    var company utils.RecruitmentCompanies
    if err := oc.DB.Where("company_name LIKE ?", "%"+companyName+"%").First(&company).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return "收藏失败：公司名不存在"
        }
        return "收藏失败：" + err.Error()
    }

    // 查询是否已经收藏
    var existingFavorite FavoriteCompany
    if err := oc.DB.Where("user_id = ? AND company_id = ?", userID, company.CompanyID).First(&existingFavorite).Error; err == nil {
        return "收藏失败：已经添加"
    }

    // 插入收藏公司表
    favorite := FavoriteCompany{
        UserID:    userID,
        CompanyID: uint(company.CompanyID),
    }
    if err := oc.DB.Create(&favorite).Error; err != nil {
        return "收藏失败：" + err.Error()
    }

    return "收藏成功！"
}

// UnCollectCompany 取消收藏公司
func (oc *OfferController) UnCollectCompany(userName, companyName string) string {
    // 根据用户名字查找用户ID
    userID, err := oc.GetUserByName(userName)
    if err != nil {
        return "取消收藏失败：" + err.Error()
    }

    // 查询招聘公司表，根据公司名称找到对应的公司ID
    var company utils.RecruitmentCompanies
    if err := oc.DB.Where("company_name LIKE ?", "%"+companyName+"%").First(&company).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return "取消收藏失败：公司名不存在"
        }
        return "取消收藏失败：" + err.Error()
    }

    // 查询是否在收藏列表内
    var existingFavorite FavoriteCompany
    if err := oc.DB.Where("user_id = ? AND company_id = ?", userID, company.CompanyID).First(&existingFavorite).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return "取消收藏失败：不在收藏列表内"
        }
        return "取消收藏失败：" + err.Error()
    }

    // 删除收藏公司表中的记录
    if err := oc.DB.Where("user_id = ? AND company_id = ?", userID, company.CompanyID).Delete(&FavoriteCompany{}).Error; err != nil {
        return "取消收藏失败：" + err.Error()
    }

    return "取消收藏成功！"
}

// 查看已收藏的企业
func (oc *OfferController) GetFavoriteCompanies(userName string) (string, error) {
    // 根据用户名字查找用户ID
    userID, err := oc.GetUserByName(userName)
    if err != nil {
        return "", err
    }

    // 查询收藏公司表，根据用户ID找到收藏的公司ID列表
    var favorites []FavoriteCompany
    if err := oc.DB.Where("user_id = ?", userID).Find(&favorites).Error; err != nil {
        return "", err
    }

    // 查询招聘公司表，根据公司ID列表找到对应的公司名称列表
    var companyNames []string
    for _, favorite := range favorites {
        var company utils.RecruitmentCompanies
        if err := oc.DB.Where("company_id = ?", favorite.CompanyID).First(&company).Error; err != nil {
            if errors.Is(err, gorm.ErrRecordNotFound) {
                return "", errors.New("Company not found")
            }
            return "", err
        }
        companyNames = append(companyNames, company.CompanyName)
    }

    content := strings.Join(companyNames, "\n")

    return content, nil
}