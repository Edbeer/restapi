package entity

import (
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User model
type User struct {
	ID          uuid.UUID `json:"user_id" db:"user_id" validate:"omitempty,uuid"`
	FirstName   string    `json:"first_name" db:"first_name" validate:"required,lte=30"`
	LastName    string    `json:"last_name" db:"last_name" validate:"required,lte=30"`
	Email       string    `json:"email" db:"email" validate:"omitempty,lte=60,email"`
	PhoneNumber *string   `json:"phone_number" db:"phone_number" validate:"omitempty,lte=20"`
	Role        *string    `json:"role" db:"role" validate:"omitempty,lte=10"`
	Address     *string   `json:"address" db:"address" validate:"omitempty,lte=250"`
	City        *string   `json:"city" db:"city" validate:"omitempty,lte=24"`
	Country     *string   `json:"country" db:"country" validate:"omitempty,lte=24"`
	Postcode    *int      `json:"postcode" db:"postcode" validate:"omitempty,lte=10"`
	Balance     float64   `json:"balance" db:"balance"`
	Avatar      *string   `json:"avatar" db:"avatar"`
	Password    string    `json:"-" db:"password" validate:"required,gte=6"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Hash user password with bcrypt
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err 
	}

	u.Password = string(hashedPassword)
	return nil
}

// Compare user password and payload
func (u *User) ComparePassword(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return err
	}
	return nil
}

// Sanitize user password
func (u *User) SanitizePassword() {
	u.Password = ""
}

// Prepare user struct for register
func (u *User) PrepareCreate() error {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
	u.Password = strings.TrimSpace(u.Password)
	
	if err := u.HashPassword(); err != nil {
		return err
	}

	if u.PhoneNumber != nil {
		*u.PhoneNumber = strings.TrimSpace(*u.PhoneNumber)
	}
	
	if u.Role != nil {
		*u.Role = strings.ToLower(strings.TrimSpace(*u.Role))
	}


	return nil
}