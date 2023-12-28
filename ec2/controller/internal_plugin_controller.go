package controller

import (
	"github.com/gofiber/fiber/v2"
)

type InternalPluginController interface {
	Register() fiber.Handler
	Login() fiber.Handler
	User() fiber.Handler
	Logout() fiber.Handler
	UpdatePassword() fiber.Handler
	UpdateInfo() fiber.Handler
	Upload() fiber.Handler
	AllOrders() fiber.Handler
	AddOrder() fiber.Handler
	CreateOrder() fiber.Handler
	UpdateOrder() fiber.Handler
	GetOrder() fiber.Handler
	DeleteOrder() fiber.Handler
	Export() fiber.Handler
	CreateFile() fiber.Handler
	AllPermissions() fiber.Handler
	AllProducts() fiber.Handler
	CreateProduct() fiber.Handler
	GetProduct() fiber.Handler
	UpdateProduct() fiber.Handler
	DeleteProduct() fiber.Handler
	AllRoles() fiber.Handler
	CreateRole() fiber.Handler
	UpdateRole() fiber.Handler
	DeleteRole() fiber.Handler
	GetRole() fiber.Handler
	AllUsers() fiber.Handler
	CreateUser() fiber.Handler
	GetUser() fiber.Handler
	UpdateUser() fiber.Handler
	DeleteUser() fiber.Handler
	Chart() fiber.Handler
}

type InternalPluginControllerImpl struct {
}

// AddOrder implements InternalPluginController.
func (InternalPluginControllerImpl) AddOrder() fiber.Handler {
	panic("unimplemented")
}

// GetOrder implements InternalPluginController.
func (InternalPluginControllerImpl) GetOrder() fiber.Handler {
	panic("unimplemented")
}

// DeleteOrder implements InternalPluginController.
func (InternalPluginControllerImpl) DeleteOrder() fiber.Handler {
	panic("unimplemented")
}

// UpdateOrder implements InternalPluginController.
func (InternalPluginControllerImpl) UpdateOrder() fiber.Handler {
	panic("unimplemented")
}

// CreateOrder implements InternalPluginController.
func (InternalPluginControllerImpl) CreateOrder() fiber.Handler {
	panic("unimplemented")
}

// AllOrders implements InternalPluginController.
func (InternalPluginControllerImpl) AllOrders() fiber.Handler {
	panic("unimplemented")
}

// AllPermissions implements InternalPluginController.
func (InternalPluginControllerImpl) AllPermissions() fiber.Handler {
	panic("unimplemented")
}

// AllProducts implements InternalPluginController.
func (InternalPluginControllerImpl) AllProducts() fiber.Handler {
	panic("unimplemented")
}

// AllRoles implements InternalPluginController.
func (InternalPluginControllerImpl) AllRoles() fiber.Handler {
	panic("unimplemented")
}

// AllUsers implements InternalPluginController.
func (InternalPluginControllerImpl) AllUsers() fiber.Handler {
	panic("unimplemented")
}

// Chart implements InternalPluginController.
func (InternalPluginControllerImpl) Chart() fiber.Handler {
	panic("unimplemented")
}

// CreateFile implements InternalPluginController.
func (InternalPluginControllerImpl) CreateFile() fiber.Handler {
	panic("unimplemented")
}

// CreateProduct implements InternalPluginController.
func (InternalPluginControllerImpl) CreateProduct() fiber.Handler {
	panic("unimplemented")
}

// CreateRole implements InternalPluginController.
func (InternalPluginControllerImpl) CreateRole() fiber.Handler {
	panic("unimplemented")
}

// CreateUser implements InternalPluginController.
func (InternalPluginControllerImpl) CreateUser() fiber.Handler {
	panic("unimplemented")
}

// DeleteProduct implements InternalPluginController.
func (InternalPluginControllerImpl) DeleteProduct() fiber.Handler {
	panic("unimplemented")
}

// DeleteRole implements InternalPluginController.
func (InternalPluginControllerImpl) DeleteRole() fiber.Handler {
	panic("unimplemented")
}

// DeleteUser implements InternalPluginController.
func (InternalPluginControllerImpl) DeleteUser() fiber.Handler {
	panic("unimplemented")
}

// Export implements InternalPluginController.
func (InternalPluginControllerImpl) Export() fiber.Handler {
	panic("unimplemented")
}

// GetProduct implements InternalPluginController.
func (InternalPluginControllerImpl) GetProduct() fiber.Handler {
	panic("unimplemented")
}

// GetRole implements InternalPluginController.
func (InternalPluginControllerImpl) GetRole() fiber.Handler {
	panic("unimplemented")
}

// GetUser implements InternalPluginController.
func (InternalPluginControllerImpl) GetUser() fiber.Handler {
	panic("unimplemented")
}

// Login implements InternalPluginController.
func (InternalPluginControllerImpl) Login() fiber.Handler {
	panic("unimplemented")
}

// Logout implements InternalPluginController.
func (InternalPluginControllerImpl) Logout() fiber.Handler {
	panic("unimplemented")
}

// Register implements InternalPluginController.
func (InternalPluginControllerImpl) Register() fiber.Handler {
	panic("unimplemented")
}

// UpdateInfo implements InternalPluginController.
func (InternalPluginControllerImpl) UpdateInfo() fiber.Handler {
	panic("unimplemented")
}

// UpdatePassword implements InternalPluginController.
func (InternalPluginControllerImpl) UpdatePassword() fiber.Handler {
	panic("unimplemented")
}

// UpdateProduct implements InternalPluginController.
func (InternalPluginControllerImpl) UpdateProduct() fiber.Handler {
	panic("unimplemented")
}

// UpdateRole implements InternalPluginController.
func (InternalPluginControllerImpl) UpdateRole() fiber.Handler {
	panic("unimplemented")
}

// UpdateUser implements InternalPluginController.
func (InternalPluginControllerImpl) UpdateUser() fiber.Handler {
	panic("unimplemented")
}

// Upload implements InternalPluginController.
func (InternalPluginControllerImpl) Upload() fiber.Handler {
	panic("unimplemented")
}

// User implements InternalPluginController.
func (InternalPluginControllerImpl) User() fiber.Handler {
	panic("unimplemented")
}

func NewInternalPluginController() InternalPluginController {
	return InternalPluginControllerImpl{}
}
