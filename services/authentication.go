package services

import (
	"fmt"
	"net/http"
	"treeforms_billing/application_types"
	"treeforms_billing/auth"
	"treeforms_billing/db"
	"treeforms_billing/dtos"
	"treeforms_billing/logger"
	"treeforms_billing/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type authenticationService struct {
	userSvc UserService
	passSvc PasswordService
	db      *gorm.DB
}

type AuthenticationService interface{
	EmailLogin(emailID string, passwordStr string) *application_types.ApplicationError
	Signup(signup dtos.SignupDTO) *application_types.ApplicationError
}

func NewAuthenticationSevice() AuthenticationService {
	return &authenticationService{
		userSvc: NewUserService(),
		passSvc: NewPasswordService(),
		db:      db.Get(),
	}
}

func (svc *authenticationService) EmailLogin(emailID string, passwordStr string) *application_types.ApplicationError {
	user, appErr := svc.userSvc.FindByEmail(emailID)
	if appErr != nil {
		return appErr
	}

	password, err := auth.GetPasswordByUserID(user.ID)
	if err != nil {
		return application_types.NewApplicationError(false, http.StatusInternalServerError, "Login using email failed", err)
	}

	if !password.VerifyPassword(passwordStr) {
		return application_types.NewApplicationError(false, http.StatusUnauthorized, "Login using email failed", fmt.Errorf("Invalid Credentials"))
	}

	return nil
}

func (svc *authenticationService) Signup(signup dtos.SignupDTO) *application_types.ApplicationError {
	if signup.Password != signup.ConfirmPassword {
		return application_types.NewApplicationError(false, http.StatusUnprocessableEntity, "Signup Failed", fmt.Errorf("Password and Confirm passwords are not identical"))
	}

	user := &models.User{
		Name:   signup.Name,
		Email:  signup.Email,
		Phone:  signup.Phone,
		Role:   "user",
		Status: "active",
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(signup.ConfirmPassword), bcrypt.DefaultCost)
	if err != nil {
		return application_types.NewApplicationError(false, http.StatusInternalServerError, "Signup Failed", fmt.Errorf("error occurred while hashing password: %w", err))
	}

	err = svc.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		passwordInsertQuery := `INSERT INTO passwords (hash, user_id) VALUES (?, ?);`
		if err := tx.Raw(passwordInsertQuery, string(hashedPass), user.ID).Error; err != nil {
			logger.HighlightedDanger("Password creation failed. Message: " + err.Error())
			return err
		}

		return nil
	})

	if err != nil {
		return application_types.NewApplicationError(false, http.StatusInternalServerError, "Signup Failed", err)
	}

	return nil
}
