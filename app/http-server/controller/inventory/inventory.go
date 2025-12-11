package inventory

import (
	"belajarGo2/app/http-server/common"
	"belajarGo2/service/inventory"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
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

func (c *Controller) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req InventoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.logger.Error("inventory.create Error", slog.Any("error", err))
		// http.Error(w, err.Error(), http.StatusBadRequest)
		common.ErrorInvalidJJSON(w)
		return
	}

	if err := validator.New().Struct(req); err != nil {
		c.logger.Error("inventory.Create Validation error", slog.Any("error", err))

		// http.Error(w, err.Error(), http.StatusBadRequest)
		common.ErrorValidation(w, err)
		return
	}

	if err := c.inventorySvc.Create(inventory.Inventory{
		Code:        req.Code,
		Name:        req.Name,
		Stock:       req.Stock,
		Description: req.Description,
		Status:      req.Status,
	}); err != nil {
		c.logger.Error("inventory.Create Error", slog.Any("error", err))

		// http.Error(w, err.Error(), http.StatusInternalServerError)
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			// w.Header().Set("Content-Type", "application/json")
			// w.WriteHeader(http.StatusUnprocessableEntity)
			// _ = json.NewEncoder(w).Encode(map[string]interface{}{"message": http.StatusText(http.StatusUnprocessableEntity), "data": map[string]interface{}{}})
			common.ErrorDataConflict(w)
			return
		}

		// w.Header().Set("Content-Type", "application/json")
		// w.WriteHeader(http.StatusInternalServerError)
		// _ = json.NewEncoder(w).Encode(map[string]interface{}{"message": http.StatusText(http.StatusInternalServerError), "data": map[string]interface{}{}})
		common.ErrorInternal(w)
		return
	}

	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusCreated)
	// _ = json.NewEncoder(w).Encode(map[string]interface{}{"message": "inventory created successfully", "data": map[string]interface{}{"code": req.Code}})
	common.ValidResponse(w, http.StatusCreated, map[string]interface{}{"code": req.Code})
}

func (c *Controller) GetAll(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	pReq := r.FormValue("page")
	lReq := r.FormValue("limit")
	page, _ := strconv.Atoi(pReq)
	limit, _ := strconv.Atoi(lReq)

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	invs, err := c.inventorySvc.GetAll(page, limit)
	if err != nil {
		c.logger.Error("inventory.GetAll Error", slog.Any("error", err))
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		common.ErrorInternal(w)
		return
	}

	if len(invs) == 0 {
		// w.Header().Set("Content-Type", "application/json")
		// w.WriteHeader(http.StatusOK)
		// _ = json.NewEncoder(w).Encode(map[string]interface{}{"message": "no inventories found", "data": []inventory.Inventory{}})
		common.ValidResponse(w, http.StatusOK, []inventory.Inventory{})
		return
	}

	// w.Header().Set("Content-Type", "application/json")
	// _ = json.NewEncoder(w).Encode(map[string]interface{}{"message": "ok", "data": invs})
	common.ValidResponse(w, http.StatusOK, invs)
}

func (c *Controller) GetByCode(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	code := p.ByName("code")
	if code == "" {
		// http.Error(w, "code is required", http.StatusBadRequest)
		common.ErrorValidation(w, errors.New("code parameter is required"))
		return
	}

	inv, err := c.inventorySvc.GetByCode(code)
	if err != nil {
		c.logger.Error("inventory.GetByCode Error", slog.Any("error", err))

		// http.Error(w, err.Error(), http.StatusBadRequest)
		common.ErrorInternal(w)
		return
	}

	if inv.Code == "" {
		// w.Header().Set("Content-Type", "application/json")
		// w.WriteHeader(http.StatusNotFound)
		// _ = json.NewEncoder(w).Encode(map[string]interface{}{"message": "inventory not found", "data": map[string]interface{}{}})
		common.ErrorDataNotFound(w)
		return
	}

	// w.Header().Set("Content-Type", "application/json")
	// _ = json.NewEncoder(w).Encode(map[string]interface{}{"message": "ok", "data": inv})
	common.ValidResponse(w, http.StatusOK, inv)
}

func (c *Controller) Update(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// code := p.ByName("code")
	// if code == "" {
	// 	http.Error(w, "code is required", http.StatusBadRequest)
	// 	return
	// }
	var req InventoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.ErrorInvalidJJSON(w)
		return
	}
	req.Code = p.ByName("code")

	if err := validator.New().Struct(req); err != nil {
		c.logger.Error("inventory.Update validation error", slog.Any("error", err))

		common.ErrorValidation(w, err)
		return
	}

	if err := c.inventorySvc.Update(inventory.Inventory{
		Code:        req.Code,
		Name:        req.Name,
		Stock:       req.Stock,
		Description: req.Description,
		Status:      req.Status,
	}); err != nil {
		c.logger.Error("inventory.Update error", slog.Any("error", err))

		if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			common.ErrorDataConflict(w)
			return
		}

		common.ErrorInternal(w)
		return
	}

	common.ValidResponse(w, http.StatusOK, nil)
}

func (c *Controller) Delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	code := p.ByName("code")
	if code == "" {
		// http.Error(w, "code is required", http.StatusBadRequest)
		common.ErrorValidation(w, errors.New("code parameter is required"))
		return
	}

	if err := c.inventorySvc.Delete(code); err != nil {
		c.logger.Error("inventory.Delete Error", slog.Any("error", err))

		common.ErrorInternal(w)
		return
	}

	common.ValidResponse(w, http.StatusOK, nil)
}
