package group

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	dto "xinde/internal/dto/group"
	"xinde/internal/middleware/auth"
	"xinde/internal/service/group"
	"xinde/pkg/logger"
	"xinde/pkg/response"
	"xinde/pkg/stderr"
)

type Controller struct {
	Service *group.Service
}

func NewGroupController() (*Controller, error) {
	service, err := group.NewGroupService()
	if err != nil {
		return nil, err
	}
	return &Controller{
		Service: service,
	}, nil
}

// Create handles the creation of a new group.
// @Summary      创建新分组
// @Description  创建一个新分组，可选择性上传图标
// @Tags         Group
// @Accept       multipart/form-data
// @Produce      json
// @Param        name formData string true "分组名称"
// @Param        parent_id formData int true "父级分组ID"
// @Param        icon formData file false "分组图标 (可选)"
// @Security     ApiKeyAuth
// @Success      200 {object} response.Response "创建成功"
// @Failure      400 {object} response.Response "请求参数错误"
// @Failure      401 {object} response.Response "Token错误"
// @Failure      403 {object} response.Response "没有管理员权限"
// @Failure      500 {object} response.Response "服务器内部错误"
// @Router       /api/v1/admin/group/create [post]
func (ctrl *Controller) Create(c *gin.Context) {
	var req dto.CreateReq
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, "绑定参数错误: "+err.Error())
		logger.Error("/admin/group/create 绑定参数错误: " + err.Error())
		return
	}
	// 从表单中获取上传的icon（可选，如果没有iconFile为nil）
	iconFile, _ := c.FormFile("icon")

	// 获取当前操作的管理员ID
	adminID, err := auth.GetCurrentUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "无法获取当前的用户ID")
		logger.Error("/admin/group/create 无法获取当前的用户ID: " + err.Error())
	}

	if !auth.IsAdmin(c) {
		response.Error(c, http.StatusForbidden, response.CodeForbidden, stderr.ErrorTokenNotAdmin)
		logger.Error(fmt.Sprintf("/admin/group/create 当前用户ID: %d 非管理员，权限不足", adminID))
	}

	// 剩余工作交由Service层处理
	err = ctrl.Service.Create(req.Name, req.ParentID, adminID, iconFile)
	if err != nil {
		switch err.Error() {
		default:
			response.Error(c, http.StatusInternalServerError, response.CodeInternalError, stderr.ErrorInternalServerError)
			logger.Error("创建分组失败: " + err.Error())
			return
		}
	}
	response.Success(c, nil)
}
