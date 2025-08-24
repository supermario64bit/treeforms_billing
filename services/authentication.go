package services

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"time"
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

type AuthenticationService interface {
	EmailLogin(emailID string, passwordStr string) (access_token, refresh_token string, appErr *application_types.ApplicationError)
	Signup(signup dtos.SignupDTO) *application_types.ApplicationError
}

func NewAuthenticationSevice() AuthenticationService {
	return &authenticationService{
		userSvc: NewUserService(),
		passSvc: NewPasswordService(),
		db:      db.Get(),
	}
}

func (svc *authenticationService) EmailLogin(emailID string, passwordStr string) (access_token, refresh_token string, appErr *application_types.ApplicationError) {
	logger.Info("Email login Service Started")
	user, appErr := svc.userSvc.FindByEmail(emailID)
	if appErr != nil {
		logger.Danger("Stopping Email login service.")
		return
	}

	password, err := auth.GetPasswordByUserID(user.ID)
	if err != nil {
		logger.HighlightedDanger("Stopping Email login service.")
		return "", "", application_types.NewApplicationError(false, http.StatusInternalServerError, "Login using email failed", err)
	}

	if password == nil {
		logger.HighlightedDanger("Password not created for the user")
		return "", "", application_types.NewApplicationError(false, http.StatusInternalServerError, "Login using email failed", fmt.Errorf("Password not created for the user"))
	}

	if !password.VerifyPassword(passwordStr) {
		logger.Info("Stopping Email login service. Message: Invalid password")
		return "", "", application_types.NewApplicationError(false, http.StatusUnauthorized, "Login using email failed", fmt.Errorf("Invalid Credentials"))
	}

	logger.Success("Email login is success")
	tokenStr, err := user.NewAccessToken()
	if err != nil {
		logger.HighlightedDanger("Error occured while signing access token")
		return "", "", application_types.NewApplicationError(false, http.StatusInternalServerError, "Access Token Signing Failed", err)
	}

	refresh_token, appErr = svc.NewRefreshToken(user.ID)
	if appErr != nil {
		logger.HighlightedDanger("Error occured while generation refresh token")
		return "", "", application_types.NewApplicationError(false, http.StatusInternalServerError, "Refresh Token Generation Failed", err)
	}
	return tokenStr, refresh_token, nil
}

func (svc *authenticationService) Signup(signup dtos.SignupDTO) *application_types.ApplicationError {
	logger.Info("User Signup Service Started")
	if signup.Password != signup.ConfirmPassword {
		logger.Danger("Confirm password and password are not identical.")
		return application_types.NewApplicationError(false, http.StatusUnprocessableEntity, "Signup Failed", fmt.Errorf("Password and Confirm passwords are not identical"))
	}

	user := &models.User{
		Name:   signup.Name,
		Email:  signup.Email,
		Phone:  signup.Phone,
		Role:   "user",
		Status: "active",
	}

	u, appErr := svc.userSvc.FindByEmail(user.Email)
	if appErr != nil && !errors.Is(appErr.GetError(), gorm.ErrRecordNotFound) {
		logger.Danger("User signup service stopped")
		return appErr
	} else if u != nil {
		logger.Warning("User signup service stopped. Message: Email already registered with another user.")
		return application_types.NewApplicationError(false, http.StatusUnprocessableEntity, "Signup failed.", fmt.Errorf("Email already registered with another user."))
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(signup.ConfirmPassword), bcrypt.DefaultCost)
	if err != nil {
		logger.HighlightedDanger("User signup failed. Unable to hash the password. Message: " + err.Error())
		return application_types.NewApplicationError(false, http.StatusInternalServerError, "Signup Failed", fmt.Errorf("error occurred while hashing password: %w", err))
	}

	err = svc.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			logger.HighlightedDanger("User signup failed. Unable to create user in db. Gorm Message: " + err.Error())
			return err
		}

		passwordInsertQuery := `INSERT INTO passwords (hash, user_id) VALUES (?, ?);`
		logger.Info("Executing query: " + passwordInsertQuery)
		if err := tx.Exec(passwordInsertQuery, string(hashedPass), user.ID).Error; err != nil {
			logger.HighlightedDanger("Password creation failed. Gorm Message: " + err.Error())
			return err
		}

		return nil
	})

	if err != nil {
		return application_types.NewApplicationError(false, http.StatusInternalServerError, "Signup Failed", err)
	}

	logger.Success("User Signup Service Success")
	return nil
}

func (svc *authenticationService) NewAccessToken(userID uint) (string, *application_types.ApplicationError) {
	logger.Info("Started New Access Token Service")
	user, appErr := svc.userSvc.FindByID(userID)
	if appErr != nil {
		logger.Danger("New Access Token Service stopped")
		return "", appErr
	}

	tokenStr, err := user.NewAccessToken()
	if err != nil {
		logger.HighlightedDanger("Error occured while signing access token")
		return "", application_types.NewApplicationError(false, http.StatusInternalServerError, "Access Token Signing Failed", err)
	}

	return tokenStr, nil
}

func (svc *authenticationService) NewRefreshToken(userID uint) (string, *application_types.ApplicationError) {
	logger.Info("Started New Refresh Token Service")
	user, appErr := svc.userSvc.FindByID(userID)
	if appErr != nil {
		logger.Danger("New Refresh Token Service stopped")
		return "", appErr
	}

	bytes := make([]byte, 35)
	_, err := rand.Read(bytes)
	if err != nil {
		logger.HighlightedDanger("Error occured while generating refresh token")
		return "", application_types.NewApplicationError(false, http.StatusInternalServerError, "Access Refresh Signing Failed", err)
	}

	tokenStr := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(bytes)
	tokenHash, err := bcrypt.GenerateFromPassword([]byte(tokenStr), bcrypt.DefaultCost)
	if err != nil {
		logger.HighlightedDanger("Error occured while hashing refresh token")
		return "", application_types.NewApplicationError(false, http.StatusInternalServerError, "Access Refresh Signing Failed", err)
	}

	if err := svc.db.Create(&models.RefreshToken{UserID: user.ID, TokenHash: string(tokenHash), ExpiresAt: time.Now().Add(7 * (24 * time.Hour))}).Error; err != nil {
		logger.HighlightedDanger("Error occured while generating refresh token. Gorm Message: " + err.Error())
		return "", application_types.NewApplicationError(false, http.StatusInternalServerError, "Access Refresh Signing Failed.", err)

	}

	return tokenStr, nil
}
