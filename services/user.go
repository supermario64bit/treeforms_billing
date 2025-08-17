package services

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
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
	Find(filter models.UserFilter) ([]*models.User, *application_types.ApplicationError)
	FindByID(id uint) (*models.User, *application_types.ApplicationError)
	UpdateByID(id uint, updatedUserData *models.User) (*models.User, *application_types.ApplicationError)
}

func NewUserService() UserService {
	return &userService{
		db: db.Get(),
	}
}

func (svc *userService) Create(user models.User) (*models.User, *application_types.ApplicationError) {
	logger.Info("Creating a new user.")
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

	logger.Success("User created succesfully.")
	return &user, nil
}

func (svc *userService) Find(filter models.UserFilter) ([]*models.User, *application_types.ApplicationError) {
	logger.Info("Finding users")
	var users []*models.User
	query := svc.db

	if strings.TrimSpace(filter.Name) != "" {
		logger.Info("Added Name filter to the user find query")
		query = query.Where("name LIKE ?", "%"+strings.TrimSpace(filter.Name)+"%")
	}

	if strings.TrimSpace(filter.Email) != "" {
		logger.Info("Added Email filter to the user find query")
		query = query.Where("email LIKE ?", "%"+strings.TrimSpace(filter.Email)+"%")

	}

	if strings.TrimSpace(filter.Phone) != "" {
		logger.Info("Added Phone filter to the user find query")
		query = query.Where("phone LIKE ?", "%"+strings.TrimSpace(filter.Phone)+"%")

	}

	if strings.TrimSpace(filter.Role) != "" {
		logger.Info("Added Role filter to the user find query")
		query = query.Where("role LIKE ?", "%"+strings.TrimSpace(filter.Role)+"%")

	}

	if strings.TrimSpace(filter.Status) != "" {
		logger.Info("Added Status filter to the user find query")
		query = query.Where("status LIKE ?", "%"+strings.TrimSpace(filter.Status)+"%")

	}

	if err := query.Find(&users).Error; err != nil {
		logger.Danger("Unable to find users. Message: " + err.Error())
		return nil, application_types.NewApplicationError(false, http.StatusInternalServerError, fmt.Errorf("Unable to find users. Message: "+err.Error()))
	}

	logger.Success("Users found successfully")
	return users, nil
}

func (svc *userService) FindByID(id uint) (*models.User, *application_types.ApplicationError) {
	var user *models.User

	if err := svc.db.First(user, id).Error; err != nil {
		logger.Danger("Unable to find user by id. Message: " + err.Error())
		return nil, application_types.NewApplicationError(false, http.StatusInternalServerError, fmt.Errorf("Unable to find user by id. Message: "+err.Error()))
	}

	logger.Success("User found by id!")
	return user, nil
}

func (svc *userService) UpdateByID(id uint, updatedUserData *models.User) (*models.User, *application_types.ApplicationError) {
	logger.Info("Started updated user by id " + strconv.FormatUint(uint64(id), 10))
	updatedUser, appErr := svc.FindByID(id)
	if appErr != nil {
		return nil, appErr
	}

	if strings.TrimSpace(updatedUserData.Name) != "" {
		logger.Info("User name changed")
		updatedUser.Name = updatedUserData.Name
	}

	if strings.TrimSpace(updatedUserData.Email) != "" {
		logger.Info("User email changed")
		updatedUser.Email = updatedUserData.Email
	}

	if strings.TrimSpace(updatedUserData.Phone) != "" {
		logger.Info("User phone changed")
		updatedUser.Phone = updatedUserData.Phone
	}

	if strings.TrimSpace(updatedUserData.Role) != "" {
		logger.Info("User role changed")
		updatedUser.Role = updatedUserData.Role
	}

	if strings.TrimSpace(updatedUserData.Status) != "" {
		logger.Info("User status changed")
		updatedUser.Status = updatedUserData.Status
	}

	if err := updatedUser.ValidateFields(); err != nil {
		appErr = application_types.NewApplicationError(false, http.StatusBadRequest, fmt.Errorf("Validation failed. Message: "+err.Error()))
		logger.Danger(appErr.GetErrorMessage())
		return nil, appErr
	}

	if err := svc.db.Save(updatedUser).Error; err != nil {
		appErr = application_types.NewApplicationError(false, http.StatusInternalServerError, fmt.Errorf("Error occured while updating user. Message: "+err.Error()))
		logger.Danger(appErr.GetErrorMessage())
		return nil, appErr
	}

	logger.Success("User updated by id " + strconv.FormatUint(uint64(id), 10))
	return updatedUser, nil
}
