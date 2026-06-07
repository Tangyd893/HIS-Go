package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"his-go/internal/prescription/model"
	"his-go/internal/prescription/service"
	"his-go/pkg/errors"
	"his-go/pkg/response"
)

// PrescriptionHandler 处方接口处理器
type PrescriptionHandler struct {
	svc *service.PrescriptionService
}

// NewPrescriptionHandler 创建处方接口处理器
func NewPrescriptionHandler(svc *service.PrescriptionService) *PrescriptionHandler {
	return &PrescriptionHandler{svc: svc}
}

// CreatePrescriptionRequest 创建处方请求
type CreatePrescriptionRequest struct {
	Prescription model.Prescription         `json:"prescription"`
	Details      []model.PrescriptionDetail `json:"details"`
}

// ReviewRequest 审核请求
type ReviewRequest struct {
	ID       string `json:"id" binding:"required"`
	Approved bool   `json:"approved"`
	Comment  string `json:"comment"`
}

// CreatePrescription 创建处方（含明细）
func (h *PrescriptionHandler) CreatePrescription(c *gin.Context) {
	var req CreatePrescriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	if err := h.svc.Create(&req.Prescription, req.Details); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			response.FailWithMsg(c, appErr.Code, appErr.Message)
			return
		}
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.Success(c, req.Prescription)
}

// GetPrescription 获取处方详情
func (h *PrescriptionHandler) GetPrescription(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	p, err := h.svc.GetByID(id)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			response.FailWithMsg(c, appErr.Code, appErr.Message)
			return
		}
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.Success(c, p)
}

// ListPrescriptions 查询处方列表
func (h *PrescriptionHandler) ListPrescriptions(c *gin.Context) {
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

// ReviewPrescription 审核处方
func (h *PrescriptionHandler) ReviewPrescription(c *gin.Context) {
	var req ReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	if err := h.svc.Review(req.ID, req.Approved, req.Comment); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			response.FailWithMsg(c, appErr.Code, appErr.Message)
			return
		}
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.Success(c, nil)
}

// CancelPrescription 取消处方
func (h *PrescriptionHandler) CancelPrescription(c *gin.Context) {
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
