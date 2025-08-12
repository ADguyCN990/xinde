package account

import (
	"github.com/gin-gonic/gin"
	"net/http"
	dto "xinde/internal/dto/account"
	"xinde/pkg/response"
	"xinde/pkg/stderr"
)

// Login handles user login.
// @Summary 用户登录
// @Description 用户登录，返回JWT token和用户基本信息
// @Tags Account
// @Accept json
// @Produce json
// @Param request body dto.LoginReq true "Login Request"
// @Success 200 {object} dto.LoginResp "登录成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "用户名或密码错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/account/login [post]
func (ctrl *Controller) Login(c *gin.Context) {
	var req dto.LoginReq
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, err.Error())
		return
	}

	// 校验参数由ShouldBind完成，剩下的交由Service层处理
	loginRespData, err := ctrl.accountService.Login(req.Username, req.Password)
	if err != nil {
		switch err.Error() {
		case stderr.ErrorUserNotPass:
			response.Error(c, http.StatusForbidden, response.CodeForbidden, err.Error())
		case stderr.ErrorUserUnauthorized:
			response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, err.Error())
		case stderr.ErrorUserBanned:
			response.Error(c, http.StatusForbidden, response.CodeForbidden, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, response.CodeInternalError, err.Error())
		}
		return
	}
	response.Success(c, loginRespData)
}
