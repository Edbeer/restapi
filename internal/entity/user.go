package entity

import (
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User model
type User struct {
	ID          uuid.UUID `json:"user_id" db:"user_id"`
	FirstName   string    `json:"first_name" db:"first_name"`
	LastName    string    `json:"last_name" db:"last_name"`
	Email       string    `json:"email" db:"email"`
	PhoneNumber *string   `json:"phone_number" db:"phone_number"`
	Role        string    `json:"role" db:"role"`
	Address     *string   `json:"address" db:"address"`
	City        *string   `json:"city" db:"city"`
	Password    string    `json:"-" db:"password"`
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
func (u *User) Prepare() error {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
	u.Password = strings.TrimSpace(u.Password)
	if u.PhoneNumber != nil {
		*u.PhoneNumber = strings.TrimSpace(*u.PhoneNumber)
	}
	u.Role = strings.ToLower(strings.TrimSpace(u.Role))

	return nil
}