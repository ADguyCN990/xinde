package group

import (
	"fmt"
	"gorm.io/gorm"
	"xinde/internal/dao/common"
	model "xinde/internal/model/group"
	"xinde/internal/store"
	"xinde/pkg/stderr"
)

type Dao struct {
	db        *gorm.DB
	commonDao *common.Dao
}

func (d *Dao) DB() *gorm.DB {
	return d.db
}

func NewGroupDao() (*Dao, error) {
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

func (d *Dao) Create(tx *gorm.DB, groupName string, parentID uint) (uint, error) {
	if tx == nil {
		return 0, fmt.Errorf(stderr.ErrorDbNil)
	}

	g := &model.Group{
		Name:     groupName,
		ParentID: parentID,
	}
	if err := tx.Model(&model.Group{}).Create(g).Error; err != nil {
		return 0, fmt.Errorf("Dao层创建分组失败: " + err.Error())
	}

	return g.ID, nil
}

func (d *Dao) GetGroupByID(tx *gorm.DB, groupID uint) (*model.Group, error) {
	if tx == nil {
		return nil, fmt.Errorf(stderr.ErrorDbNil)
	}
	var group *model.Group
	if err := tx.Model(&model.Group{}).Where("id = ?", groupID).First(&group).Error; err != nil {
		return nil, fmt.Errorf("Dao层根据ID查找分组失败: " + err.Error())
	}
	return group, nil
}

func (d *Dao) GetGroupsByIDs(tx *gorm.DB, groupIDs []uint) ([]*model.Group, error) {
	if tx == nil {
		return nil, fmt.Errorf(stderr.ErrorDbNil)
	}
	var groups []*model.Group
	if err := tx.Model(&model.Group{}).Where("id IN (?)", groupIDs).Find(&groups).Error; err != nil {
		return nil, fmt.Errorf("Dao层根据ID列表查找分组失败: " + err.Error())
	}
	return groups, nil
}

func (d *Dao) UpdateGroupByID(tx *gorm.DB, groupID uint, updateData map[string]interface{}) error {
	if tx == nil {
		return fmt.Errorf(stderr.ErrorDbNil)
	}
	if err := tx.Model(&model.Group{}).Where("id = ?", groupID).Updates(updateData).Error; err != nil {
		return fmt.Errorf("Dao层更新分组失败: " + err.Error())
	}
	return nil
}

func (d *Dao) GetAll(tx *gorm.DB) ([]*model.Group, error) {
	if tx == nil {
		return nil, fmt.Errorf(stderr.ErrorDbNil)
	}
	var groups []*model.Group
	if err := tx.Model(&model.Group{}).Order("id asc").Find(&groups).Error; err != nil {
		return nil, fmt.Errorf("Dao层查找所有分组失败: " + err.Error())
	}
	return groups, nil
}

// FindAllDescendantIDs 接收一个 groupID，并返回一个包含其所有子孙ID 的列表。
func (d *Dao) FindAllDescendantIDs(tx *gorm.DB, groupID uint) ([]uint, error) {
	if tx == nil {
		return nil, fmt.Errorf(stderr.ErrorDbNil)
	}
	var descendantIDs []uint

	// 使用原生SQL的递归查询
	// 1. 初始集是 groupID 自身。
	// 2. 递归部分是不断地将子节点的 id 加入结果集。
	// 3. 最后我们选择除 groupID 自身外的所有 id。
	sql := `
		WITH RECURSIVE descendants AS (
			SELECT id FROM t_group WHERE id = ?
			UNION ALL
			SELECT g.id FROM t_group g JOIN descendants d ON g.parent_id = d.id
		)
		SELECT id FROM descendants WHERE id != ?;
	`

	err := tx.Raw(sql, groupID, groupID).Scan(&descendantIDs).Error
	if err != nil {
		return nil, fmt.Errorf("查找子孙分组失败: " + err.Error())
	}

	return descendantIDs, nil
}

func (d *Dao) DeleteGroupsByIDs(tx *gorm.DB, groupIDs []uint) error {
	if tx == nil {
		return fmt.Errorf(stderr.ErrorDbNil)
	}
	if len(groupIDs) == 0 {
		return nil
	}
	if err := tx.Delete(&model.Group{}, groupIDs).Error; err != nil {
		return fmt.Errorf("删除分组失败: " + err.Error())
	}
	return nil
}
