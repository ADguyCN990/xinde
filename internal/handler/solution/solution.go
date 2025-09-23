package solution

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	dto "xinde/internal/dto/solution"
	"xinde/internal/middleware/auth"
	"xinde/internal/service/solution"
	"xinde/pkg/logger"
	"xinde/pkg/response"
	"xinde/pkg/stderr"
)

type Controller struct {
	service *solution.Service
}

func NewSolutionController() (*Controller, error) {
	service, err := solution.NewSolutionService()
	if err != nil {
		return nil, fmt.Errorf("NewSolutionService() 创建service实例失败: %v", err)
	}
	return &Controller{service: service}, nil
}

// Query handles the dynamic querying of solutions.
// @Summary      查询/筛选选型方案
// @Description  根据用户提供的筛选条件，动态查询方案列表，并返回下一步可用的筛选选项。传入空的筛选对象可获取初始状态。
// @Tags         Solution
// @Accept       json
// @Produce      json
// @Param        body body dto.QueryReq true "查询请求体"
// @Security     ApiKeyAuth
// @Success      200 {object} response.Response{data=dto.QueryResp} "查询成功"
// @Failure      400 {object} response.Response "请求参数错误"
// @Failure      401 {object} response.Response "Token错误"
// @Failure      500 {object} response.Response "服务器内部错误"
// @Router       /api/v1/solutions/query [post]
func (ctrl *Controller) Query(c *gin.Context) {
	var req *dto.QueryReq
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, "绑定参数错误: "+err.Error())
		logger.Error("/solutions/query 绑定参数错误: " + err.Error())
		return
	}

	userID, err := auth.GetCurrentUserID(c)
	if err != nil {
		response.Error(c, http.StatusForbidden, response.CodeForbidden, "无法获取当前的用户ID: "+err.Error())
		logger.Error("/solutions/query 无法获取当前的用户ID: " + err.Error())
		return
	}

	resp, err := ctrl.service.Query(userID, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternalError, stderr.ErrorInternalServerError)
		logger.Error("/solutions/query 查询方案失败: " + err.Error())
		return
	}
	response.Success(c, resp)
}
