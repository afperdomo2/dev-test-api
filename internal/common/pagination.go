package common

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	DefaultPage      = 1
	DefaultPerPage   = 20
	MaxPerPage       = 100
	DefaultSortOrder = "desc"
)

type PaginationParams struct {
	Page      int
	PerPage   int
	SortBy    string
	SortOrder string
}

type SortConfig struct {
	Allowed []string
	Default string
}

func (cfg SortConfig) OrderClause(sortBy, sortOrder string) string {
	if sortBy == "" {
		return cfg.Default
	}
	return sortBy + " " + sortOrder
}

func ParsePagination(c *gin.Context, cfg SortConfig) (PaginationParams, error) {
	page := DefaultPage
	if p := c.Query("page"); p != "" {
		v, err := strconv.Atoi(p)
		if err != nil || v < 1 {
			return PaginationParams{}, fmt.Errorf("page debe ser un número entero positivo")
		}
		page = v
	}

	perPage := DefaultPerPage
	if pp := c.Query("perPage"); pp != "" {
		v, err := strconv.Atoi(pp)
		if err != nil || v < 1 {
			return PaginationParams{}, fmt.Errorf("perPage debe ser un número entero positivo")
		}
		if v > MaxPerPage {
			return PaginationParams{}, fmt.Errorf("perPage no puede exceder %d", MaxPerPage)
		}
		perPage = v
	}

	sortBy := c.Query("sortBy")
	sortOrder := c.Query("sortOrder")

	if sortBy != "" {
		valid := false
		for _, a := range cfg.Allowed {
			if sortBy == a {
				valid = true
				break
			}
		}
		if !valid {
			return PaginationParams{}, fmt.Errorf("sortBy '%s' no es válido. Valores permitidos: %s",
				sortBy, strings.Join(cfg.Allowed, ", "))
		}
	}

	if sortOrder != "" && sortOrder != "asc" && sortOrder != "desc" {
		return PaginationParams{}, fmt.Errorf("sortOrder debe ser 'asc' o 'desc'")
	}
	if sortBy != "" && sortOrder == "" {
		sortOrder = DefaultSortOrder
	}

	return PaginationParams{
		Page:      page,
		PerPage:   perPage,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}, nil
}
