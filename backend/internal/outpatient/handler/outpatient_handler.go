package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"his-go/internal/outpatient/model"
	"his-go/internal/outpatient/service"
	"his-go/pkg/errors"
	"his-go/pkg/response"
)

// OutpatientHandler 院外患者接口处理器
type OutpatientHandler struct {
	svc *service.OutpatientService
}

// NewOutpatientHandler 创建院外患者接口处理器
func NewOutpatientHandler(svc *service.OutpatientService) *OutpatientHandler {
	return &OutpatientHandler{svc: svc}
}

// CreateConsultation 创建在线问诊
func (h *OutpatientHandler) CreateConsultation(c *gin.Context) {
	var consultation model.Consultation
	if err := c.ShouldBindJSON(&consultation); err != nil {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	if err := h.svc.CreateConsultation(&consultation); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			response.FailWithMsg(c, appErr.Code, appErr.Message)
			return
		}
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.Success(c, consultation)
}

// GetConsultation 获取问诊详情
func (h *OutpatientHandler) GetConsultation(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	consultation, err := h.svc.GetConsultation(id)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			response.FailWithMsg(c, appErr.Code, appErr.Message)
			return
		}
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.Success(c, consultation)
}

// ListConsultations 查询问诊列表
func (h *OutpatientHandler) ListConsultations(c *gin.Context) {
	patientID := c.Query("patientId")
	doctorID := c.Query("doctorId")
	status, _ := strconv.Atoi(c.DefaultQuery("status", "-1"))
	page, _ := strconv.Atoi(c.DefaultQuery(response.QueryKeyPage, response.DefaultPage))
	pageSize, _ := strconv.Atoi(c.DefaultQuery(response.QueryKeyPageSize, response.DefaultPageSize))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	list, total, err := h.svc.ListConsultations(patientID, doctorID, status, page, pageSize)
	if err != nil {
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.SuccessPage(c, list, total, page, pageSize)
}

// SendMessage 发送问诊消息
func (h *OutpatientHandler) SendMessage(c *gin.Context) {
	var msg model.ConsultationMessage
	if err := c.ShouldBindJSON(&msg); err != nil {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	if err := h.svc.SendMessage(&msg); err != nil {
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.Success(c, msg)
}

// GetMessages 获取问诊消息
func (h *OutpatientHandler) GetMessages(c *gin.Context) {
	consultationID := c.Query("consultationId")
	if consultationID == "" {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	list, err := h.svc.GetMessages(consultationID)
	if err != nil {
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.Success(c, list)
}

// CreateChronicContract 创建慢病签约
func (h *OutpatientHandler) CreateChronicContract(c *gin.Context) {
	var contract model.ChronicContract
	if err := c.ShouldBindJSON(&contract); err != nil {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	if err := h.svc.CreateChronicContract(&contract); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			response.FailWithMsg(c, appErr.Code, appErr.Message)
			return
		}
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.Success(c, contract)
}

// ReportHealthData 上报健康数据
func (h *OutpatientHandler) ReportHealthData(c *gin.Context) {
	var data model.HealthData
	if err := c.ShouldBindJSON(&data); err != nil {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	if err := h.svc.ReportHealthData(&data); err != nil {
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.Success(c, data)
}

// ListHealthData 查询健康数据列表
func (h *OutpatientHandler) ListHealthData(c *gin.Context) {
	patientID := c.Query("patientId")
	if patientID == "" {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	list, err := h.svc.ListHealthData(patientID)
	if err != nil {
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.Success(c, list)
}
