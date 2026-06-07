package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"his-go/internal/billing/model"
	"his-go/internal/billing/service"
	"his-go/pkg/errors"
	"his-go/pkg/response"
)

// BillingHandler 收费结算HTTP处理器
type BillingHandler struct {
	svc *service.BillingService
}

// NewBillingHandler 创建收费结算HTTP处理器
func NewBillingHandler(svc *service.BillingService) *BillingHandler {
	return &BillingHandler{svc: svc}
}

// createBillReq 创建账单请求
type createBillReq struct {
	PatientID      string             `json:"patient_id" binding:"required"`
	RegistrationID string             `json:"registration_id"`
	BillNo         string             `json:"bill_no" binding:"required"`
	Details        []model.BillDetail `json:"details"`
}

// CreateBill 创建账单
func (h *BillingHandler) CreateBill(c *gin.Context) {
	var req createBillReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	var totalAmount float64
	for _, d := range req.Details {
		totalAmount += float64(d.Quantity) * d.UnitPrice
	}

	bill := &model.Bill{
		PatientID:      req.PatientID,
		RegistrationID: req.RegistrationID,
		BillNo:         req.BillNo,
		TotalAmount:    totalAmount,
		Status:         0,
	}
	if err := h.svc.CreateBill(bill, req.Details); err != nil {
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.Success(c, bill)
}

// GetBill 获取账单详情
func (h *BillingHandler) GetBill(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	bill, err := h.svc.GetBill(id)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			response.Fail(c, appErr.Code)
		} else {
			response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		}
		return
	}

	details, _ := h.svc.GetBillDetails(id)
	response.Success(c, gin.H{"bill": bill, "details": details})
}

// PayBill 支付账单
func (h *BillingHandler) PayBill(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		PayMethod int8 `json:"pay_method" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	if err := h.svc.Pay(id, req.PayMethod); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			response.Fail(c, appErr.Code)
		} else {
			response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		}
		return
	}
	response.Success(c, gin.H{"message": "支付成功"})
}

// RefundBill 退款
func (h *BillingHandler) RefundBill(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	if err := h.svc.Refund(id); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			response.Fail(c, appErr.Code)
		} else {
			response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		}
		return
	}
	response.Success(c, gin.H{"message": "退款成功"})
}

// ListBills 分页查询账单
func (h *BillingHandler) ListBills(c *gin.Context) {
	patientID := c.Query("patient_id")
	if patientID == "" {
		response.Fail(c, errors.CodeParamInvalid)
		return
	}

	status, _ := strconv.Atoi(c.DefaultQuery("status", "-1"))
	page, _ := strconv.Atoi(c.DefaultQuery(response.QueryKeyPage, response.DefaultPage))
	pageSize, _ := strconv.Atoi(c.DefaultQuery(response.QueryKeyPageSize, response.DefaultPageSize))

	bills, total, err := h.svc.ListBills(patientID, status, page, pageSize)
	if err != nil {
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}
	response.SuccessPage(c, bills, total, page, pageSize)
}
