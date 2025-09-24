package group

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xinde/pkg/logger"
	"xinde/pkg/response"
	"xinde/pkg/stderr"
)

// GetTree handles fetching the group tree structure.
// @Summary      获取树状分组列表。用于前台展示分组，和后台需要树状分组的地方。
// @Description  获取一个完整的、嵌套的树状分组结构
// @Tags         Group
// @Tags         Solution
// @Accept       json
// @Produce      json
// @Param        icon query string true "是否包含图标URL (必填，true或者false)"
// @Security     ApiKeyAuth
// @Success      200 {object} dto.TreeResp "成功返回分组树"
// @Failure      500 {object} response.Response "服务器内部错误"
// @Router       /api/v1/admin/group/tree [get]
// @Router       /api/v1/groups/tree [get]
func (ctrl *Controller) GetTree(c *gin.Context) {
	//var req dto.TreeReq
	//if err := c.ShouldBind(&req); err != nil {
	//	response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, "绑定参数错误: "+err.Error())
	//	logger.Error("/admin/group/tree 绑定参数错误: " + err.Error())
	//	return
	//}
	includeIcon := "true"
	includeIcon = c.Query("include_icon")

	// 剩余的工作交由service处理
	tree, err := ctrl.Service.GetTree(includeIcon)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternalError, stderr.ErrorInternalServerError)
		logger.Error("/admin/group/tree 获取树状分组列表出错: " + err.Error())
		return
	}
	response.Success(c, tree)
}
