package account

import (
	"errors"
	"fmt"
	registerDao "xinde/internal/dao/account"
	dto "xinde/internal/model/dto/account"
	"xinde/pkg/util"
)

type AccountService struct {
	dao *registerDao.RegisterDao
}

func NewAccountService() (*AccountService, error) {
	dao, err := registerDao.NewRegisterDao()
	if err != nil {
		return nil, fmt.Errorf("创建 DAO 实例失败: %w", err)
	}

	return &AccountService{
		dao: dao,
	}, nil
}

func (s *AccountService) Register(req *dto.RegisterReq) (uint, error) {
	// 检查用户是否存在
	exists, err := s.dao.IsExistUser(req.Username)
	if err != nil {
		return 0, fmt.Errorf("检查用户是否存在失败: %w", err)
	}

	if exists {
		return 0, errors.New("用户已存在，注册失败")
	}

	// 加密密码
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return 0, fmt.Errorf("密码加密失败: %w", err)
	}

	// 写入数据库
	userID, err := s.dao.CreateUser(req.Username, req.Email, req.Name, req.CompanyName, req.CompanyAddress, hashedPassword, req.Phone)
	if err != nil {
		return 0, fmt.Errorf("创建用户失败: %w", err)
	}

	return userID, nil
}
