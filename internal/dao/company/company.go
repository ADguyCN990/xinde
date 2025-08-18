package company

import (
	"fmt"
	"gorm.io/gorm"
	"xinde/internal/dao/common"
	model "xinde/internal/model/account"
	"xinde/internal/store"
	"xinde/pkg/stderr"
)

type Dao struct {
	db        *gorm.DB
	commonDao *common.Dao
}

func NewCompanyDao() (*Dao, error) {
	db := store.GetDB()
	if db == nil {
		return nil, fmt.Errorf("数据库连接未初始化，请先调用 store.InitDB()")
	}

	commonDao, err := common.NewCommonDao()
	if err != nil {
		return nil, err
	}

	return &Dao{
		db:        db,
		commonDao: commonDao,
	}, nil
}

// DB 返回原始的 gorm.DB 实例，以便 Service 层可以开启事务
func (d *Dao) DB() *gorm.DB {
	return d.db
}

// CountCompanies 统计公司总数
func (d *Dao) CountCompanies(tx *gorm.DB) (int64, error) {
	if tx == nil {
		return 0, fmt.Errorf(stderr.ErrorDbNil)
	}

	var count int64
	err := tx.Model(&model.Company{}).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("统计公司总数失败: " + err.Error())
	}
	return count, nil
}

// GetCompanyByID 根据公司ID获取公司
func (d *Dao) GetCompanyByID(tx *gorm.DB, id uint) (*model.Company, error) {
	if tx == nil {
		return nil, fmt.Errorf(stderr.ErrorDbNil)
	}

	var company model.Company
	err := tx.Model(&model.Company{}).Where("id = ?", id).First(&company).Error
	if err != nil {
		return nil, fmt.Errorf("查询公司失败" + err.Error())
	}
	return &company, nil
}

// IsExistCompanyByID 根据ID判断公司是否存在
func (d *Dao) IsExistCompanyByID(tx *gorm.DB, id uint) (bool, error) {
	if tx == nil {
		return false, fmt.Errorf(stderr.ErrorDbNil)
	}

	var count int64
	err := tx.Model(&model.Company{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, nil
	}
	return count > 0, nil
}

// FindCompanyListWithPagination 分页查找公司列表
func (d *Dao) FindCompanyListWithPagination(tx *gorm.DB, page, pageSize int) ([]*model.Company, error) {
	if tx == nil {
		return nil, fmt.Errorf(stderr.ErrorDbNil)
	}

	var companies []*model.Company
	offset := (page - 1) * pageSize
	err := tx.Model(&model.Company{}).Order("id asc").
		Limit(pageSize).Offset(offset).
		Find(&companies).Error
	if err != nil {
		return nil, fmt.Errorf("查找公司列表失败: " + err.Error())
	}
	return companies, nil
}

func (d *Dao) UpdateCompany(tx *gorm.DB, id uint, updateData map[string]interface{}) error {
	if tx == nil {
		return fmt.Errorf(stderr.ErrorDbNil)
	}

	err := tx.Model(&model.Company{}).Where("id = ?", id).Updates(updateData).Error
	if err != nil {
		return fmt.Errorf("更新公司失败: " + err.Error())
	}
	return nil
}
