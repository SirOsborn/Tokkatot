package api

import (
	"middleware/services"
	"middleware/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var (
	farmService      = services.NewFarmService()
	coopService      = services.NewCoopService()
	alertService     = services.NewAlertService()
	deviceService    = services.NewDeviceService()
	authService      = services.NewAuthService()
	scheduleService  = services.NewScheduleService()
	analyticsService = services.NewAnalyticsService()
	adminService     = services.NewAdminService()
)

// checkFarmAccess is a helper to verify farm membership/role
func checkFarmAccess(userID, farmID uuid.UUID, minRole string) error {
	err := farmService.CheckAccess(userID, farmID, minRole)
	if err == services.ErrFarmAccessDenied {
		return utils.Forbidden(nil, "Access denied")
	}
	return err
}

// verifyFarmAccess version that takes context
func verifyFarmAccess(c *fiber.Ctx, userID, farmID uuid.UUID, minRole string) error {
	err := farmService.CheckAccess(userID, farmID, minRole)
	if err == services.ErrFarmAccessDenied {
		return utils.Forbidden(c, "You do not have permission to access this farm")
	}
	if err != nil {
		return utils.InternalError(c, "Failed to verify farm access")
	}
	return nil
}
