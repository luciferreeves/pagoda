package meta

import (
	"shrine/types/common"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Pagination struct {
	Page    int
	PerPage int
}

func Paginate(context *fiber.Ctx) Pagination {
	request := Request(context)

	pageStr, _ := request.Query("page")
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}

	perPageStr, _ := request.Query("per_page")
	perPage, _ := strconv.Atoi(perPageStr)
	if perPage < 1 || perPage > 50 {
		perPage = 20
	}

	return Pagination{Page: page, PerPage: perPage}
}

func (self Pagination) Apply(query *gorm.DB) *gorm.DB {
	return query.Offset((self.Page - 1) * self.PerPage).Limit(self.PerPage)
}

func (self Pagination) Response(items any, total int64) common.PaginatedResponse {
	totalPages := int(total) / self.PerPage
	if int(total)%self.PerPage > 0 {
		totalPages++
	}

	return common.PaginatedResponse{
		Items:      items,
		Total:      total,
		Page:       self.Page,
		PerPage:    self.PerPage,
		TotalPages: totalPages,
	}
}