package services

import (
	"fmt"
	"shrine/enums"
	"shrine/messages"
	"shrine/models"
	"shrine/repositories"
	"shrine/types/audit"
	"shrine/types/council"
	"shrine/types/hypertext"
	"shrine/types/user"
	"shrine/utils/auth"
	"shrine/utils/meta"
	"shrine/utils/crypto"
	"shrine/utils/emails"
	"shrine/utils/logger"
	"shrine/utils/sanitize"
	"shrine/utils/validators"
	"strings"
	"time"
)

func ListUsers(pagination meta.Pagination, search string) ([]user.AdminUserResponse, int64) {
	citizens, total := repositories.ListUsers(pagination, search)

	items := make([]user.AdminUserResponse, len(citizens))
	for index, citizen := range citizens {
		items[index] = citizen.ToAdminResponse()
	}

	return items, total
}

func GetUser(username string) (*user.AdminUserResponse, *hypertext.ServiceError) {
	citizen, err := repositories.FindUserByUsername(username)
	if err != nil {
		return nil, fail(enums.NotFound, messages.UserNotFound)
	}

	response := citizen.ToAdminResponse()
	return &response, nil
}

func BanUser(admin *models.User, target *models.User, request council.BanRequest) (*user.AdminUserResponse, *hypertext.ServiceError) {
	if hierarchyErr := auth.ValidateHierarchy(admin, target, "ban"); hierarchyErr != nil {
		return nil, fail(enums.Forbidden, hierarchyErr.Error())
	}

	ref := crypto.SystemRef()
	sanitizedMessage := sanitize.HTML(request.Message)

	now := time.Now()
	target.AccountBanned = true
	target.BannedAt = &now
	target.BannedReason = request.Reason
	target.BannedBy = &admin.ID

	if err := repositories.UpdateUser(target); err != nil {
		return nil, fail(enums.Internal, messages.FailedBanUser)
	}

	if sanitizedMessage != "" {
		repositories.CreateSystemLetter(target.ID, fmt.Sprintf("Account Banned [%s]", ref), sanitizedMessage, ref)
	}

	repositories.LogAction(admin.ID, "user.ban", "user", target.Username, fmt.Sprintf("Banned user @%s", target.Username), audit.BanDetails{
		Reason:    request.Reason,
		SystemRef: ref,
	})

	if err := emails.SendBannedNotification(target.Email, target.Username, request.Reason); err != nil {
		logger.Errorf("Council", "Failed to send ban notification to %s: %v", target.Email, err)
	}

	response := target.ToAdminResponse()
	return &response, nil
}

func UnbanUser(admin *models.User, target *models.User) (*user.AdminUserResponse, *hypertext.ServiceError) {
	target.AccountBanned = false
	target.BannedAt = nil
	target.BannedReason = ""
	target.BannedBy = nil

	if err := repositories.UpdateUser(target); err != nil {
		return nil, fail(enums.Internal, messages.FailedUnbanUser)
	}

	repositories.LogAction(admin.ID, "user.unban", "user", target.Username, fmt.Sprintf("Unbanned user @%s", target.Username), nil)

	response := target.ToAdminResponse()
	return &response, nil
}

func DisableUser(admin *models.User, target *models.User, request council.DisableRequest) (*user.AdminUserResponse, *hypertext.ServiceError) {
	if hierarchyErr := auth.ValidateHierarchy(admin, target, "disable"); hierarchyErr != nil {
		return nil, fail(enums.Forbidden, hierarchyErr.Error())
	}

	ref := crypto.SystemRef()
	sanitizedMessage := sanitize.HTML(request.Message)

	now := time.Now()
	target.AccountDisabled = true
	target.DisabledAt = &now
	target.DisabledReason = request.Reason
	target.DisabledBy = &admin.ID

	var disabledUntilStr string
	if request.DisabledUntil != nil {
		parsed, err := time.Parse(time.RFC3339, *request.DisabledUntil)
		if err != nil {
			return nil, fail(enums.BadRequest, messages.InvalidDisabledUntil)
		}
		target.DisabledUntil = &parsed
		disabledUntilStr = parsed.Format("January 2, 2006")
	}

	if err := repositories.UpdateUser(target); err != nil {
		return nil, fail(enums.Internal, messages.FailedDisableUser)
	}

	if sanitizedMessage != "" {
		repositories.CreateSystemLetter(target.ID, fmt.Sprintf("Account Disabled [%s]", ref), sanitizedMessage, ref)
	}

	repositories.LogAction(admin.ID, "user.disable", "user", target.Username, fmt.Sprintf("Disabled user @%s", target.Username), audit.DisableDetails{
		Reason:        request.Reason,
		DisabledUntil: request.DisabledUntil,
		SystemRef:     ref,
	})

	if err := emails.SendDisabledNotification(target.Email, target.Username, request.Reason, disabledUntilStr); err != nil {
		logger.Errorf("Council", "Failed to send disable notification to %s: %v", target.Email, err)
	}

	response := target.ToAdminResponse()
	return &response, nil
}

func EnableUser(admin *models.User, target *models.User) (*user.AdminUserResponse, *hypertext.ServiceError) {
	target.AccountDisabled = false
	target.DisabledAt = nil
	target.DisabledReason = ""
	target.DisabledBy = nil
	target.DisabledUntil = nil

	if err := repositories.UpdateUser(target); err != nil {
		return nil, fail(enums.Internal, messages.FailedEnableUser)
	}

	repositories.LogAction(admin.ID, "user.enable", "user", target.Username, fmt.Sprintf("Enabled user @%s", target.Username), nil)

	response := target.ToAdminResponse()
	return &response, nil
}

