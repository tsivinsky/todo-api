package router

import (
	"errors"
	"log"
	"time"
	"todo-app/internal/db"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthRegisterBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthLoginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func handleAuthRegister(c *fiber.Ctx) error {
	body := new(AuthRegisterBody)

	if err := c.BodyParser(body); err != nil {
		log.Print("error parsing auth register body", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var emailExists db.User
	if err := db.Db.First(&emailExists, "email = ?", body.Email).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Print("error checking if email exists", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "email already registered",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), 8)
	if err != nil {
		log.Print("error hashing password in register", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	user := db.User{
		Email:    body.Email,
		Password: string(hashedPassword),
	}
	if err := db.Db.Create(&user).Error; err != nil {
		log.Print("error creating user in register", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	accessToken, refreshToken, err := generateAuthTokens(user.ID)
	if err != nil {
		log.Print("error generating tokens in register", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"user":         user,
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func handleAuthLogin(c *fiber.Ctx) error {
	body := new(AuthLoginBody)

	if err := c.BodyParser(body); err != nil {
		log.Print("error parsing body in login", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	user := db.User{}
	if err := db.Db.First(&user, "email", body.Email).Error; err != nil {
		log.Print("error finding user by email in login", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid email",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		log.Print("error comparing hash and password in login", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid password",
		})
	}

	accessToken, refreshToken, err := generateAuthTokens(user.ID)
	if err != nil {
		log.Print("error generating tokens in login", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user":         user,
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func generateAuthTokens(userId uint) (string, string, error) {
	now := time.Now()

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"nbf":    now.Add(1 * time.Hour).Unix(),
	}).SignedString("JWT SECRET")
	if err != nil {
		log.Print("error generating accessToken", err)
		return "", "", err
	}

	refreshToken, err := jwt.New(jwt.SigningMethodHS256).SignedString("JWT SECRET")
	if err != nil {
		log.Print("error generating refreshToken", err)
		return "", "", err
	}

	return accessToken, refreshToken, err
}
