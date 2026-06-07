package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"his-go/internal/clinic/model"
	"his-go/internal/clinic/service"
	"his-go/pkg/errors"
	"his-go/pkg/response"
)

// ClinicHandler 门诊诊疗接口处理器
type ClinicHandler struct {
	svc *service.ClinicService
}

// NewClinicHandler 创建门诊诊疗接口处理器
func NewClinicHandler(svc *service.ClinicService) *ClinicHandler {
	return &ClinicHandler{svc: svc}
}

// CreateClinicRecord 创建门诊诊疗记录
func (h *ClinicHandler) CreateClinicRecord(c *gin.Context) {
	var record model.ClinicRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	if err := h.svc.CreateRecord(&record); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			response.FailWithMsg(c, appErr.Code, appErr.Message)
			return
		}
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.Success(c, record)
}

// GetClinicRecord 获取门诊诊疗记录详情
func (h *ClinicHandler) GetClinicRecord(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	record, err := h.svc.GetByID(id)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			response.FailWithMsg(c, appErr.Code, appErr.Message)
			return
		}
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.Success(c, record)
}

// ListClinicRecords 查询诊疗记录列表
func (h *ClinicHandler) ListClinicRecords(c *gin.Context) {
	patientID := c.Query("patientId")
	page, _ := strconv.Atoi(c.DefaultQuery(response.QueryKeyPage, response.DefaultPage))
	pageSize, _ := strconv.Atoi(c.DefaultQuery(response.QueryKeyPageSize, response.DefaultPageSize))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	list, total, err := h.svc.ListByPatient(patientID, page, pageSize)
	if err != nil {
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.SuccessPage(c, list, total, page, pageSize)
}

// CreateExaminationRequest 创建检查申请单
func (h *ClinicHandler) CreateExaminationRequest(c *gin.Context) {
	var req model.ExaminationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	if err := h.svc.CreateExamRequest(&req); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			response.FailWithMsg(c, appErr.Code, appErr.Message)
			return
		}
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.Success(c, req)
}
