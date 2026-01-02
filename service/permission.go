package service

import (
	"alfdwirhmn/inventory/model"
	"errors"
)

type PermissionService interface {
	// USER
	CanManageUsers(role string) bool
	CanCreateUser(currentRole, targetRole string) error
	CanUpdateUser(currentUser, targetUser *model.User, newRole string) error
	CanDeleteUser(currentUser, targetUser *model.User) error

	// MASTER DATA
	CanReadMasterData(role string) bool
	CanCreateMasterData(role string) bool
	CanUpdateMasterData(role string) bool
	CanDeleteMasterData(role string) bool

	// STOCK
	CanUpdateStock(role string) bool
	CanCheckMinStock(role string) bool

	// SALE
	CanCreateSale(role string) bool
	CanViewSale(role string) bool

	// REPORT
	CanAccessReports(role string) bool
}

type permissionService struct{}

func NewPermissionService() PermissionService {
	return &permissionService{}
}

func (s *permissionService) CanManageUsers(role string) bool {
	// hanya super admin dan admin
	return role == "super_admin" || role == "admin"
}

func (s *permissionService) CanCreateUser(currentRole, targetRole string) error {
	switch currentRole {

	// super_admin bebas membuat role apapun
	case "super_admin":
		return nil

	case "admin":
		// admin tidak boleh membuat super_admin
		if targetRole == "super_admin" {
			return errors.New("admin cannot create super_admin")
		}
		// super admin bisa buat admin dan staff
		if targetRole == "admin" || targetRole == "staff" {
			return nil
		}
		return errors.New("invalid target role")

	default:
		return errors.New("insufficient permission")
	}
}

func (s *permissionService) CanUpdateUser(currentUser, targetUser *model.User, newRole string) error {
	// user tidak boleh mengganti role sendiri
	if currentUser.ID == targetUser.ID && newRole != "" && newRole != currentUser.Role {
		return errors.New("cannot change your own role")
	}

	// super_admin bebas update siapa pun
	if currentUser.Role == "super_admin" {
		return nil
	}

	// admin tidak boleh menyentuh super_admin
	if currentUser.Role == "admin" {
		if targetUser.Role == "super_admin" || newRole == "super_admin" {
			return errors.New("admin cannot modify super_admin")
		}
		return nil
	}

	return errors.New("insufficient permission")
}

func (s *permissionService) CanDeleteUser(currentUser, targetUser *model.User) error {
	// user tidak boleh menghapus akunnya sendiri
	if currentUser.ID == targetUser.ID {
		return errors.New("cannot delete your own account")
	}

	if currentUser.Role == "super_admin" {
		return nil
	}

	// admin tidak boleh hapus super_admin
	if currentUser.Role == "admin" {
		if targetUser.Role == "super_admin" {
			return errors.New("admin cannot delete super_admin")
		}
		return nil
	}

	return errors.New("insufficient permission")
}

// master data, hanya admin dan super admin yang bisa crud
func (s *permissionService) CanReadMasterData(role string) bool {
	return role == "super_admin" || role == "admin" || role == "staff"
}

func (s *permissionService) CanCreateMasterData(role string) bool {
	return role == "super_admin" || role == "admin"
}

func (s *permissionService) CanUpdateMasterData(role string) bool {
	return role == "super_admin" || role == "admin"
}

func (s *permissionService) CanDeleteMasterData(role string) bool {
	return role == "super_admin" || role == "admin"
}

// stock, semua role bisa update stock
func (s *permissionService) CanUpdateStock(role string) bool {
	return role == "super_admin" || role == "admin" || role == "staff"
}

func (s *permissionService) CanCheckMinStock(role string) bool {
	return role == "super_admin" || role == "admin" || role == "staff"
}

// sale, semua role bisa input trasaksi
func (s *permissionService) CanCreateSale(role string) bool {
	return role == "super_admin" || role == "admin" || role == "staff"
}

func (s *permissionService) CanViewSale(role string) bool {
	return role == "super_admin" || role == "admin" || role == "staff"
}

// report, hanya admin dan super yang bisa akses report
func (s *permissionService) CanAccessReports(role string) bool {
	return role == "super_admin" || role == "admin"
}
