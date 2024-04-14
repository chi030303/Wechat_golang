package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
	"Wechat-project/config"
	"gorm.io/gorm"
)
// 招聘企业薪资结构体
type CompanySalaries struct {
	CompanyID             int    `json:"id" gorm:"company_id"`
	CompanyName    string `json:"公司名称" gorm:"company_name"`
	CompanyLocation   string `json:"工作城市" gorm:"company_location"`
	JobPosition    string `json:"岗位" gorm:"job_position"`
	EducationBackground      string `json:"学历要求" gorm:"education_background"`
	Salary         string `json:"薪资" gorm:"salary"`
	Majors          string `json:"专业" gorm:"majors"`
	CompanyProfile string `json:"公司简介" gorm:"company_profile"`
}

func (CompanySalaries) TableName() string {
	return "Company_Salaries"
}

// 招聘企业结构体
type RecruitmentCompanies struct {
	CompanyID                     int    `json:"id" gorm:"company_id"`
	CompanyName            string `json:"公司名称" gorm:"company_name"`
	RecruitmentTarget      string `json:"招聘对象" gorm:"recruitment_target"`
	JobPosition    		   string `json:"招聘岗位" gorm:"job_position"`
	RecruitmentMajors string `json:"招聘专业/条件/要求" gorm:"recruitment_majors"`
	WorkLocation           string `json:"工作城市" gorm:"work_location"`
	RecruitmentSession     string `json:"招聘届数" gorm:"recruitment_session"`
	DeliveryTime           string `json:"投递时间" gorm:"delivery_time"`
}

func (RecruitmentCompanies) TableName() string {
	return "Recruitment_Companies"
}

// 插入爬虫数据
func InsertData() error {
	// 初始化数据库连接
	_, err := config.InitDB()
	if err != nil {
		return err
	}
	db := config.DB

	// 读取薪资数据并插入数据库
	if err := insertSalaries(db); err != nil {
		return err
	}

	// 读取招聘公司数据并插入数据库
	if err := insertRecruitmentCompanies(db); err != nil {
		return err
	}

	return nil
}

// 插入薪资
func insertSalaries(db *gorm.DB) error {
	// 读取 JSON 文件
	filePath := filepath.Join("data", "salaries_data.json")
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 解码 JSON 数据
	var salaries []CompanySalaries
	err = json.NewDecoder(file).Decode(&salaries)
	if err != nil {
		return err
	}

	// 将数据插入数据库
	for _, salary := range salaries {
		result := db.Create(&salary)
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}

// 插入招聘企业
func insertRecruitmentCompanies(db *gorm.DB) error {
	// 读取 JSON 文件
	filePath := filepath.Join("data", "companies_data.json")
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 解码 JSON 数据
	var companies []RecruitmentCompanies
	err = json.NewDecoder(file).Decode(&companies)
	if err != nil {
		return err
	}

	// 将数据插入数据库
	for _, company := range companies {
		result := db.Create(&company)
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}