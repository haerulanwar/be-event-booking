package controllers

import (
	"bytes"
	"encoding/json"
	"event-booking/common/request"
	"event-booking/config"
	"event-booking/models"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func TestLogin(t *testing.T) {
	app := fiber.New()
	app.Post("/login", Login)
	if err := godotenv.Load(filepath.Join("../", ".env")); err != nil {
		t.Fatal(err)
	}
	config.ConnectDB()

	// Test with correct credentials
	t.Run("Correct credentials", func(t *testing.T) {
		password := "testpassword"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		user := models.User{Username: "testuser", Password: string(hashedPassword), Role: "HR"}
		config.ConnectDB()
		config.DB.Create(&user)
		defer config.DB.Delete(&user)

		reqBody := request.LoginRequest{
			Username: "testuser",
			Password: password,
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != fiber.StatusOK {
			t.Errorf("Expected status code %d, got %d", fiber.StatusOK, resp.StatusCode)
		}

		var respBody map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&respBody)

		if _, ok := respBody["token"]; !ok {
			t.Error("Missing token in response body")
		}
		if respBody["role"] != user.Role {
			t.Errorf("Expected Role to be %s but got %s", user.Role, respBody["role"])
		}
	})

	// Test with incorrect password
	t.Run("Incorrect Password", func(t *testing.T) {
		password := "testpassword"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		user := models.User{Username: "testuser", Password: string(hashedPassword), Role: "HR"}
		config.DB.Create(&user)
		defer config.DB.Delete(&user)

		reqBody := request.LoginRequest{
			Username: "testuser",
			Password: "wrongpassword",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != fiber.StatusUnauthorized {
			t.Errorf("Expected status code %d, got %d", fiber.StatusUnauthorized, resp.StatusCode)
		}
	})

	// Test with non-existent user
	t.Run("Non-existent user", func(t *testing.T) {
		reqBody := request.LoginRequest{
			Username: "nonexistentuser",
			Password: "password",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != fiber.StatusUnauthorized {
			t.Errorf("Expected status code %d, got %d", fiber.StatusUnauthorized, resp.StatusCode)
		}
	})

	t.Run("Invalid Input", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"username": 123,
			"password": "test",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		if resp.StatusCode != fiber.StatusBadRequest {
			t.Fatalf("Expected status code %d but got %d", fiber.StatusBadRequest, resp.StatusCode)
		}
	})
}
