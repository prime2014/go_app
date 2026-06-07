package accounts

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"
	"golang.org/x/crypto/argon2"
	"gorm.io/gorm"
	"errors"
	"strings"
)

// The parameters to use in argon2
type params struct {
	memory uint32
	iterations uint32
	parallelism uint8
	saltLength uint32
	keyLength uint32
}


type Users struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Firstname string         `json:"firstname" validate:"required"`
	Lastname  *string         `json:"lastname" validate:"required"`
	Email     string         `json:"email" validate:"required,email" gorm:"unique;not null;index"`
	Password  string         `json:"password,omitempty" validate:"required"`
	PasswordResetToken *string  `json:"password_reset_token,omitempty"`
	PasswordResetExpiry *time.Time ` json:"password_reset_expiry,omitempty"`
	IsActive  bool           `json:"is_active" gorm:"default:false"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"` // Useful for soft-deletes
}


func generateFromPassword(password string, p *params) (string, error) {
	// Generate a cryptographically secure random salt.
	salt, err := generateRandomBytes(p.saltLength)

	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	// Encode parameters, salt, and hash int oa single standard string format
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
			       argon2.Version, p.memory, p.iterations, p.parallelism, b64Salt, b64Hash)


	return encoded, nil
}


func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)

	if _, err := rand.Read(b); err != nil {
		return nil, err
	}

	return b, nil
}


func (u *Users) BeforeCreate(tx *gorm.DB) (err error) {
	p := &params{
		memory: 64 * 1024,
		iterations: 3,
		parallelism: 2,
		saltLength: 16,
		keyLength: 32,
	}

	// Clean up string spaces and check if empty instead of nil
	u.Password = strings.TrimSpace(u.Password)

	if u.Password == "" {
		return errors.New("password cannot be empty")
	}

	formattedHash, err := generateFromPassword(u.Password, p)
	if err != nil {
		return err
	}

	u.Password = string(formattedHash)
	return nil
}
