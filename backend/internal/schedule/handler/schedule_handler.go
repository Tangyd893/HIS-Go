package handler

import (
	"github.com/gin-gonic/gin"

	"his-go/internal/schedule/model"
	"his-go/internal/schedule/service"
	"his-go/pkg/errors"
	"his-go/pkg/response"
)

// ScheduleHandler 排班接口处理器
type ScheduleHandler struct {
	svc *service.ScheduleService
}

// NewScheduleHandler 创建排班接口处理器
func NewScheduleHandler(svc *service.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{svc: svc}
}

// GenerateSchedulesRequest 生成排班请求
type GenerateSchedulesRequest struct {
	StartDate string `json:"startDate" binding:"required"`
	EndDate   string `json:"endDate" binding:"required"`
	DeptID    string `json:"deptId" binding:"required"`
}

// GenerateWeeklySchedules 生成一周排班
func (h *ScheduleHandler) GenerateWeeklySchedules(c *gin.Context) {
	var req GenerateSchedulesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	schedules, err := h.svc.GenerateWeekly(req.StartDate, req.EndDate, req.DeptID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			response.FailWithMsg(c, appErr.Code, appErr.Message)
			return
		}
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.Success(c, schedules)
}

// ListSchedules 查询排班列表
func (h *ScheduleHandler) ListSchedules(c *gin.Context) {
	deptID := c.Query("deptId")
	date := c.Query("date")
	doctorID := c.Query("doctorId")

	if doctorID != "" {
		list, err := h.svc.ListByDoctor(doctorID, date)
		if err != nil {
			response.FailWithMsg(c, errors.CodeInternalError, err.Error())
			return
		}
		response.Success(c, list)
		return
	}

	list, err := h.svc.ListByDept(deptID, date)
	if err != nil {
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.Success(c, list)
}

// UpdateSchedule 更新排班
func (h *ScheduleHandler) UpdateSchedule(c *gin.Context) {
	var schedule model.ScheduleInfo
	if err := c.ShouldBindJSON(&schedule); err != nil {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	if err := h.svc.UpdateSchedule(&schedule); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			response.FailWithMsg(c, appErr.Code, appErr.Message)
			return
		}
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.Success(c, nil)
}

// CancelSchedule 取消排班
func (h *ScheduleHandler) CancelSchedule(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	if err := h.svc.CancelSchedule(id); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			response.FailWithMsg(c, appErr.Code, appErr.Message)
			return
		}
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.Success(c, nil)
}
