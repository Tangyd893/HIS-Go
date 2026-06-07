package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"his-go/internal/followup/model"
	"his-go/internal/followup/service"
	"his-go/pkg/errors"
	"his-go/pkg/response"
)

// FollowupHandler 随访接口处理器
type FollowupHandler struct {
	svc *service.FollowupService
}

// NewFollowupHandler 创建随访接口处理器
func NewFollowupHandler(svc *service.FollowupService) *FollowupHandler {
	return &FollowupHandler{svc: svc}
}

// ExecuteTaskRequest 执行任务请求
type ExecuteTaskRequest struct {
	TaskID string `json:"task_id" binding:"required"`
	Result string `json:"result"`
}

// CreatePlan 创建随访计划
func (h *FollowupHandler) CreatePlan(c *gin.Context) {
	var plan model.FollowupPlan
	if err := c.ShouldBindJSON(&plan); err != nil {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	if err := h.svc.CreatePlan(&plan); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			response.FailWithMsg(c, appErr.Code, appErr.Message)
			return
		}
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.Success(c, plan)
}

// GetPlan 获取随访计划详情
func (h *FollowupHandler) GetPlan(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	plan, err := h.svc.GetPlan(id)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			response.FailWithMsg(c, appErr.Code, appErr.Message)
			return
		}
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.Success(c, plan)
}

// ListPlans 查询随访计划列表
func (h *FollowupHandler) ListPlans(c *gin.Context) {
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

	list, total, err := h.svc.ListPlans(patientID, status, page, pageSize)
	if err != nil {
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.SuccessPage(c, list, total, page, pageSize)
}

// ExecuteTask 执行随访任务
func (h *FollowupHandler) ExecuteTask(c *gin.Context) {
	var req ExecuteTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	if err := h.svc.ExecuteTask(req.TaskID, req.Result); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			response.FailWithMsg(c, appErr.Code, appErr.Message)
			return
		}
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.Success(c, nil)
}

// SubmitSurvey 提交满意度调查
func (h *FollowupHandler) SubmitSurvey(c *gin.Context) {
	var survey model.SatisfactionSurvey
	if err := c.ShouldBindJSON(&survey); err != nil {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	if err := h.svc.SubmitSurvey(&survey); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			response.FailWithMsg(c, appErr.Code, appErr.Message)
			return
		}
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.Success(c, survey)
}
