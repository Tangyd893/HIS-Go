package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"his-go/internal/emr/model"
	"his-go/internal/emr/service"
	"his-go/pkg/errors"
	"his-go/pkg/response"
)

// EMRHandler 电子病历HTTP处理器
type EMRHandler struct {
	svc *service.EMRService
}

// NewEMRHandler 创建电子病历HTTP处理器
func NewEMRHandler(svc *service.EMRService) *EMRHandler {
	return &EMRHandler{svc: svc}
}

// CreateRecord 创建病历
func (h *EMRHandler) CreateRecord(c *gin.Context) {
	var record model.MedicalRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	if err := h.svc.CreateRecord(&record); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			response.Fail(c, appErr.Code)
		} else {
			response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		}
		return
	}
	response.Success(c, record)
}

// GetRecord 获取病历详情
func (h *EMRHandler) GetRecord(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	record, err := h.svc.GetRecord(id)
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

// ListRecords 分页查询患者病历
func (h *EMRHandler) ListRecords(c *gin.Context) {
	patientID := c.Query("patient_id")
	if patientID == "" {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery(response.QueryKeyPage, response.DefaultPage))
	pageSize, _ := strconv.Atoi(c.DefaultQuery(response.QueryKeyPageSize, response.DefaultPageSize))

	records, total, err := h.svc.ListRecords(patientID, page, pageSize)
	if err != nil {
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.SuccessPage(c, records, total, page, pageSize)
}

// QualityControl 病历质控
func (h *EMRHandler) QualityControl(c *gin.Context) {
	recordID := c.Param("id")

	var req struct {
		ReviewerID string `json:"reviewer_id" binding:"required"`
		Level      int    `json:"level" binding:"required"`
		Comment    string `json:"comment"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	if err := h.svc.QualityControl(recordID, req.ReviewerID, req.Level, req.Comment); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			response.Fail(c, appErr.Code)
		} else {
			response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		}
		return
	}
	response.Success(c, gin.H{"message": "质控完成"})
}

// ListTemplates 获取病历模板列表
func (h *EMRHandler) ListTemplates(c *gin.Context) {
	templates, err := h.svc.ListTemplates()
	if err != nil {
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.Success(c, templates)
}

// CDSSCheck 临床决策支持检查
func (h *EMRHandler) CDSSCheck(c *gin.Context) {
	var req struct {
		PatientID string `json:"patient_id" binding:"required"`
		DrugID    string `json:"drug_id" binding:"required"`
		Diagnosis string `json:"diagnosis" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	warnings, err := h.svc.CDSSCheck(req.PatientID, req.DrugID, req.Diagnosis)
	if err != nil {
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	if len(warnings) == 0 {
		warnings = []string{}
	}
	response.Success(c, gin.H{"warnings": warnings})
}
