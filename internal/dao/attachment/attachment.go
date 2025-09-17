package attachment

import (
	"fmt"
	"gorm.io/gorm"
	"xinde/internal/dao/common"
	model "xinde/internal/model/attachment"
	"xinde/internal/store"
	"xinde/pkg/stderr"
)

type Dao struct {
	db        *gorm.DB
	commonDao *common.Dao
}

type ListParams struct {
	Page     int
	PageSize int
	Filename string
}

func NewAttachmentDao() (*Dao, error) {
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

func (d *Dao) GetAttachmentByID(tx *gorm.DB, id uint) (*model.Attachment, error) {
	if tx == nil {
		return nil, fmt.Errorf(stderr.ErrorDbNil)
	}

	var attachment model.Attachment
	err := tx.Model(&model.Attachment{}).Where("id = ?", id).First(&attachment).Error
	if err != nil {
		return nil, err
	}
	return &attachment, nil
}

func (d *Dao) GetAttachmentsByBusinessType(tx *gorm.DB, businessType string) ([]*model.Attachment, error) {
	if tx == nil {
		return nil, fmt.Errorf(stderr.ErrorDbNil)
	}
	var attachments []*model.Attachment
	err := tx.Model(&model.Attachment{}).Where("business_type = ?", businessType).Find(&attachments).Error
	if err != nil {
		return nil, err
	}
	return attachments, nil
}

func (d *Dao) GetAttachmentByBusinessType(tx *gorm.DB, businessType string, id uint) (*model.Attachment, error) {
	if tx == nil {
		return nil, fmt.Errorf(stderr.ErrorDbNil)
	}
	var attachment *model.Attachment
	err := tx.Model(&model.Attachment{}).Where("business_type = ? AND business_id = ?", businessType, id).First(&attachment).Error
	if err != nil {
		return nil, err
	}
	return attachment, nil
}

func (d *Dao) GetAttachmentsByBusinessAndID(tx *gorm.DB, businessType string, id uint) ([]*model.Attachment, error) {
	if tx == nil {
		return nil, fmt.Errorf(stderr.ErrorDbNil)
	}
	var attachments []*model.Attachment
	err := tx.Model(&model.Attachment{}).Where("business_type = ? AND business_id = ?", businessType, id).Find(&attachments).Error
	if err != nil {
		return nil, fmt.Errorf("Dao层查询失败, 无法根据业务类型和业务ID获取附件列表: " + err.Error())
	}
	return attachments, nil
}

func (d *Dao) GetAttachmentsByBusinessAndIDs(tx *gorm.DB, businessType string, ids []uint) ([]*model.Attachment, error) {
	if tx == nil {
		return nil, fmt.Errorf(stderr.ErrorDbNil)
	}
	var attachments []*model.Attachment
	err := tx.Model(&model.Attachment{}).Where("business_type = ? AND business_id IN (?)", businessType, ids).Find(&attachments).Error
	if err != nil {
		return nil, fmt.Errorf("Dao层查询失败，无法根据业务类型和业务ID列表获取附件列表: " + err.Error())
	}
	return attachments, nil
}

func (d *Dao) DeleteAttachmentByID(tx *gorm.DB, id uint) error {
	if tx == nil {
		return fmt.Errorf(stderr.ErrorDbNil)
	}
	result := tx.Delete(&model.Attachment{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf(stderr.ErrorAttachmentNotFound)
	}
	return nil
}

// DeleteAttachmentsByBusinessTypeAndIDs deletes all attachments matching a business type and a list of business IDs.
func (d *Dao) DeleteAttachmentsByBusinessTypeAndIDs(tx *gorm.DB, bizType string, bizIDs []uint) error {
	if tx == nil {
		return fmt.Errorf(stderr.ErrorDbNil)
	}
	if len(bizIDs) == 0 {
		return nil // 如果没有ID，无需操作
	}

	// 使用 Where(...).Delete(...) 执行批量删除
	err := tx.Delete(&model.Attachment{}, "business_type = ? AND business_id IN ?", bizType, bizIDs).Error

	if err != nil {
		return fmt.Errorf("Dao层根据业务类型和业务ID批量删除附件列表失败:" + err.Error())
	}

	return nil
}

func (d *Dao) CountWithParams(tx *gorm.DB, params *ListParams) (int64, error) {
	if tx == nil {
		return 0, fmt.Errorf(stderr.ErrorDbNil)
	}
	queryBuilder := tx.Model(model.Attachment{})
	if params.Filename != "" {
		queryBuilder = queryBuilder.Where("t_attachment.filename LIKE ?", "%"+params.Filename+"%")
	}
	var count int64
	if err := queryBuilder.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("统计附件总数失败: " + err.Error())
	}
	return count, nil
}

func (d *Dao) GetAllAttachments(tx *gorm.DB) ([]*model.Attachment, error) {
	if tx == nil {
		return nil, fmt.Errorf(stderr.ErrorDbNil)
	}
	var list []*model.Attachment
	err := tx.Model(&model.Attachment{}).
		Select("t_attachment.*, t_user.name as uploader_name").
		Joins("LEFT JOIN t_user ON t_user.uid = t_attachment.uploaded_by_uid").
		Order("t_attachment.id asc").
		Find(&list).Error
	if err != nil {
		return nil, fmt.Errorf("获取所有附件记录失败: " + err.Error())
	}
	return list, nil
}

func (d *Dao) FindAttachmentListWithPagination(tx *gorm.DB, params *ListParams) ([]*model.Attachment, error) {
	if tx == nil {
		return nil, fmt.Errorf(stderr.ErrorDbNil)
	}
	var list []*model.Attachment
	queryBuilder := tx.Model(&model.Attachment{}).
		Select("t_attachment.*, t_user.name as uploader_name").
		Joins("LEFT JOIN t_user ON t_user.uid = t_attachment.uploaded_by_uid")
	if params.Filename != "" {
		queryBuilder = queryBuilder.Where("t_attachment.filename LIKE ?", "%"+params.Filename+"%")
	}

	offset := params.PageSize * (params.Page - 1)
	if err := queryBuilder.Order("t_attachment.id desc").Offset(offset).Limit(params.PageSize).Find(&list).Error; err != nil {
		return nil, fmt.Errorf("分页查找价格列表失败: " + err.Error())
	}
	return list, nil
}

func (d *Dao) Create(tx *gorm.DB, att *model.Attachment) error {
	if tx == nil {
		return fmt.Errorf(stderr.ErrorDbNil)
	}

	err := tx.Model(&model.Attachment{}).Create(att).Error
	if err != nil {
		return fmt.Errorf("创建attachment失败: " + err.Error())
	}
	return nil
}
