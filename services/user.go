package services

import (
	"fmt"
	"net/http"
	"treeforms_billing/application_types"
	"treeforms_billing/db"
	"treeforms_billing/logger"
	"treeforms_billing/models"

	"gorm.io/gorm"
)

type userService struct {
	db *gorm.DB
}

type UserService interface {
	Create(user models.User) (*models.User, *application_types.ApplicationError)
}

func NewUserService() UserService {
	return &userService{
		db: db.Get(),
	}
}

func (svc *userService) Create(user models.User) (*models.User, *application_types.ApplicationError) {
	// Validation checks
	err := user.ValidateFields()
	if err != nil {
		appErr := application_types.NewApplicationError(false, http.StatusBadRequest,
			fmt.Errorf("Validation failed for creating the user. Message: "+err.Error()))
		logger.Warning(appErr.GetErrorMessage())
		return nil, appErr
	}

	// TODO -> Check email and mobile already exist

	// Creating the user
	tx := svc.db.Create(&user)
	if tx.Error != nil {
		appErr := application_types.NewApplicationError(false, http.StatusInternalServerError,
			fmt.Errorf("User creation failed. Message: "+err.Error()))
		logger.Danger(appErr.GetErrorMessage())
		return nil, appErr
	}

	return &user, nil
}
