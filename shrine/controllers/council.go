package controllers

import (
	"errors"
	"shrine/enums"
	"shrine/models"
	"shrine/repositories"
	"shrine/types"
	"shrine/utils/auth"
	"shrine/utils/meta"
	"time"

	"github.com/gofiber/fiber/v2"
)

func findTargetUser(context *fiber.Ctx) (*models.User, error) {
	username := meta.Request(context).MustHave().Param("username")
	return repositories.FindUserByUsername(username)
}

func ListUsersController(context *fiber.Ctx) error {
	p := meta.Paginate(context)
	search, _ := meta.Request(context).Query("search")

	users, total := repositories.ListUsers(p, search)

	items := make([]types.AdminUserResponse, len(users))
	for i, u := range users {
		items[i] = u.ToAdminResponse()
	}

	return Success(context, p.Response(items, total))
}

func GetUserController(context *fiber.Ctx) error {
	user, err := findTargetUser(context)
	if err != nil {
		return NotFound(context, errors.New("User not found."))
	}

	return Success(context, user.ToAdminResponse())
}

func BanUserController(context *fiber.Ctx) error {
	admin := auth.GetUser(context)

	user, err := findTargetUser(context)
	if err != nil {
		return NotFound(context, errors.New("User not found."))
	}

	if user.ID == admin.ID {
		return BadRequest(context, errors.New("You cannot ban yourself."))
	}

	if user.IsOwner() {
		return BadRequest(context, errors.New("You cannot ban the owner."))
	}

	if user.IsAdmin() && !admin.IsOwner() {
		return BadRequest(context, errors.New("Only the owner can ban an administrator."))
	}

	body, _ := meta.Body[types.BanUserRequest](context)

	now := time.Now()
	user.AccountBanned = true
	user.BannedAt = &now
	user.BannedReason = body.Reason
	user.BannedBy = &admin.ID

	if err := repositories.UpdateUser(user); err != nil {
		return InternalServerError(context, errors.New("Failed to ban user."))
	}

	return Success(context, user.ToAdminResponse())
}

func UnbanUserController(context *fiber.Ctx) error {
	user, err := findTargetUser(context)
	if err != nil {
		return NotFound(context, errors.New("User not found."))
	}

	user.AccountBanned = false
	user.BannedAt = nil
	user.BannedReason = ""
	user.BannedBy = nil

	if err := repositories.UpdateUser(user); err != nil {
		return InternalServerError(context, errors.New("Failed to unban user."))
	}

	return Success(context, user.ToAdminResponse())
}

func DisableUserController(context *fiber.Ctx) error {
	admin := auth.GetUser(context)

	user, err := findTargetUser(context)
	if err != nil {
		return NotFound(context, errors.New("User not found."))
	}

	if user.ID == admin.ID {
		return BadRequest(context, errors.New("You cannot disable yourself."))
	}

	if user.IsOwner() {
		return BadRequest(context, errors.New("You cannot disable the owner."))
	}

	if user.IsAdmin() && !admin.IsOwner() {
		return BadRequest(context, errors.New("Only the owner can disable an administrator."))
	}

	body, _ := meta.Body[types.DisableUserRequest](context)

	now := time.Now()
	user.AccountDisabled = true
	user.DisabledAt = &now
	user.DisabledReason = body.Reason
	user.DisabledBy = &admin.ID

	if err := repositories.UpdateUser(user); err != nil {
		return InternalServerError(context, errors.New("Failed to disable user."))
	}

	return Success(context, user.ToAdminResponse())
}

func EnableUserController(context *fiber.Ctx) error {
	user, err := findTargetUser(context)
	if err != nil {
		return NotFound(context, errors.New("User not found."))
	}

	user.AccountDisabled = false
	user.DisabledAt = nil
	user.DisabledReason = ""
	user.DisabledBy = nil

	if err := repositories.UpdateUser(user); err != nil {
		return InternalServerError(context, errors.New("Failed to enable user."))
	}

	return Success(context, user.ToAdminResponse())
}

func ChangeRoleController(context *fiber.Ctx) error {
	admin := auth.GetUser(context)

	user, err := findTargetUser(context)
	if err != nil {
		return NotFound(context, errors.New("User not found."))
	}

	if user.ID == admin.ID {
		return BadRequest(context, errors.New("You cannot change your own role."))
	}

	if user.IsOwner() {
		return BadRequest(context, errors.New("You cannot change the owner's role."))
	}

	body, err := meta.Body[types.ChangeRoleRequest](context)
	if err != nil {
		return BadRequest(context, errors.New("Invalid request body."))
	}

	role := enums.UserRole(body.Role)
	switch role {
	case enums.Member, enums.Moderator:
		user.Role = role
	case enums.Admin:
		if !admin.IsOwner() {
			return BadRequest(context, errors.New("Only the owner can assign the admin role."))
		}
		user.Role = role
	default:
		return BadRequest(context, errors.New("Invalid role."))
	}

	if err := repositories.UpdateUser(user); err != nil {
		return InternalServerError(context, errors.New("Failed to change role."))
	}

	return Success(context, user.ToAdminResponse())
}