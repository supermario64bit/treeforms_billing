package services

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"treeforms_billing/application_types"
	"treeforms_billing/db"
	"treeforms_billing/dtos"
	"treeforms_billing/logger"
	"treeforms_billing/models"

	"gorm.io/gorm"
)

type userService struct {
	db *gorm.DB
}

type UserService interface {
	Create(user *dtos.UserDTO) (*models.User, *application_types.ApplicationError)
	Find(filter models.UserFilter) ([]*models.User, *application_types.ApplicationError)
	FindByID(id uint) (*models.User, *application_types.ApplicationError)
	UpdateByID(id uint, updatedUserData *dtos.UserDTO) (*models.User, *application_types.ApplicationError)
	DeleteByID(id uint) *application_types.ApplicationError
	FindByEmail(email string) (*models.User, *application_types.ApplicationError)
	FindByPhone(phone string) (*models.User, *application_types.ApplicationError)
}

func NewUserService() UserService {
	return &userService{
		db: db.Get(),
	}
}

func (svc *userService) Create(userDTO *dtos.UserDTO) (*models.User, *application_types.ApplicationError) {
	logger.Info("Creating a new user.")
	// Data transfering from DTO to model
	user := &models.User{Name: userDTO.Name, Email: userDTO.Email, Phone: userDTO.Phone,
		Role: userDTO.Role, Status: userDTO.Status}

	// Validation checks
	logger.Info("Validating new user fields.")
	err := user.ValidateFields()
	if err != nil {
		appErr := application_types.NewApplicationError(false, http.StatusBadRequest, "Invalid request",
			fmt.Errorf("Validation failed for creating the user. Message: "+err.Error()))
		logger.Warning(appErr.GetErrorMessage())
		return nil, appErr
	}

	logger.Info("Checking given email is enrolled by any other user")
	if err := svc.db.Where("email =?", user.Email).First(&models.User{}).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Danger("Error occured while checking the user email enrolled by any other user")
		return nil, application_types.NewApplicationError(false, http.StatusUnprocessableEntity, "User update failed", fmt.Errorf("Unable to find user by email. Message: "+err.Error()))
	} else if err == nil {
		logger.Warning("Given email id already in use")
		return nil, application_types.NewApplicationError(false, http.StatusUnprocessableEntity, "User update failed", fmt.Errorf("Email ID is already registered with another user"))
	}

	logger.Info("Checking given phone is enrolled by any other user")
	if err := svc.db.Where("phone =?", user.Phone).First(&models.User{}).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Danger("Error occured while checking the user phone enrolled by any other user")
		return nil, application_types.NewApplicationError(false, http.StatusUnprocessableEntity, "User update failed", fmt.Errorf("Unable to find user by phone. Message: "+err.Error()))
	} else if err == nil {
		logger.Warning("Given phone already in use")
		return nil, application_types.NewApplicationError(false, http.StatusUnprocessableEntity, "User update failed", fmt.Errorf("Phone number is already registered with another user"))
	}

	// Creating the user
	tx := svc.db.Create(&user)
	if tx.Error != nil {
		appErr := application_types.NewApplicationError(false, http.StatusInternalServerError, "User creation failed",
			fmt.Errorf("User creation failed. Message: "+err.Error()))
		logger.Danger(appErr.GetErrorMessage())
		return nil, appErr
	}

	logger.Success("User created succesfully.")
	return user, nil
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
		return nil, application_types.NewApplicationError(false, http.StatusInternalServerError, "User find failed!",
			fmt.Errorf("Unable to find users. Message: "+err.Error()))
	}

	logger.Success("Users found successfully")
	return users, nil
}

func (svc *userService) FindByID(id uint) (*models.User, *application_types.ApplicationError) {
	user := &models.User{}

	if err := svc.db.First(user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Danger("No user found for the id " + strconv.FormatUint(uint64(id), 10))
			return nil, application_types.NewApplicationError(false, http.StatusNotFound, "No user found", fmt.Errorf("No user found for the given id"))

		}
		logger.Danger("Unable to find user by id. Message: " + err.Error())
		return nil, application_types.NewApplicationError(false, http.StatusInternalServerError, "Unable to find user with id",
			fmt.Errorf("Unable to find user by id. Message: "+err.Error()))
	}

	logger.Success("User found by id!")
	return user, nil
}

