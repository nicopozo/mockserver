package controller

import (
	"net/http"

	mockscontext "github.com/nicopozo/mockserver/internal/context"
	"github.com/nicopozo/mockserver/internal/model"
	"github.com/nicopozo/mockserver/internal/service"
	httputils "github.com/nicopozo/mockserver/internal/utils/http"
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

// GetLogs returns captured log entries with pagination.
func (controller *LogController) GetLogs(writer http.ResponseWriter, request *http.Request) {
	reqContext := mockscontext.New(request)
	logger := mockscontext.Logger(reqContext)

	logger.Debug(controller, nil, "Entering LogController GetLogs()")

	paging, err := getPagingFromRequest(request)
	if err != nil {
		logger.Error(controller, nil, err, "Error parsing pagination params")
		httputils.WriteError(writer, model.ValidationError, "Error parsing pagination params: %s", err.Error())

		return
	}

	logs := controller.LogService.GetAll(*paging)

	httputils.WriteJSON(writer, http.StatusOK, logs)
}

// ClearLogs deletes all captured log entries.
func (controller *LogController) ClearLogs(writer http.ResponseWriter, request *http.Request) {
	reqContext := mockscontext.New(request)
	logger := mockscontext.Logger(reqContext)

	logger.Debug(controller, nil, "Entering LogController ClearLogs()")

	controller.LogService.Clear()
	writer.WriteHeader(http.StatusNoContent)
}
