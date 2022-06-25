package service

import (
	"errors"
	"fmt"
	"mygram-go/entity"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type Userinterf interface {
	Register(user *entity.User) (*entity.User, error)
	Login(user *entity.User, tempPassword string) (*entity.User, error)
	GetToken(id uint, email string, password string) string
	CheckToken(compareToken string, id uint, email string, password string) error
	VerivyToken(TempToken string) float64
}

type UserService struct {
}

func NewUserService() Userinterf {
	return &UserService{}
}

func (servicuser *UserService) Register(user *entity.User) (*entity.User, error) {
	if user.Username == "" {
		return nil, errors.New("Username cannot be empty")
	}
	if user.Email == "" {
		return nil, errors.New("Email cannot be empty")
	}
	if len(user.Password) < 6 {
		return nil, errors.New("password must be minimum 6 characters")
	}
	if user.Age < 8 {
		return nil, errors.New("age must be greater than 8")
	}
	fmt.Println("cek user")
	return user, nil
}

func (servicuser *UserService) Login(user *entity.User, tempPassword string) (*entity.User, error) {
	if user.Email == "" {
		return nil, errors.New("email cannot be empty")
	}
	password := []byte(tempPassword)
	//check password salah
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), password); err != nil {
		return nil, errors.New("password salah")
	}
	return user, nil
}

func (servicuser *UserService) GetToken(id uint, email string, password string) string {
	claims := jwt.MapClaims{
		"id":    id,
		"email": email,
	}

	parseToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := parseToken.SignedString([]byte(password))

	return signedToken
}

func (servicuser *UserService) CheckToken(compareToken string, id uint, email string, password string) error {
	token := servicuser.GetToken(id, email, password)
	if compareToken == token {
		fmt.Println("berhasil")
		return nil
	} else {
		fmt.Println("tidak berhasil")
		return errors.New("username cannot be empty")
	}
	//compare
}

func (servicuser *UserService) VerivyToken(TempToken string) float64 {
	tokenString := TempToken
	claims := jwt.MapClaims{}
	var verivykey []byte
	token, _ := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return verivykey, nil
	})
	payload := token.Claims.(jwt.MapClaims)
	id := payload["id"].(float64)
	// fmt.Println(token.Claims.Valid())
	return id
}
