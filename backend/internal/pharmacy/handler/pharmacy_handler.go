package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"his-go/internal/pharmacy/service"
	"his-go/pkg/errors"
	"his-go/pkg/response"
)

// PharmacyHandler 药房管理HTTP处理器
type PharmacyHandler struct {
	svc *service.PharmacyService
}

// NewPharmacyHandler 创建药房管理HTTP处理器
func NewPharmacyHandler(svc *service.PharmacyService) *PharmacyHandler {
	return &PharmacyHandler{svc: svc}
}

// ListDrugs 分页查询药品列表
func (h *PharmacyHandler) ListDrugs(c *gin.Context) {
	name := c.Query("name")
	page, _ := strconv.Atoi(c.DefaultQuery(response.QueryKeyPage, response.DefaultPage))
	pageSize, _ := strconv.Atoi(c.DefaultQuery(response.QueryKeyPageSize, response.DefaultPageSize))

	drugs, total, err := h.svc.ListDrugs(name, page, pageSize)
	if err != nil {
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.SuccessPage(c, drugs, total, page, pageSize)
}

// GetDrug 获取药品详情
func (h *PharmacyHandler) GetDrug(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	drug, err := h.svc.GetDrug(id)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			response.Fail(c, appErr.Code)
		} else {
			response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		}
		return
	}
	response.Success(c, drug)
}

// AddStock 增加药品库存
func (h *PharmacyHandler) AddStock(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Quantity int `json:"quantity" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	if err := h.svc.AddStock(id, req.Quantity); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			response.Fail(c, appErr.Code)
		} else {
			response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		}
		return
	}
	response.Success(c, gin.H{"message": "入库成功"})
}

// DispenseDrug 发药
func (h *PharmacyHandler) DispenseDrug(c *gin.Context) {
	var req struct {
		PrescriptionID string `json:"prescription_id" binding:"required"`
		DrugID         string `json:"drug_id" binding:"required"`
		Quantity       int    `json:"quantity" binding:"required,min=1"`
		DispatcherID   string `json:"dispenser_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	if err := h.svc.DispenseDrug(req.PrescriptionID, req.DrugID, req.Quantity, req.DispatcherID); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			response.Fail(c, appErr.Code)
		} else {
			response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		}
		return
	}
	response.Success(c, gin.H{"message": "发药成功"})
}

// CheckExpiredDrugs 检查过期药品
func (h *PharmacyHandler) CheckExpiredDrugs(c *gin.Context) {
	h.svc.CheckAndAlertExpired()
	response.Success(c, gin.H{"message": "过期药品扫描已触发"})
}
