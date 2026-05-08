package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"his-go/internal/inpatient/model"
	"his-go/internal/inpatient/service"
	"his-go/pkg/errors"
	"his-go/pkg/response"
)

// InpatientHandler 住院管理HTTP处理器
type InpatientHandler struct {
	svc *service.InpatientService
}

// NewInpatientHandler 创建住院管理HTTP处理器
func NewInpatientHandler(svc *service.InpatientService) *InpatientHandler {
	return &InpatientHandler{svc: svc}
}

// AdmitPatient 入院登记
func (h *InpatientHandler) AdmitPatient(c *gin.Context) {
	var record model.InpatientRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	if err := h.svc.AdmitPatient(&record); err != nil {
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.Success(c, record)
}

// DischargePatient 出院
func (h *InpatientHandler) DischargePatient(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	if err := h.svc.DischargePatient(id); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			response.Fail(c, appErr.Code)
		} else {
			response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		}
		return
	}
	response.Success(c, gin.H{"message": "出院成功"})
}

// GetInpatient 获取住院记录详情
func (h *InpatientHandler) GetInpatient(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	record, err := h.svc.GetInpatient(id)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			response.Fail(c, appErr.Code)
		} else {
			response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		}
		return
	}
	response.Success(c, record)
}

// ListInpatients 分页查询住院列表
func (h *InpatientHandler) ListInpatients(c *gin.Context) {
	deptID := c.Query("dept_id")
	status, _ := strconv.Atoi(c.DefaultQuery("status", "-1"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	records, total, err := h.svc.ListInpatients(deptID, status, page, pageSize)
	if err != nil {
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.SuccessPage(c, records, total, page, pageSize)
}

// CreateMedicalOrder 创建医嘱
func (h *InpatientHandler) CreateMedicalOrder(c *gin.Context) {
	var order model.MedicalOrder
	if err := c.ShouldBindJSON(&order); err != nil {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	if err := h.svc.CreateMedicalOrder(&order); err != nil {
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.Success(c, order)
}

// CreateNursingRecord 创建护理记录
func (h *InpatientHandler) CreateNursingRecord(c *gin.Context) {
	var record model.NursingRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	if err := h.svc.CreateNursingRecord(&record); err != nil {
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.Success(c, record)
}
