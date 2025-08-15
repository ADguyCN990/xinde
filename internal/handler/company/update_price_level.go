package company

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	dto "xinde/internal/dto/company"
	"xinde/internal/handler/common"
	"xinde/pkg/logger"
	"xinde/pkg/response"
	"xinde/pkg/stderr"
)

// UpdatePriceLevel handles admin update company's price level
// @Summary 修改价格等级
// @Description 管理员根据公司ID修改公司的价格等级
// @Tags Company
// @Accept json
// @Produce json
// @Param request body dto.UpdatePriceLevelReq true "UpdatePriceLevel　Request"
// @Success 200 {object} response.Response "修改价格等级成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "access_token有错误"
// @Failure 403 {object} response.Response "没有管理员权限"
// @Failure 404 {object} response.Response "没有该公司"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/admin/company/price/level/{id} [patch]
func (ctrl *Controller) UpdatePriceLevel(c *gin.Context) {

	var req dto.UpdatePriceLevelReq
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, err.Error())
		logger.Error("/admin/company/price/level/patch 绑定参数错误: " + err.Error())
		return
	}

	id, err := common.GetIDFromUrl(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, stderr.ErrorCompanyIDInvalid)
		logger.Error("/admin/company/price/level/patch 无效的公司ID格式: " + err.Error())
		return
	}

	// 完成参数校验，将剩余的工作交由service处理
	err = ctrl.companyService.UpdatePriceLevel(id, req.PriceLevel)
	if err != nil {
		switch err.Error() {
		case stderr.ErrorCompanyNotFound:
			response.Error(c, http.StatusNotFound, response.CodeNotFound, stderr.ErrorCompanyNotFound)
			logger.Error(fmt.Sprintf("/admin/company/price/level/patch 修改公司价格等级失败! 公司ID: %d 错误: %s", id, err.Error()))
		default:
			response.Error(c, http.StatusInternalServerError, response.CodeInternalError, stderr.ErrorInternalServerError)
			logger.Error(fmt.Sprintf("/admin/company/price/level/patch 修改公司价格等级失败! 公司ID: %d 错误: %s", id, err.Error()))
		}
		return
	}

	response.Success(c, nil)
}
