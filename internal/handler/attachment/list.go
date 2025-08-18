package attachment

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	dto "xinde/internal/dto/attachment"
	"xinde/internal/service/attachment"
	"xinde/pkg/logger"
	"xinde/pkg/response"
	"xinde/pkg/stderr"
)

type Controller struct {
	attachmentService *attachment.Service
}

func NewAttachmentController() (*Controller, error) {
	service, err := attachment.NewAttachmentService()
	if err != nil {
		return nil, fmt.Errorf("创建Service实例失败: " + err.Error())
	}
	return &Controller{attachmentService: service}, nil
}

// List handles fetching a paginated list of attachments.
// @Summary      获取附件列表
// @Description  根据筛选条件分页获取已上传的附件列表
// @Tags		 Attachment
// @Accept       json
// @Produce      json
// @Param        page query int true "页码"
// @Param        pageSize query int false "每页数量"
// @Param        business_type query string false "业务类型 (e.g., price_import)"
// @Param        filename query string false "文件名 (模糊搜索)"
// @Security     ApiKeyAuth
// @Success      200 {object} dto.ListResp "成功返回附件列表"
// @Failure      400 {object} response.Response "请求参数错误"
// @Failure      401 {object} response.Response "Token错误"
// @Failure      403 {object} response.Response "没有管理员权限"
// @Failure      500 {object} response.Response "服务器内部错误"
// @Router       /api/v1/admin/attachment/list [get]
func (ctrl *Controller) List(c *gin.Context) {
	var req *dto.ListReq
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, "绑定参数错误: "+err.Error())
		logger.Error("/admin/attachment/list 绑定参数错误: " + err.Error())
		return
	}

	if req.PageSize == 0 {
		req.PageSize = viper.GetInt("page.defaultPageSize")
	}

	// 参数校验完毕，剩余的工作交由service处理
	list, err := ctrl.attachmentService.List(req.Page, req.PageSize, req.Filename)
	if err != nil {
		switch err.Error() {
		case stderr.ErrorDbNil:
			response.Error(c, http.StatusInternalServerError, response.CodeInternalError, stderr.ErrorInternalServerError)
			logger.Error("/admin/attachment/list " + err.Error())
		case stderr.ErrorOverLargePage:
			response.SuccessWithMessage(c, fmt.Sprintf("%s, 跳转至最后一页", stderr.ErrorOverLargePage), list)
		case stderr.ErrorOverSmallPage:
			response.SuccessWithMessage(c, fmt.Sprintf("%s, 跳转至第一页", stderr.ErrorOverSmallPage), list)
		default:
			response.Error(c, http.StatusInternalServerError, response.CodeInternalError, stderr.ErrorInternalServerError)
			logger.Error("/admin/attachment/list " + err.Error())
		}
		return
	}
	response.Success(c, list)
}
