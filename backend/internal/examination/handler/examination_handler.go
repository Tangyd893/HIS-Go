package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"his-go/internal/examination/model"
	"his-go/internal/examination/service"
	"his-go/pkg/errors"
	"his-go/pkg/response"
)

// ExaminationHandler 检查检验接口处理器
type ExaminationHandler struct {
	svc *service.ExaminationService
}

// NewExaminationHandler 创建检查检验接口处理器
func NewExaminationHandler(svc *service.ExaminationService) *ExaminationHandler {
	return &ExaminationHandler{svc: svc}
}

// ReviewRequest 审核请求
type ReviewRequest struct {
	ReportID   string `json:"report_id" binding:"required"`
	ReviewerID string `json:"reviewer_id" binding:"required"`
	Approved   bool   `json:"approved"`
	Comment    string `json:"comment"`
}

// CreateReport 创建检查报告
func (h *ExaminationHandler) CreateReport(c *gin.Context) {
	var report model.ExaminationReport
	if err := c.ShouldBindJSON(&report); err != nil {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	if err := h.svc.CreateReport(&report); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			response.FailWithMsg(c, appErr.Code, appErr.Message)
			return
		}
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.Success(c, report)
}

// GetReport 获取检查报告详情
func (h *ExaminationHandler) GetReport(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	report, err := h.svc.GetByID(id)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			response.FailWithMsg(c, appErr.Code, appErr.Message)
			return
		}
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.Success(c, report)
}

// ListReports 查询检查报告列表
func (h *ExaminationHandler) ListReports(c *gin.Context) {
	patientID := c.Query("patientId")
	status, _ := strconv.Atoi(c.DefaultQuery("status", "-1"))
	page, _ := strconv.Atoi(c.DefaultQuery(response.QueryKeyPage, response.DefaultPage))
	pageSize, _ := strconv.Atoi(c.DefaultQuery(response.QueryKeyPageSize, response.DefaultPageSize))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	list, total, err := h.svc.ListByPatient(patientID, status, page, pageSize)
	if err != nil {
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.SuccessPage(c, list, total, page, pageSize)
}

// ReviewReport 审核检查报告
func (h *ExaminationHandler) ReviewReport(c *gin.Context) {
	var req ReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	if err := h.svc.Review(req.ReportID, req.ReviewerID, req.Approved, req.Comment); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			response.FailWithMsg(c, appErr.Code, appErr.Message)
			return
		}
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.Success(c, nil)
}
