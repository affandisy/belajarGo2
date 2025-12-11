package inventory

import (
	"belajarGo2/service/inventory"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Controller struct {
	logger       *slog.Logger
	inventorySvc inventory.Service
}

func NewController(logger *slog.Logger, s inventory.Service) *Controller {
	return &Controller{
		logger:       logger,
		inventorySvc: s,
	}
}

type InventoryRequest struct {
	Code        string `json:"code" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Stock       int    `json:"stock"`
	Description string `json:"description"`
	Status      string `json:"status" validate:"required,oneof=active broken"`
}

func (ctrl *Controller) Create(c echo.Context) error {
	var req InventoryRequest
	if err := c.Bind(&req); err != nil {
		ctrl.logger.Error("inventory.Create Bind Error", slog.Any("error", err))
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	if err := validator.New().Struct(req); err != nil {
		ctrl.logger.Error("inventory.Create Validation Error", slog.Any("error", err))
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Validation error"})
	}

	if err := ctrl.inventorySvc.Create(inventory.Inventory{
		Code:        req.Code,
		Name:        req.Name,
		Stock:       req.Stock,
		Description: req.Description,
		Status:      req.Status,
	}); err != nil {
		ctrl.logger.Error("inventory.Create Service Error", slog.Any("error", err))

		if strings.Contains(err.Error(), "duplicate key") {
			return c.JSON(http.StatusConflict, map[string]string{"message": "Data conflict"})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{"message": "OK", "data": "data"})
}

func (ctrl *Controller) GetAll(c echo.Context) error {
	pReq := c.QueryParam("page")
	lReq := c.QueryParam("limit")
	page, _ := strconv.Atoi(pReq)
	limit, _ := strconv.Atoi(lReq)

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	invs, err := ctrl.inventorySvc.GetAll(page, limit)
	if err != nil {
		ctrl.logger.Error("inventory.GetAll Service Error", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
	}

	if len(invs) == 0 {
		return c.JSON(http.StatusOK, []inventory.Inventory{})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "OK", "data": invs})
}

func (ctrl *Controller) GetByCode(c echo.Context) error {
	code := c.Param("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Code parameter is required"})
	}

	inv, err := ctrl.inventorySvc.GetByCode(code)
	if err != nil {
		ctrl.logger.Error("inventory.GetByCode Service Error", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
	}

	if inv.Code == "" {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Data not found"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "OK", "data": inv})
}

func (ctrl *Controller) Update(c echo.Context) error {
	var req InventoryRequest
	if err := c.Bind(&req); err != nil {
		ctrl.logger.Error("inventory.Update Bind Error", slog.Any("error", err))
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}
	req.Code = c.Param("code")

	if err := validator.New().Struct(req); err != nil {
		ctrl.logger.Error("inventory.Update Validation Error", slog.Any("error", err))
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Validation error"})
	}

	if err := ctrl.inventorySvc.Update(inventory.Inventory{
		Code:        req.Code,
		Name:        req.Name,
		Stock:       req.Stock,
		Description: req.Description,
		Status:      req.Status,
	}); err != nil {
		ctrl.logger.Error("inventory.Update Service Error", slog.Any("error", err))

		if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			return c.JSON(http.StatusConflict, map[string]string{"message": "Data conflict"})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "OK", "data": map[string]string{}})
}

func (ctrl *Controller) Delete(c echo.Context) error {
	code := c.Param("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Code parameter is required"})
	}

	if err := ctrl.inventorySvc.Delete(code); err != nil {
		ctrl.logger.Error("inventory.Delete Service Error", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "OK", "data": map[string]string{}})
}
