package handler

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"his-go/internal/registration/service"
	"his-go/pkg/errors"
	"his-go/pkg/response"
)

// RegistrationHandler 挂号接口处理器
type RegistrationHandler struct {
	svc *service.RegistrationService
}

// NewRegistrationHandler 创建挂号接口处理器
func NewRegistrationHandler(svc *service.RegistrationService) *RegistrationHandler {
	return &RegistrationHandler{svc: svc}
}

// RegistrationRequest 挂号请求
type RegistrationRequest struct {
	PatientID   string `json:"patientId" binding:"required"`
	PatientName string `json:"patientName"`
	ScheduleID  string `json:"scheduleId" binding:"required"`
}

// ListRegistrations 分页查询挂号记录（管理端/患者端通用）
// 若传入 patientId 则按患者过滤，否则查询全部
func (h *RegistrationHandler) ListRegistrations(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	date := c.Query("date")
	patientID := c.Query("patientId")

	var list interface{}
	var total int64
	var err error

	if patientID != "" {
		list, total, err = h.svc.ListByPatient(patientID, page, pageSize)
	} else {
		list, total, err = h.svc.ListAll(page, pageSize, nil, date)
	}

	if err != nil {
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.Success(c, gin.H{
		"list":     list,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// ListSchedules 查询某科室某天的号源
func (h *RegistrationHandler) ListSchedules(c *gin.Context) {
	deptID := c.Query("deptId")
	date := c.Query("date")

	schedules, err := h.svc.ListSchedules(deptID, date)
	if err != nil {
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.Success(c, schedules)
}

// Register 挂号
func (h *RegistrationHandler) Register(c *gin.Context) {
	var req RegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	reg, err := h.svc.Register(req.PatientID, req.PatientName, req.ScheduleID, time.Now().Format("2006-01-02"))
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			response.FailWithMsg(c, appErr.Code, appErr.Message)
			return
		}
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.Success(c, reg)
}

// CancelRegistration 取消挂号
func (h *RegistrationHandler) CancelRegistration(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	if err := h.svc.Cancel(id); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			response.FailWithMsg(c, appErr.Code, appErr.Message)
			return
		}
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.Success(c, nil)
}

// SignIn 签到
func (h *RegistrationHandler) SignIn(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	if err := h.svc.SignIn(id); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			response.FailWithMsg(c, appErr.Code, appErr.Message)
			return
		}
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.Success(c, nil)
}

// GetQueueStatus 查询排队状态
func (h *RegistrationHandler) GetQueueStatus(c *gin.Context) {
	registrationID := c.Query("registrationId")
	if registrationID == "" {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	rank, err := h.svc.GetQueueStatus(registrationID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			response.FailWithMsg(c, appErr.Code, appErr.Message)
			return
		}
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.Success(c, gin.H{"registrationId": registrationID, "rank": rank})
}
