package attachment

import "github.com/gin-gonic/gin"

func (ctrl *Controller) ScanInvalid(c *gin.Context) {

	_, _ = ctrl.attachmentService.ScanInvalid()
}
