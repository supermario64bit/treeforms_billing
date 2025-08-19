package services

import (
	"fmt"
	"net/http"
	"strconv"
	"treeforms_billing/application_types"
	"treeforms_billing/auth"
	"treeforms_billing/logger"

	"gorm.io/gorm"
)

type passwordService struct {
	db *gorm.DB
}

type PasswordService interface {
	Create(userID uint, plainPassword string) *application_types.ApplicationError
	ChangePassword(userID uint, currentPassword, newPlainPassword string) *application_types.ApplicationError
	ChangePasswordWithoutConfirmingCurrentPassword(userID uint, plainPassword string) *application_types.ApplicationError
}

func (svc *passwordService) Create(userID uint, plainPassword string) *application_types.ApplicationError {
	password, err := auth.NewPassword(userID, plainPassword)
	if err != nil {
		return application_types.NewApplicationError(false, http.StatusInternalServerError, "Password creation failed",
			fmt.Errorf("Unable to create password for the userid "+strconv.FormatUint(uint64(userID), 10)+". Message: "+err.Error()))
	}

	if err := svc.db.Create(password).Error; err != nil {
		fmt.Errorf("Unable to create password for the userid " + strconv.FormatUint(uint64(userID), 10) + ". Message: " + err.Error())
	}
	return nil
}

func (svc *passwordService) ChangePassword(userID uint, currentPassword, newPlainPassword string) *application_types.ApplicationError {
	password, err := auth.GetPasswordByUserID(userID)
	if err != nil {
		return application_types.NewApplicationError(false, http.StatusInternalServerError, "Password creation failed",
			fmt.Errorf("Unable to find password for the userid "+strconv.FormatUint(uint64(userID), 10)+". Message: "+err.Error()))
	}

	if !password.VerifyPassword(currentPassword) {
		return application_types.NewApplicationError(false, http.StatusUnprocessableEntity, "Unable to change password", fmt.Errorf("Current password is not matching"))
	}

	if err := password.Delete(); err != nil {
		return application_types.NewApplicationError(false, http.StatusInternalServerError, "Unable to change password", err)
	}

	_, err = auth.NewPassword(userID, newPlainPassword)
	if err != nil {
		return application_types.NewApplicationError(false, http.StatusInternalServerError, "Unable to change password", err)
	}
	return nil
}

func (svc *passwordService) ChangePasswordWithoutConfirmingCurrentPassword(userID uint, plainPassword string) *application_types.ApplicationError {
	password, err := auth.GetPasswordByUserID(userID)
	if err != nil {
		return application_types.NewApplicationError(false, http.StatusInternalServerError, "Password creation failed",
			fmt.Errorf("Unable to find password for the userid "+strconv.FormatUint(uint64(userID), 10)+". Message: "+err.Error()))
	}

	if password == nil {
		logger.Warning("Didnt find password for the user. Proceeding to create new one")
	} else {
		logger.Info("Deleting existing password for the user.")
		err = password.Delete()
		if err != nil {
			logger.Warning("Password deleting stopped")
			return application_types.NewApplicationError(false, http.StatusInternalServerError, "Password deletion failed", err)
		}
	}

	password, err = auth.NewPassword(userID, plainPassword)
	if err != nil {
		return application_types.NewApplicationError(false, http.StatusInternalServerError, "Password creation failed",
			fmt.Errorf("Unable to create password for the userid "+strconv.FormatUint(uint64(userID), 10)+". Message: "+err.Error()))
	}

	if err := svc.db.Create(password).Error; err != nil {
		fmt.Errorf("Unable to create password for the userid " + strconv.FormatUint(uint64(userID), 10) + ". Message: " + err.Error())
	}
	return nil
}
