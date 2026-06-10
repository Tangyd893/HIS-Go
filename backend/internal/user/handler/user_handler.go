package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"his-go/internal/user/model"
	"his-go/internal/user/service"
	apperrors "his-go/pkg/errors"
	"his-go/pkg/response"
	"his-go/pkg/security/auth"
)

// UserHandler 用户管理接口处理器
type UserHandler struct {
	svc *service.UserService
}

// NewUserHandler 创建用户管理接口处理器
func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// ---- 患者接口 ----

// ListPatients 患者列表查询接口
func (h *UserHandler) ListPatients(c *gin.Context) {
	name := c.Query("name")
	phone := c.Query("phone")
	page, _ := strconv.Atoi(c.DefaultQuery(response.QueryKeyPage, response.DefaultPage))
	pageSize, _ := strconv.Atoi(c.DefaultQuery(response.QueryKeyPageSize, response.DefaultPageSize))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	list, total, err := h.svc.ListPatients(name, phone, page, pageSize)
	if err != nil {
		response.Fail(c, apperrors.CodeInternalError)
		return
	}
	response.SuccessPage(c, list, total, page, pageSize)
}

// GetMyPatient 获取当前登录用户对应的患者档案
func (h *UserHandler) GetMyPatient(c *gin.Context) {
	userCtx := auth.GetUserContext(c)
	if userCtx == nil || userCtx.UserID == "" {
		response.Fail(c, apperrors.CodeUnauthorized)
		return
	}

	patient, err := h.svc.GetPatientByUserID(userCtx.UserID)
	if err != nil {
		response.Fail(c, apperrors.CodeNotFound)
		return
	}
	response.Success(c, patient)
}

// GetPatient 患者详情查询接口
func (h *UserHandler) GetPatient(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Fail(c, apperrors.CodeParamInvalid)
		return
	}

	patient, err := h.svc.GetPatient(id)
	if err != nil {
		response.Fail(c, apperrors.CodeNotFound)
		return
	}
	response.Success(c, patient)
}

// CreatePatient 创建患者接口
func (h *UserHandler) CreatePatient(c *gin.Context) {
	var patient model.Patient
	if err := c.ShouldBindJSON(&patient); err != nil {
		response.Fail(c, apperrors.CodeParamInvalid)
		return
	}

	if err := h.svc.CreatePatient(&patient); err != nil {
		response.Fail(c, apperrors.CodeInternalError)
		return
	}
	response.Success(c, patient)
}

// UpdatePatient 更新患者接口
func (h *UserHandler) UpdatePatient(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Fail(c, apperrors.CodeParamInvalid)
		return
	}

	var patient model.Patient
	if err := c.ShouldBindJSON(&patient); err != nil {
		response.Fail(c, apperrors.CodeParamInvalid)
		return
	}
	patient.ID = id

	if err := h.svc.UpdatePatient(&patient); err != nil {
		response.Fail(c, apperrors.CodeInternalError)
		return
	}
	response.Success(c, patient)
}

// DeletePatient 删除患者接口
func (h *UserHandler) DeletePatient(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Fail(c, apperrors.CodeParamInvalid)
		return
	}

	if err := h.svc.DeletePatient(id); err != nil {
		response.Fail(c, apperrors.CodeInternalError)
		return
	}
	response.Success(c, nil)
}

// ---- 科室接口 ----

// ListDepartments 科室列表查询接口
func (h *UserHandler) ListDepartments(c *gin.Context) {
	depts, err := h.svc.ListDepartments()
	if err != nil {
		response.Fail(c, apperrors.CodeInternalError)
		return
	}
	response.Success(c, depts)
}

// ---- 员工接口 ----

// ListEmployees 员工列表查询接口
func (h *UserHandler) ListEmployees(c *gin.Context) {
	deptID := c.Query("deptId")
	name := c.Query("name")
	page, _ := strconv.Atoi(c.DefaultQuery(response.QueryKeyPage, response.DefaultPage))
	pageSize, _ := strconv.Atoi(c.DefaultQuery(response.QueryKeyPageSize, response.DefaultPageSize))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	list, total, err := h.svc.ListEmployees(deptID, name, page, pageSize)
	if err != nil {
		response.Fail(c, apperrors.CodeInternalError)
		return
	}
	response.SuccessPage(c, list, total, page, pageSize)
}

// GetEmployee 员工详情查询接口
func (h *UserHandler) GetEmployee(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Fail(c, apperrors.CodeParamInvalid)
		return
	}

	emp, err := h.svc.GetEmployee(id)
	if err != nil {
		response.Fail(c, apperrors.CodeNotFound)
		return
	}
	response.Success(c, emp)
}
