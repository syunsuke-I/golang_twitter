package models

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Page struct {
	Number        int `json:"number"`
	Size          int `json:"size"`
	TotalElements int `json:"total_elements"`
	TotalPages    int `json:"total_pages"`
}

type Pagination struct{}

func (pagination Pagination) Pagination(page Page) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page.Number - 1) * page.Size
		return db.Offset(offset).Limit(page.Size)
	}
}

func ConvertContextAndTotalElementsToPage(context *gin.Context, totalElements int) Page {
	page, _ := strconv.Atoi(context.Query("page"))
	if page == 0 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(context.Query("size"))
	switch {
	case pageSize > totalElements:
		pageSize = totalElements
	case pageSize > 100:
		pageSize = 100
	case pageSize <= 0:
		if totalElements < 5 {
			pageSize = totalElements
		} else {
			pageSize = 5
		}
	}
	totalPages := int(math.Ceil(float64(totalElements) / float64(pageSize)))

	return Page{Number: page, Size: pageSize, TotalElements: totalElements, TotalPages: totalPages}
}
