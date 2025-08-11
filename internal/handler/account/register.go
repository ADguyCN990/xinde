package account

import (
	"github.com/gin-gonic/gin"
	"net/http"
	dto "xinde/internal/dto/account"
	service "xinde/internal/service/account"
	"xinde/pkg/response"
)

type Controller struct {
	accountService *service.Service
}

func NewAccountController() (*Controller, error) {
	accountService, err := service.NewAccountService()
	if err != nil {
		return nil, err
	}

	return &Controller{
		accountService: accountService,
	}, nil
}

// Register handles user registration.
// @Summary 注册一个新用户
// @Description 用户名，真实姓名，公司名称，公司地址（可选），密码，手机号，邮箱（可选）
// @Tags Account
// @Accept json
// @Produce json
// @Param request body dto.RegisterReq true "Register Request"
// @Success 200 {object} dto.RegisterResp "User registered successfully"
// @Failure 400 {object} dto.RegisterResp "Bad Request"
// @Failure 403 {object} dto.RegisterResp "Forbidden"
// @Failure 500 {object} dto.RegisterResp "Internal Server Error"
// @Router /api/v1/account/register [post]
func (ctrl *Controller) Register(c *gin.Context) {
	var req dto.RegisterReq
	var err error
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, err.Error())
		return
	}

	// 重复输入密码需一致
	if req.Password != req.ConfirmedPassword {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, "两次输入的密码不一致，请重新输入")
		return
	}

	// 校验完参数后，交由service层处理
	_, err = ctrl.accountService.Register(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternalError, err.Error())
		return
	}
	response.Success(c, nil)
}
