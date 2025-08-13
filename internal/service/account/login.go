package account

import (
	"fmt"
	dto "xinde/internal/dto/account"
	"xinde/pkg/stderr"
	"xinde/pkg/util"
)

func (s *Service) Login(username, password string) (*dto.LoginData, error) {
	tx := s.dao.DB()

	// 根据用户名查找用户
	user, err := s.dao.FindUserByUsername(tx, username)
	if err != nil {
		return nil, err
	}

	// 验证用户密码是否正确
	isPasswordOK := util.CheckPasswordHash(password, user.Password)
	if !isPasswordOK {
		return nil, fmt.Errorf(stderr.ErrorUserUnauthorized)
	}

	// 验证管理员是否通过了用户的注册申请
	switch user.IsUser {
	case 0:
		return nil, fmt.Errorf(stderr.ErrorUserNotPass)
	case 2:
		return nil, fmt.Errorf(stderr.ErrorUserBanned)
	}

	// 一切正常，生成JWT Token
	token, err := s.jwt.GenerateToken(user.UID, user.Username, user.IsAdmin == 1)
	if err != nil {
		return nil, fmt.Errorf("user: %s Login, 生成token错误：%s", user.Username, err.Error())
	}

	loginData := dto.LoginData{
		Username:    user.Username,
		Name:        user.Name,
		Phone:       user.Phone,
		Email:       *user.UserEmail,
		AccessToken: token,
	}
	return &loginData, nil

}
