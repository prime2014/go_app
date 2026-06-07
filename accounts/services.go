package accounts

import (
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
	"gorm.io/gorm"
)

type UserService struct {
	Db *gorm.DB
}

func (u *UserService) SignupUser(dto SignUpDto) (*Users, error) {
	fmt.Println("Inside SignupUser service")
	user := &Users{
		Firstname: dto.Firstname,
		Lastname:  &dto.Lastname,
		Email:     dto.Email,
		Password:  dto.Password,
	}

	result := u.Db.Create(user)
	fmt.Println("After creating user")

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil

}

func verifyPassword(password, hashedPassword string) (bool, error) {
	// 1. Split the encoded string into components
	// Expected format: ["", "argon2id", "v=19", "m=65536,t=3,p=2", "saltString", "hashString"]
	parts := strings.Split(hashedPassword, "$")
	if len(parts) != 6 {
		return false, errors.New("invalid encoded hash format")
	}

	// 2. Parse the parameters out of the 4th segment
	var version int
	_, err := fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil {
		return false, err
	}

	var p params
	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations, &p.parallelism)
	if err != nil {
		return false, err
	}

	// 3. Decode the base64-encoded salt and hash back to raw bytes
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	existingHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}
	p.keyLength = uint32(len(existingHash))

	// 4. Derive the hash from the incoming plain text password using the decoded parameters/salt
	comparisonHash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	// 5. Use subtle.ConstantTimeCompare to protect against timing attacks
	if subtle.ConstantTimeCompare(existingHash, comparisonHash) == 1 {
		return true, nil
	}

	return false, nil
}

func (u *UserService) LoginUser(dto LoginDto) (*Users, error) {
	var user Users

	result := u.Db.Where("email = ?", dto.Email).First(&user)

	if result.Error != nil {
		return nil, errors.New("Invalid email!")
	}

	isValid, err := verifyPassword(dto.Password, user.Password)

	if err != nil {
		return nil, errors.New("internal server error during verification")
	}

	if !isValid {
		return nil, errors.New("invalid email or password")
	}

	return &user, nil
}
