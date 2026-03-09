package files

import (
	"shrine/utils/collections"
	"strings"
)

var archiveMimes = collections.SetOf(
	"application/zip",
	"application/x-rar-compressed",
	"application/gzip",
	"application/x-7z-compressed",
	"application/x-tar",
	"application/x-bzip2",
	"application/x-xz",
	"application/x-compress",
)

var documentMimes = collections.SetOf(
	"application/msword",
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	"application/vnd.oasis.opendocument.text",
	"application/rtf",
	"application/epub+zip",
)

var spreadsheetMimes = collections.SetOf(
	"application/vnd.ms-excel",
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	"application/vnd.oasis.opendocument.spreadsheet",
	"text/csv",
)

var presentationMimes = collections.SetOf(
	"application/vnd.ms-powerpoint",
	"application/vnd.openxmlformats-officedocument.presentationml.presentation",
	"application/vnd.oasis.opendocument.presentation",
)

var codeMimes = collections.SetOf(
	"application/javascript",
	"application/json",
	"application/xml",
	"application/x-httpd-php",
	"application/x-sh",
	"application/x-python",
	"application/typescript",
)

var databaseMimes = collections.SetOf(
	"application/x-sqlite3",
	"application/vnd.ms-access",
)


func DetectCategory(contentType string) string {
	if strings.HasPrefix(contentType, "image/") {
		return "image"
	}
	if strings.HasPrefix(contentType, "video/") {
		return "video"
	}
	if strings.HasPrefix(contentType, "audio/") {
		return "audio"
	}
	if strings.Contains(contentType, "font") {
		return "font"
	}
	if contentType == "application/pdf" {
		return "document"
	}
	if archiveMimes.Has(contentType) {
		return "archive"
	}
	if documentMimes.Has(contentType) {
		return "document"
	}
	if spreadsheetMimes.Has(contentType) {
		return "spreadsheet"
	}
	if presentationMimes.Has(contentType) {
		return "presentation"
	}
	if codeMimes.Has(contentType) {
		return "code"
	}
	if databaseMimes.Has(contentType) {
		return "database"
	}
	if strings.HasPrefix(contentType, "text/x-") {
		return "code"
	}
	if strings.HasPrefix(contentType, "text/") {
		return "document"
	}
	return "other"
}