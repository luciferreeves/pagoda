package validators

import (
	"regexp"
	"shrine/utils/collections"
	"strings"
)

var validUsernamePattern = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

var reservedUsernames = collections.SetOf(
	"admin", "administrator", "owner",
	"mod", "moderator", "janitor", "staff",
	"api", "www", "mail", "email",
	"support", "help", "helpdesk",
	"about", "contact", "privacy", "terms", "tos",
	"null", "undefined", "system", "bot", "guest",
	"login", "register", "signup", "signin", "logout",
	"profile", "settings", "account", "dashboard",
	"shifoo", "pagoda", "deleted", "anonymous",
	"root", "webmaster", "postmaster", "hostmaster",
	"noreply", "no_reply", "mailer", "daemon",
	"abuse", "security", "info", "status",
	"home", "search", "explore", "discover",
	"forums", "forum", "chat", "bazaar", "districts",
	"members", "online", "buttons", "webring",
	"guestbook", "hitcounter", "letters",
	"rules", "faq", "blog", "news", "feed",
	"static", "assets", "uploads", "images", "media", "default", "defaults",
	"test", "testing", "debug", "dev", "staging",
	"everyone", "all", "here", "channel",
)

func IsValidUsername(username string, minimumLength int) bool {
	if len(username) < minimumLength || len(username) > 32 {
		return false
	}
	return validUsernamePattern.MatchString(username)
}

func IsReservedUsername(username string) bool {
	return reservedUsernames.Has(strings.ToLower(username))
}