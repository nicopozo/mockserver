package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	mockscontext "github.com/nicopozo/mockserver/internal/context"
	"github.com/nicopozo/mockserver/internal/service"
)

// LogController exposes the in-memory request/response logs.
type LogController struct {
	LogService service.LogService
}

func NewLogController(logService service.LogService) *LogController {
	return &LogController{
		LogService: logService,
	}
}

// GetLogs returns all captured log entries (newest first).
func (controller *LogController) GetLogs(context *gin.Context) {
	reqContext := mockscontext.New(context)
	logger := mockscontext.Logger(reqContext)

	logger.Debug(controller, nil, "Entering LogController GetLogs()")

	logs := controller.LogService.GetAll()
	context.JSON(http.StatusOK, logs)
}

// ClearLogs deletes all captured log entries.
func (controller *LogController) ClearLogs(context *gin.Context) {
	reqContext := mockscontext.New(context)
	logger := mockscontext.Logger(reqContext)

	logger.Debug(controller, nil, "Entering LogController ClearLogs()")

	controller.LogService.Clear()
	context.Status(http.StatusNoContent)
}
