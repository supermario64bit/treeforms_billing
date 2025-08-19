package auth

import (
	"fmt"
	"strconv"
	"treeforms_billing/db"
	"treeforms_billing/logger"

	"golang.org/x/crypto/bcrypt"
)

type password struct {
	id     uint
	hash   []byte
	userID uint
}

func NewPassword(userID uint, plainPassword string) (*password, error) {
	db := db.Get()
	hashed, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error occurred while hashing password: %w", err)
	}

	var id uint
	insertQuery := `INSERT INTO passwords (hash, user_id) VALUES (?, ?) RETURNING id;`

	if err := db.Raw(insertQuery, string(hashed), userID).Scan(&id).Error; err != nil {
		logger.HighlightedDanger("Password creation failed. Message: " + err.Error())
		return nil, err
	}

	password := &password{
		id:     id,
		userID: userID,
		hash:   hashed,
	}

	return password, nil
}

func GetPasswordByUserID(userID int) (*password, error) {
	var p *password
	db := db.Get()

	selectQuery := `SELECT id, hash, user_id FROM passwords WHERE user_id = ?`

	rows, err := db.Raw(selectQuery, userID).Rows()
	if err != nil {
		logger.HighlightedDanger("Query execution failed for getting password using userid. Message: " + err.Error())
		return nil, fmt.Errorf("Query execution failed for getting password using userid. Message: " + err.Error())
	}

	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&p.id, &p.hash, &p.userID); err != nil {
			logger.HighlightedDanger("Scan failed getting password using userid. Message: " + err.Error())
			return nil, fmt.Errorf("Scan failed getting password using userid. Message: " + err.Error())
		}
	}
	return p, nil
}

func (p *password) VerifyPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(p.hash), []byte(password)) == nil
}

func (p *password) GetUserID() uint {
	return p.userID
}

func (p *password) Delete() error {
	db := db.Get()

	res := db.Exec(`DELETE FROM passwords WHERE id = ?`, p.id)
	if res.Error != nil {
		logger.HighlightedDanger("Deleting password failed for the user id " + strconv.FormatUint(uint64(p.userID), 10) + ". Message: " + res.Error.Error())
		return fmt.Errorf("Deleting password failed for the user id " + strconv.FormatUint(uint64(p.userID), 10) + ". Message: " + res.Error.Error())
	}

	logger.Warning("Deleting password for the user id " + strconv.FormatUint(uint64(p.userID), 10))
	return nil
}
