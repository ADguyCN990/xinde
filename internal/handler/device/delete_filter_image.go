package device

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"xinde/internal/handler/common"
	"xinde/pkg/logger"
	"xinde/pkg/response"
	"xinde/pkg/stderr"
)

// DeleteFilterImage handles the deletion of a filter image configuration.
// @Summary      删除筛选条件图片
// @Description  删除一条筛选条件图片配置及其关联的图片文件
// @Tags         FilterImage
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "筛选条件图片配置 ID"
// @Security     ApiKeyAuth
// @Success      200 {object} response.Response "删除成功"
// @Failure      400 {object} response.Response "请求参数错误或无效ID"
// @Failure      404 {object} response.Response "配置不存在"
// @Failure      500 {object} response.Response "服务器内部错误"
// @Router       /api/v1/admin/filter_image/{id} [delete]
func (ctrl *Controller) DeleteFilterImage(c *gin.Context) {
	id, err := common.GetIDFromUrl(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, stderr.ErrorFilterImageIDInvalid)
		logger.Error("/admin/filter_image/delete 无效的下拉筛选列表图片ID格式: " + err.Error())
		return
	}

	// 剩余的工作交由service处理
	err = ctrl.service.DeleteFilterImage(id)
	if err != nil {
		switch err.Error() {
		case gorm.ErrRecordNotFound.Error():
			response.Error(c, http.StatusNotFound, response.CodeNotFound, stderr.ErrorFilterImageNotFound)
		default:
			response.Error(c, http.StatusInternalServerError, response.CodeInternalError, stderr.ErrorInternalServerError)
			logger.Error("/admin/filter_image/delete 删除下拉筛选列表失败: " + err.Error())
		}
		return
	}
	response.Success(c, nil)
}
