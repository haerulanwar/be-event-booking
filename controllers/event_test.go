package controllers

import (
	"bytes"
	"encoding/json"
	"event-booking/common/constant"
	"event-booking/common/request"
	"event-booking/config"
	"event-booking/middleware"
	"event-booking/models"
	"fmt"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestGetEvents(t *testing.T) {
	if err := godotenv.Load(filepath.Join("../", ".env")); err != nil {
		t.Fatal(err)
	}
	config.ConnectDB()
	app := fiber.New()
	app.Get("/api/events", middleware.JWTMiddleware, GetEvents)

	// Prepare test data
	userHR := models.User{Username: "testuserHR", Password: "testpassword", Role: constant.HR}
	config.DB.Create(&userHR)
	userVendor := models.User{Username: "testuserVendor", Password: "testpassword", Role: constant.VENDOR}
	config.DB.Create(&userVendor)

	event1 := models.Event{CompanyName: "Company A", ProposedDates: "2024-07-20", Location: "Location A", EventName: "Event A", CreatedBy: userHR.ID}
	config.DB.Create(&event1)
	event2 := models.Event{CompanyName: "Company B", ProposedDates: "2024-07-21", Location: "Location B", EventName: "Event B", VendorID: userVendor.ID}
	config.DB.Create(&event2)

	// Test cases
	testCases := []struct {
		description        string
		role               string
		userId             uint
		expectedEvents     []models.Event
		expectedStatusCode int
	}{
		{
			description:        "HR user gets their created events",
			role:               constant.HR,
			userId:             userHR.ID,
			expectedEvents:     []models.Event{event1},
			expectedStatusCode: fiber.StatusOK,
		},
		{
			description:        "Vendor user gets their assigned events",
			role:               constant.VENDOR,
			userId:             userVendor.ID,
			expectedEvents:     []models.Event{event2},
			expectedStatusCode: fiber.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			token := generateTestToken(tc.userId, tc.role)

			req := httptest.NewRequest("GET", "/api/events", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			resp, _ := app.Test(req)

			assert.Equal(t, tc.expectedStatusCode, resp.StatusCode)

			if tc.expectedStatusCode == fiber.StatusOK {
				var events []models.Event
				json.NewDecoder(resp.Body).Decode(&events)
				assert.Equal(t, len(tc.expectedEvents), len(events))
				for i, expectedEvent := range tc.expectedEvents {
					assert.Equal(t, expectedEvent.CompanyName, events[i].CompanyName)
					assert.Equal(t, expectedEvent.ProposedDates, events[i].ProposedDates)
					assert.Equal(t, expectedEvent.Location, events[i].Location)
				}
			}
		})
	}
	config.DB.Delete(&event1)
	config.DB.Delete(&event2)
	config.DB.Delete(&userHR)
	config.DB.Delete(&userVendor)
}

func TestApproveEvent(t *testing.T) {
	if err := godotenv.Load(filepath.Join("../", ".env")); err != nil {
		t.Fatal(err)
	}
	config.ConnectDB()
	app := fiber.New()
	app.Post("/api/events/:id/approve", middleware.JWTMiddleware, ApproveEvent)

	userHR := models.User{Username: "testhr", Password: "password", Role: constant.HR}
	config.DB.Create(&userHR)

	event := models.Event{CompanyName: "Company A", ProposedDates: "2024-07-20", Location: "Location A", EventName: "Event A", CreatedBy: userHR.ID}
	config.DB.Create(&event)

	reqBody := request.ApproveEventRequest{ConfirmedDate: "2024-07-22"}
	body, _ := json.Marshal(reqBody)

	token := generateTestToken(userHR.ID, constant.HR)
	req := httptest.NewRequest("POST", fmt.Sprintf("/api/events/%d/approve", event.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var updatedEvent models.Event
	config.DB.First(&updatedEvent, event.ID)
	assert.Equal(t, constant.APPROVED, updatedEvent.Status)
	assert.Equal(t, "2024-07-22", updatedEvent.ConfirmedDate)
	config.DB.Delete(&event)
	config.DB.Delete(&userHR)

}

func TestRejectEvent(t *testing.T) {
	if err := godotenv.Load(filepath.Join("../", ".env")); err != nil {
		t.Fatal(err)
	}
	config.ConnectDB()
	app := fiber.New()
	app.Post("/api/events/:id/reject", middleware.JWTMiddleware, RejectEvent)

	// Create HR user for authentication
	hrUser := models.User{Username: "testuser2", Password: "password", Role: constant.HR}
	config.DB.Create(&hrUser)
	defer config.DB.Delete(&hrUser)

	// Create a test event
	event := models.Event{CompanyName: "Test Company", ProposedDates: "2024-08-15", Location: "Test Location", EventName: "Test Event", Status: constant.PENDING}
	config.DB.Create(&event)
	defer config.DB.Delete(&event)

	testCases := []struct {
		description  string
		requestBody  request.RejectEventRequest
		expectedCode int
		expectedMsg  string
	}{
		{
			description: "Valid request",
			requestBody: request.RejectEventRequest{
				Remarks: "Not suitable",
			},
			expectedCode: fiber.StatusOK,
			expectedMsg:  "Event rejected successfully",
		},
		{
			description:  "invalid request body",
			requestBody:  request.RejectEventRequest{},
			expectedCode: fiber.StatusOK,
			expectedMsg:  "Event rejected successfully",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			body, _ := json.Marshal(tc.requestBody)
			req := httptest.NewRequest(fiber.MethodPost, "/api/events/"+strconv.Itoa(int(event.ID))+"/reject", bytes.NewReader(body))
			req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

			token := generateTestToken(hrUser.ID, hrUser.Role)
			req.Header.Set("Authorization", "Bearer "+token)

			resp, err := app.Test(req)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tc.expectedCode, resp.StatusCode)

			var response map[string]string
			json.NewDecoder(resp.Body).Decode(&response)
			assert.Equal(t, tc.expectedMsg, response["message"])

			updatedEvent := models.Event{}
			config.DB.First(&updatedEvent, event.ID)

			assert.Equal(t, constant.REJECTED, updatedEvent.Status)
			assert.Equal(t, tc.requestBody.Remarks, updatedEvent.Remarks)
		})
	}
}

func generateTestToken(userId uint, role string) string {
	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString
}
