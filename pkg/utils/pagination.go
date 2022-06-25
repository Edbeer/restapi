package utils

import (
	"math"
	"strconv"

	"github.com/labstack/echo/v4"
)

const (
	defaultSize = 10
)

// Pagination with out offset
type PaginationQuery struct {
	Size       int    `json:"size,omitempty"`
	Page       int    `json:"page,omitempty"`
	Difference int    `json:"difference,omitempty"`
	OrderBy    string `json:"orderBy,omitempty"`
}

// Get pagination query struct from
func GetPaginationFromCtx(c echo.Context) (*PaginationQuery, error) {
	q := &PaginationQuery{}
	if err := q.SetPage(c.QueryParam("page")); err != nil {
		return nil, err
	}
	if err := q.SetSize(c.QueryParam("size")); err != nil {
		return nil, err
	}

	return q, nil
}

// Get total pages int
func GetTotalPages(totalCount int, pageSize int) int {
	d := float64(totalCount) / float64(pageSize)
	return int(math.Ceil(d))
}

// Set page size
func (q *PaginationQuery) SetSize(sizeQuery string) error {
	if sizeQuery == "" {
		q.Size = defaultSize
		return nil
	}
	n, err := strconv.Atoi(sizeQuery)
	if err != nil {
		return err
	}
	q.Size = n

	return nil
}

// Set page number
func (q *PaginationQuery) SetPage(pageQuery string) error {
	if pageQuery == "" {
		q.Size = 0
		return nil
	}
	n, err := strconv.Atoi(pageQuery)
	if err != nil {
		return err
	}
	q.Page = n

	return nil
}

// Set difference
func (q *PaginationQuery) SetDifference(differenceQuery string) error {
	if differenceQuery == "" {
		q.Difference = defaultSize
		return nil
	}
	n, err := strconv.Atoi(differenceQuery)
	if err != nil {
		return err
	}
	q.Difference = n

	return nil
}

// Set order by
func (q *PaginationQuery) SetOrderBy(orderByQuery string) {
	q.OrderBy = orderByQuery
}

// Get limit
func (q *PaginationQuery) GetLimit() int {
	return q.Size
}

// Get page
func (q *PaginationQuery) GetPage() int {
	return q.Page
}

// Get size
func (q *PaginationQuery) GetSize() int {
	return q.Size
}

// Get difference
func (q *PaginationQuery) GetDifference() int {
	return q.Difference
}

// Get OrderBy
func (q *PaginationQuery) GetOrderBy() string {
	return q.OrderBy
}