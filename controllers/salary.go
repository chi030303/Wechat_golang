package controllers

import (
	"Wechat-project/utils"
	"strings"

	"gorm.io/gorm"
)

type SalaryController struct {
	DB *gorm.DB
}

func NewSalaryController(db *gorm.DB) *SalaryController {
	return &SalaryController{DB: db}
}

// 查询具体的薪资信息
func (oc *SalaryController) GetSalaryInfo(companyName string) string {
	var salary utils.CompanySalaries
	var result strings.Builder

	// 查询数据库，获取符合条件的公司信息
	if err := oc.DB.Where("company_name LIKE ?", "%"+companyName+"%").First(&salary).Error; err != nil {
		// 处理查询错误
		if err == gorm.ErrRecordNotFound {
			return "未找到名称包含 \"" + companyName + "\" 的公司信息"
		}
		// 其他错误处理
		return "查询公司信息时出现错误"
	}

	// 组织返回信息
	result.WriteString("公司名称：" + salary.CompanyName + "\n")
	if salary.CompanyLocation != "" {
		result.WriteString("公司所在城市：" + salary.CompanyLocation + "\n")
	}
	if salary.JobPosition != "" {
		result.WriteString("职位：" + salary.JobPosition + "\n")
	}
	if salary.Majors != "" {
		result.WriteString("招聘专业：" + salary.Majors + "\n")
	}
	if salary.EducationBackground != "" {
		result.WriteString("学历：" + salary.EducationBackground + "\n")
	}
	if salary.Salary != "" {
		result.WriteString("薪资：" + salary.Salary + "\n")
	}
	if salary.CompanyProfile != "" {
		result.WriteString("公司简介：" + salary.CompanyProfile + "\n")
	}

	return result.String()
}
