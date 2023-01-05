package helper

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/nanpipat/golang-template-hexagonal/consts"

	"github.com/gofiber/fiber/v2"
)

type Pagination struct {
	Page  int64       `json:"page" example:"1"`
	Total int64       `json:"total" example:"45"`
	Limit int64       `json:"limit" example:"30"`
	Count int64       `json:"count" example:"30"`
	Items interface{} `json:"items"`
}

type PageOptions struct {
	Q       string
	Limit   int64
	Page    int64
	OrderBy []string
}

func (p *PageOptions) SetOrderDefault(orders ...string) {
	if len(p.OrderBy) == 0 {
		p.OrderBy = orders
	}
}

type PageResponse struct {
	Total   int64
	Limit   int64
	Count   int64
	Page    int64
	Q       string
	OrderBy []string
}

func NewPagination(items interface{}, options *PageResponse) *Pagination {
	m := &Pagination{}
	if options != nil {
		m.Limit = options.Limit
		m.Page = options.Page
		m.Total = options.Total
		m.Count = options.Count
	}

	if items == nil {
		m.Items = make([]interface{}, 0)
	} else {
		m.Items = items
	}

	return m
}

func GetPageOptions(c *fiber.Ctx) *PageOptions {
	limit, _ := strconv.ParseInt(c.Query("limit"), 10, 64)
	page, _ := strconv.ParseInt(c.Query("page"), 10, 64)

	if limit <= 0 {
		limit = consts.PageLimitDefault
	}

	if limit > consts.PageLimitMax {
		limit = consts.PageLimitMax
	}

	if page < 1 {
		page = 1
	}

	return &PageOptions{
		Q:       c.Query("q"),
		Limit:   limit,
		Page:    page,
		OrderBy: genOrderBy(c.Query("order_by")),
	}
}

func genOrderBy(s string) []string {
	orderBy := make([]string, 0)
	fields := strings.Split(s, ",")
	for _, field := range fields {
		spaceParameters := strings.Split(field, " ")
		bracketParameters := strings.Split(field, "(")
		if len(spaceParameters) == 1 && len(bracketParameters) == 1 && spaceParameters[0] != "" {
			orderBy = append(orderBy, fmt.Sprintf("%s desc", spaceParameters[0]))
		} else if len(spaceParameters) == 2 {
			name := spaceParameters[0]
			if name != "" {
				shortingParameter := spaceParameters[1]
				if shortingParameter == "asc" {
					orderBy = append(orderBy, fmt.Sprintf("%s %s", name, shortingParameter))
				} else {
					orderBy = append(orderBy, fmt.Sprintf("%s desc", name))
				}
			}
		} else if len(bracketParameters) == 2 {
			name := strings.TrimSuffix(bracketParameters[1], ")")
			if name != "" {
				shortingParameter := bracketParameters[0]
				if shortingParameter == "asc" {
					orderBy = append(orderBy, fmt.Sprintf("%s %s", name, shortingParameter))
				} else {
					orderBy = append(orderBy, fmt.Sprintf("%s desc", name))
				}
			}
		}
	}
	return orderBy
}