func (svc *userService) UpdateByID(id uint, updatedUserData *dtos.UserDTO) (*models.User, *application_types.ApplicationError) {
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
		appErr = application_types.NewApplicationError(false, http.StatusBadRequest, "User update failed",
			fmt.Errorf("Validation failed. Message: "+err.Error()))
		logger.Danger(appErr.GetErrorMessage())
		return nil, appErr
	}
	logger.Info("Checking given email is enrolled by any other user")
	if err := svc.db.Where("email = ? AND id <> ?", updatedUser.Email, id).First(&models.User{}).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Danger("Error occured while checking the user email enrolled by any other user")
		return nil, application_types.NewApplicationError(false, http.StatusUnprocessableEntity, "User update failed", fmt.Errorf("Unable to find user by email. Message: "+err.Error()))
	} else if err == nil {
		logger.Warning("Given email id already in use")
		return nil, application_types.NewApplicationError(false, http.StatusUnprocessableEntity, "User update failed", fmt.Errorf("Email ID is already registered with another user"))
	}

	logger.Info("Checking given phone is enrolled by any other user")
	if err := svc.db.Where("phone = ? AND id <> ?", updatedUser.Phone, id).First(&models.User{}).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Danger("Error occured while checking the user phone enrolled by any other user")
		return nil, application_types.NewApplicationError(false, http.StatusUnprocessableEntity, "User update failed", fmt.Errorf("Unable to find user by phone. Message: "+err.Error()))
	} else if err == nil {
		logger.Warning("Given phone already in use")
		return nil, application_types.NewApplicationError(false, http.StatusUnprocessableEntity, "User update failed", fmt.Errorf("Phone number is already registered with another user"))
	}

	if err := svc.db.Save(updatedUser).Error; err != nil {
		appErr = application_types.NewApplicationError(false, http.StatusInternalServerError, "User update failed.",
			fmt.Errorf("Error occured while updating user. Message: "+err.Error()))
		logger.Danger(appErr.GetErrorMessage())
		return nil, appErr
	}

	logger.Success("User updated by id " + strconv.FormatUint(uint64(id), 10))
	return updatedUser, nil
}

func (svc *userService) DeleteByID(id uint) *application_types.ApplicationError {
	logger.Info("Deleting a user with id " + strconv.FormatUint(uint64(id), 10))

	user, appErr := svc.FindByID(id)
	if appErr != nil {
		return appErr
	}

	if err := svc.db.Delete(user).Error; err != nil {
		appErr = application_types.NewApplicationError(false, http.StatusInternalServerError, "User delete failed.",
			fmt.Errorf("Unable to delete user if id"+strconv.FormatUint(uint64(id), 10)+". Message: "+err.Error()))
		logger.Danger(appErr.GetErrorMessage())
		return appErr
	}

	logger.Success("Deleted user with id " + strconv.FormatUint(uint64(id), 10))
	return nil
}

func (svc *userService) FindByEmail(email string) (*models.User, *application_types.ApplicationError) {
	logger.Info("Finding a user with email id " + email)

	var user models.User
	if err := svc.db.Where("email = ?", email).First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Danger("Unable to find user by email.")
			return nil, application_types.NewApplicationError(false, http.StatusNotFound, "No user found", fmt.Errorf("No user found for the given email"))
		}
		logger.Danger("Unable to find user by email. Message: " + err.Error())
		return nil, application_types.NewApplicationError(false, http.StatusInternalServerError, "User find failed",
			fmt.Errorf("Unable to find user by email. Message: "+err.Error()))
	}

	logger.Success("User found by email")
	return &user, nil
}

func (svc *userService) FindByPhone(phone string) (*models.User, *application_types.ApplicationError) {
	logger.Info("Finding a user with phone " + phone)

	var user models.User
	if err := svc.db.Where("phone = ?", phone).First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Danger("Unable to find user by phone number.")
			return nil, application_types.NewApplicationError(false, http.StatusNotFound, "No user found", fmt.Errorf("No user found for the given phone number"))
		}
		logger.Danger("Unable to find user by phone. Message: " + err.Error())
		return nil, application_types.NewApplicationError(false, http.StatusInternalServerError, "User find failed",
			fmt.Errorf("Unable to find user by phone. Message: "+err.Error()))
	}

	logger.Success("User found by phone")
	return &user, nil
}
