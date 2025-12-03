package router

import (
	"belajarGo2/app/echo-server/controller/inventory"
	"belajarGo2/app/echo-server/controller/user"
	"belajarGo2/app/echo-server/middleware"
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterPath(e *echo.Echo, jwtSecret string, ctrlInv *inventory.Controller, ctrlUser *user.Controller) {
	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"meesage": "pong",
		})
	})

	jwtMiddleware := middleware.JWTMiddleware(jwtSecret)
	// userNAdmin := middleware.RBACMiddleware([]string{"user", "admin"})
	// adminOnly := middleware.RBACMiddleware([]string{"admin"})
	// superadminOnly := middleware.RBACMiddleware([]string{"superadmin"})

	// init ACL
	userNAdminAccess := middleware.ACLMiddleware(map[string]bool{
		"admin": true,
		"user":  true,
	})
	adminAccess := middleware.ACLMiddleware(map[string]bool{
		"admin": true,
	})
	superadminAccess := middleware.ACLMiddleware(map[string]bool{
		"superadmin": true,
	})

	// user endpoint
	userEndpoint := e.Group("/users")
	userEndpoint.POST("/register", ctrlUser.Register)
	userEndpoint.POST("/login", ctrlUser.Login)

	// inventory endpoint
	inventoryEndpoint := e.Group("/inventories", jwtMiddleware)
	// inventoryEndpoint.GET("", ctrlInv.GetAll, userNAdmin)
	// inventoryEndpoint.GET("/:code", ctrlInv.GetByCode, userNAdmin)
	// inventoryEndpoint.POST("", ctrlInv.Create, adminOnly)
	// inventoryEndpoint.PUT("/:code", ctrlInv.Update, adminOnly)
	// inventoryEndpoint.DELETE("/:code", ctrlInv.Delete, superadminOnly)
	inventoryEndpoint.GET("", ctrlInv.GetAll, userNAdminAccess)
	inventoryEndpoint.GET("/:code", ctrlInv.GetByCode, userNAdminAccess)
	inventoryEndpoint.POST("", ctrlInv.Create, adminAccess)
	inventoryEndpoint.PUT("/:code", ctrlInv.Update, adminAccess)
	inventoryEndpoint.DELETE("/:code", ctrlInv.Delete, superadminAccess)
}