func ChangeRole(admin *models.User, target *models.User, request council.ChangeRoleRequest) (*user.AdminUserResponse, *hypertext.ServiceError) {
	if target.ID == admin.ID {
		return nil, fail(enums.BadRequest, messages.CannotChangeOwnRole)
	}

	if target.IsOwner() {
		return nil, fail(enums.BadRequest, messages.CannotChangeOwnerRole)
	}

	oldRole := string(target.Role)
	role := enums.UserRole(request.Role)

	switch role {
	case enums.Member, enums.Moderator:
		target.Role = role
	case enums.Admin:
		if !admin.IsOwner() {
			return nil, fail(enums.Forbidden, messages.OnlyOwnerAssignAdmin)
		}
		target.Role = role
	default:
		return nil, fail(enums.BadRequest, messages.InvalidRole)
	}

	if err := repositories.UpdateUser(target); err != nil {
		return nil, fail(enums.Internal, messages.FailedChangeRole)
	}

	repositories.LogAction(admin.ID, "user.role_change", "user", target.Username, fmt.Sprintf("Changed role of @%s from %s to %s", target.Username, oldRole, request.Role), audit.RoleChangeDetails{
		OldRole: oldRole,
		NewRole: request.Role,
	})

	response := target.ToAdminResponse()
	return &response, nil
}

func EditUser(admin *models.User, target *models.User, request council.EditUserRequest) (*user.AdminUserResponse, *hypertext.ServiceError) {
	var changes []audit.FieldChange

	if request.Email != nil && !admin.IsOwner() {
		return nil, fail(enums.Forbidden, messages.OnlyOwnerChangeEmail)
	}

	if request.Username != nil && *request.Username != target.Username {
		username := strings.TrimSpace(*request.Username)
		if !validators.IsValidUsername(username, 3) {
			return nil, fail(enums.BadRequest, messages.InvalidUsername)
		}
		if validators.IsReservedUsername(username) {
			return nil, fail(enums.BadRequest, messages.UsernameNotAvailable)
		}
		changes = append(changes, audit.FieldChange{Field: "username", Old: target.Username, New: username})
		target.Username = username
	}

	if request.DisplayName != nil {
		changes = append(changes, audit.FieldChange{Field: "display_name", Old: target.DisplayName, New: *request.DisplayName})
		target.DisplayName = *request.DisplayName
	}

	if request.Email != nil {
		changes = append(changes, audit.FieldChange{Field: "email", Old: target.Email, New: *request.Email})
		target.Email = *request.Email
	}

	if request.Bio != nil {
		target.Bio = *request.Bio
		changes = append(changes, audit.FieldChange{Field: "bio", Old: nil, New: nil})
	}

	if request.Birthday != nil {
		if *request.Birthday == "" {
			target.Birthday = nil
		} else {
			parsed, err := time.Parse("2006-01-02", *request.Birthday)
			if err != nil {
				return nil, fail(enums.BadRequest, "Invalid birthday format. Use YYYY-MM-DD.")
			}
			target.Birthday = &parsed
		}
		changes = append(changes, audit.FieldChange{Field: "birthday", Old: nil, New: nil})
	}

	if request.AvatarURL != nil {
		target.AvatarURL = *request.AvatarURL
		changes = append(changes, audit.FieldChange{Field: "avatar_url", Old: nil, New: nil})
	}

	if request.BlinkieURL != nil {
		target.BlinkieURL = *request.BlinkieURL
		changes = append(changes, audit.FieldChange{Field: "blinkie_url", Old: nil, New: nil})
	}

	if request.Website != nil {
		target.Website = *request.Website
		changes = append(changes, audit.FieldChange{Field: "website", Old: nil, New: nil})
	}

	if request.Location != nil {
		target.Location = *request.Location
		changes = append(changes, audit.FieldChange{Field: "location", Old: nil, New: nil})
	}

	if request.Pronouns != nil {
		target.Pronouns = *request.Pronouns
		changes = append(changes, audit.FieldChange{Field: "pronouns", Old: nil, New: nil})
	}

	if request.Signature != nil {
		target.Signature = *request.Signature
		changes = append(changes, audit.FieldChange{Field: "signature", Old: nil, New: nil})
	}

	if request.Jade != nil {
		changes = append(changes, audit.FieldChange{Field: "jade", Old: target.Jade, New: *request.Jade})
		target.Jade = *request.Jade
	}

	if request.Honor != nil {
		changes = append(changes, audit.FieldChange{Field: "honor", Old: target.Honor, New: *request.Honor})
		target.Honor = *request.Honor
	}

	if len(changes) == 0 {
		return nil, fail(enums.BadRequest, messages.NoChangesProvided)
	}

	if err := repositories.UpdateUser(target); err != nil {
		return nil, fail(enums.Internal, messages.FailedUpdateUser)
	}

	repositories.LogAction(admin.ID, "user.edit", "user", target.Username, fmt.Sprintf("Edited user @%s", target.Username), audit.EditUserDetails{
		Changes: changes,
	})

	response := target.ToAdminResponse()
	return &response, nil
}