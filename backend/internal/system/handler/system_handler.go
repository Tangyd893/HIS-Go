package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"his-go/internal/system/model"
	"his-go/internal/system/service"
	"his-go/pkg/errors"
	"his-go/pkg/response"
)

// SystemHandler 系统管理接口处理器
type SystemHandler struct {
	svc *service.SystemService
}

// NewSystemHandler 创建系统管理接口处理器
func NewSystemHandler(svc *service.SystemService) *SystemHandler {
	return &SystemHandler{svc: svc}
}

// ListDictTypes 查询字典类型列表
func (h *SystemHandler) ListDictTypes(c *gin.Context) {
	list, err := h.svc.ListDictTypes()
	if err != nil {
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.Success(c, list)
}

// ListDictItems 查询字典项列表
func (h *SystemHandler) ListDictItems(c *gin.Context) {
	dictType := c.Query("dictType")
	if dictType == "" {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	list, err := h.svc.ListDictItems(dictType)
	if err != nil {
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.Success(c, list)
}

// CreateDictItem 创建字典项
func (h *SystemHandler) CreateDictItem(c *gin.Context) {
	var item model.DictItem
	if err := c.ShouldBindJSON(&item); err != nil {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	if err := h.svc.CreateDictItem(&item); err != nil {
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.Success(c, item)
}

// ListParams 查询系统参数列表
func (h *SystemHandler) ListParams(c *gin.Context) {
	list, err := h.svc.ListParams()
	if err != nil {
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.Success(c, list)
}

// UpdateParam 更新系统参数
func (h *SystemHandler) UpdateParam(c *gin.Context) {
	var param model.SystemParam
	if err := c.ShouldBindJSON(&param); err != nil {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	if err := h.svc.UpdateParam(&param); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			response.FailWithMsg(c, appErr.Code, appErr.Message)
			return
		}
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.Success(c, nil)
}

// ListOperationLogs 查询操作日志
func (h *SystemHandler) ListOperationLogs(c *gin.Context) {
	userID := c.Query("userId")
	module := c.Query("module")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	list, total, err := h.svc.ListOperationLogs(userID, module, page, pageSize)
	if err != nil {
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.SuccessPage(c, list, total, page, pageSize)
}
