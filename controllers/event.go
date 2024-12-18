package controllers

import (
	"event-booking/common/constant"
	"event-booking/common/request"
	"event-booking/config"
	"event-booking/models"

	"github.com/gofiber/fiber/v2"
)

// @Summary Get Events
// @Description Fetch events based on user role (HR or Vendor)
// @Tags Event
// @Produce json
// @Success 200 {array} models.Event
// @Router /api/events [get]
// @Security Bearer
func GetEvents(c *fiber.Ctx) error {
	// var events []models.Event
	var events []models.EventWithVendorName
	role := c.Locals("role").(string)
	userId := uint(c.Locals("user_id").(float64))

	if role == constant.HR {
		// if err := config.DB.Where("created_by = ?", userId).Find(&events).Error; err != nil {
		// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to fetch events"})
		// }
		if err := config.DB.Model(&models.Event{}).Select("events.id, events.company_name, events.proposed_dates, events.location, events.event_name, events.status, events.remarks, events.confirmed_date, events.created_by, events.created_at, users.full_name as vendor_name").Where("events.created_by = ?", userId).Joins("JOIN users ON events.vendor_id = users.id").Scan(&events).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to fetch events"})
		}
	} else if role == constant.VENDOR {
		// if err := config.DB.Where("vendor_id = ?", userId).Joins("JOIN users ON events.created_by = users.id").Find(&events).Error; err != nil {
		// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to fetch events"})
		// }
		if err := config.DB.Model(&models.Event{}).Select("events.id, events.company_name, events.proposed_dates, events.location, events.event_name, events.status, events.remarks, events.confirmed_date, events.created_by, events.created_at, users.full_name as vendor_name").Where("events.vendor_id = ?", userId).Joins("JOIN users ON events.vendor_id = users.id").Scan(&events).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to fetch events"})
		}
	}

	return c.JSON(events)
}

// @Summary Approve Event
// @Description Approve an event and set a confirmed date
// @Tags Event
// @Accept json
// @Produce json
// @Param id path int true "Event ID"
// @Param request body request.ApproveEventRequest true "Confirmed date"
// @Success 200 {object} map[string]string
// @Router /api/events/{id}/approve [post]
// @Security Bearer
func ApproveEvent(c *fiber.Ctx) error {
	id := c.Params("id")
	var input request.ApproveEventRequest

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid input"})
	}

	if err := config.DB.Model(&models.Event{}).Where("id = ?", id).
		Updates(map[string]interface{}{"status": constant.APPROVED, "confirmed_date": input.ConfirmedDate}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to approve event"})
	}

	return c.JSON(fiber.Map{"message": "Event approved successfully"})
}

// @Summary Reject Event
// @Description Reject an event with remarks
// @Tags Event
// @Accept json
// @Produce json
// @Param id path int true "Event ID"
// @Param request body request.RejectEventRequest true "Remarks"
// @Success 200 {object} map[string]string
// @Router /api/events/{id}/reject [post]
// @Security Bearer
func RejectEvent(c *fiber.Ctx) error {
	id := c.Params("id")
	var input request.RejectEventRequest

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid input"})
	}

	if err := config.DB.Model(&models.Event{}).Where("id = ?", id).
		Updates(map[string]interface{}{"status": constant.REJECTED, "remarks": input.Remarks}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to reject event"})
	}

	return c.JSON(fiber.Map{"message": "Event rejected successfully"})
}
