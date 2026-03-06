package meta

import (
	"shrine/types"
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

func (p Pagination) Apply(query *gorm.DB) *gorm.DB {
	return query.Offset((p.Page - 1) * p.PerPage).Limit(p.PerPage)
}

func (p Pagination) Response(items any, total int64) types.PaginatedResponse {
	totalPages := int(total) / p.PerPage
	if int(total)%p.PerPage > 0 {
		totalPages++
	}

	return types.PaginatedResponse{
		Items:      items,
		Total:      total,
		Page:       p.Page,
		PerPage:    p.PerPage,
		TotalPages: totalPages,
	}
}