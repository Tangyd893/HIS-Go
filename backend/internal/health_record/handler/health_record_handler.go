package handler

import (
	"github.com/gin-gonic/gin"

	"his-go/internal/health_record/model"
	"his-go/internal/health_record/service"
	"his-go/pkg/errors"
	"his-go/pkg/response"
)

// HealthRecordHandler 健康档案接口处理器
type HealthRecordHandler struct {
	svc *service.HealthRecordService
}

// NewHealthRecordHandler 创建健康档案接口处理器
func NewHealthRecordHandler(svc *service.HealthRecordService) *HealthRecordHandler {
	return &HealthRecordHandler{svc: svc}
}

// GetSummary 获取档案摘要
func (h *HealthRecordHandler) GetSummary(c *gin.Context) {
	patientID := c.Query("patientId")
	if patientID == "" {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	summary, err := h.svc.GetSummary(patientID)
	if err != nil {
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.Success(c, summary)
}

// GetTimeline 获取时间轴
func (h *HealthRecordHandler) GetTimeline(c *gin.Context) {
	patientID := c.Query("patientId")
	if patientID == "" {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	events, err := h.svc.GetTimeline(patientID)
	if err != nil {
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.Success(c, events)
}

// GrantAuthorization 授权查看档案
func (h *HealthRecordHandler) GrantAuthorization(c *gin.Context) {
	var auth model.RecordAuthorization
	if err := c.ShouldBindJSON(&auth); err != nil {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	if err := h.svc.GrantAuthorization(&auth); err != nil {
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.Success(c, auth)
}

// RevokeAuthorization 撤销授权
func (h *HealthRecordHandler) RevokeAuthorization(c *gin.Context) {
	patientID := c.Query("patientId")
	doctorID := c.Query("doctorId")
	if patientID == "" || doctorID == "" {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	if err := h.svc.RevokeAuthorization(patientID, doctorID); err != nil {
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.Success(c, nil)
}
