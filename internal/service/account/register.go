package account

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	registerDao "xinde/internal/dao/account"
	dto "xinde/internal/dto/account"
	"xinde/pkg/jwt"
	"xinde/pkg/stderr"
	"xinde/pkg/util"
)

type Service struct {
	dao *registerDao.Dao
	jwt *jwt.JWTService
}

func NewAccountService() (*Service, error) {
	dao, err := registerDao.NewRegisterDao()
	if err != nil {
		return nil, fmt.Errorf("创建 DAO 实例失败: %w", err)
	}

	jwtService := jwt.NewJWTService()

	return &Service{
		dao: dao,
		jwt: jwtService,
	}, nil
}

func (s *Service) Register(req *dto.RegisterReq) (uint, error) {
	// 启动一个事务
	tx := s.dao.DB().Begin() // 假设 s.dao.DB() 可以返回 *gorm.DB 实例
	if tx.Error != nil {
		return 0, fmt.Errorf("开启事务失败: %w", tx.Error)
	}

	// 使用 defer 来确保事务在函数退出时能被处理（提交或回滚）
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback() // 如果发生 panic，回滚事务
		}
	}()

	// --- 在事务中执行所有数据库操作 ---

	// 检查用户是否存在
	exists, err := s.dao.IsExistUser(tx, req.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return 0, err
	}
	if exists {
		tx.Rollback()
		return 0, fmt.Errorf(stderr.ErrorUserAlreadyExist)
	}

	// 加密密码
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return 0, fmt.Errorf("密码加密失败: %w", err)
	}

	// 判断公司有没有注册过，如果没有注册过的话，还要新建公司
	companyID, err := s.dao.FindOrCreateCompany(tx, req.CompanyName, req.CompanyAddress)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// 写入数据库
	userID, err := s.dao.CreateUser(tx, req.Username, req.Email, req.Name, req.CompanyName, req.CompanyAddress, hashedPassword, req.Phone, companyID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// 所有操作成功，提交事务
	if err := tx.Commit().Error; err != nil {
		// 提交失败也需要回滚（虽然很少见），并记录错误
		tx.Rollback()
		return 0, fmt.Errorf("提交事务失败: %w", err)
	}

	return userID, nil
}
