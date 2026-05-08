package handler

import (
	"github.com/gin-gonic/gin"

	"his-go/internal/statistics/service"
	"his-go/pkg/errors"
	"his-go/pkg/response"
)

// StatisticsHandler 数据统计接口处理器
type StatisticsHandler struct {
	svc *service.StatisticsService
}

// NewStatisticsHandler 创建数据统计接口处理器
func NewStatisticsHandler(svc *service.StatisticsService) *StatisticsHandler {
	return &StatisticsHandler{svc: svc}
}

// DateRangeRequest 日期范围请求
type DateRangeRequest struct {
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
}

// GetOperationStats 获取运营统计
func (h *StatisticsHandler) GetOperationStats(c *gin.Context) {
	var req DateRangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	stats, err := h.svc.GetOperationStats(req.StartDate, req.EndDate)
	if err != nil {
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.Success(c, stats)
}

// GetDeptWorkload 获取科室工作负载
func (h *StatisticsHandler) GetDeptWorkload(c *gin.Context) {
	var req DateRangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	list, err := h.svc.GetDeptWorkload(req.StartDate, req.EndDate)
	if err != nil {
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.Success(c, list)
}

// GetRevenueTrend 获取收入趋势
func (h *StatisticsHandler) GetRevenueTrend(c *gin.Context) {
	var req DateRangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	list, err := h.svc.GetRevenueTrend(req.StartDate, req.EndDate)
	if err != nil {
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.Success(c, list)
}
