package utils

import (
	"context"
	"errors"
	"face-recognition-svc/gateway/app/model"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/microcosm-cc/bluemonday"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/metadata"
)

func LogError(c echo.Context, err error, stack []byte) error {
	log.Error().Err(err).Bytes("stack", stack).Msg("Error occurred")
	data, ok := err.(*model.ErrorResponse)
	if ok {
		return c.JSON(data.Code, model.Response{
			Code:    data.Code,
			Message: err.Error(),
			Data:    nil,
		})
	}
	return c.JSON(http.StatusInternalServerError, model.Response{
		Code:    http.StatusInternalServerError,
		Message: err.Error(),
		Data:    nil,
	})
}

func GetMetadata(c context.Context) (*model.MetadataUser, error) {
	var metaData = &model.MetadataUser{}
	md, ok := metadata.FromIncomingContext(c)
	if !ok {
		return nil, errors.New("Error")
	}

	if t, ok := md["user_id"]; ok {
		metaData.UserID = sanitizer(t[0])
	}

	if t, ok := md["username"]; ok {
		metaData.Username = sanitizer(t[0])
	}

	if t, ok := md["role_ids"]; ok {
		raw := sanitizer(t[0])
		if raw != "" {
			metaData.RoleIDs = strings.Split(raw, ",")
		}
	}

	if t, ok := md["role_id"]; ok {
		roleID := sanitizer(t[0])
		if roleID != "" && len(metaData.RoleIDs) == 0 {
			metaData.RoleIDs = []string{roleID}
		}
	}

	if t, ok := md["institution_id"]; ok {
		metaData.InstitutionID = sanitizer(t[0])
	}

	return metaData, nil
}

var sanitize = bluemonday.NewPolicy()

func sanitizer(s string) string {
	const replacement = ""

	var replacer = strings.NewReplacer(
		"\r\n", replacement,
		"\r", replacement,
		"\n", replacement,
		"\v", replacement,
		"\f", replacement,
		"\u0085", replacement,
		"\u2028", replacement,
		"\u2029", replacement,
	)
	out := replacer.Replace(s)
	return sanitize.Sanitize(out)
}

func Contains(arr []string, str string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}

var jakartaLoc *time.Location

func InitTimeLocation() {
	var err error
	jakartaLoc, err = time.LoadLocation("Asia/Jakarta")
	if err != nil {
		fmt.Println("Warning: Failed to load Jakarta timezone, using fixed UTC+7 instead")
		jakartaLoc = time.FixedZone("WIB", 7*60*60) // UTC+7 for Jakarta
	}
}

func LocalTime() time.Time {
	return time.Now().In(jakartaLoc)
}

// ParseFilterFromQuery parses filter parameters from echo context query parameters
func ParseFilterFromQuery(c echo.Context) *model.Filter {
	search := c.QueryParam("search")
	sortBy := c.QueryParam("sort_by")
	sortOrder := c.QueryParam("sort_order")

	// Set defaults
	if sortBy == "" {
		sortBy = "id"
	}

	if sortOrder == "" {
		sortOrder = "ASC"
	} else if sortOrder != "ASC" && sortOrder != "DESC" {
		sortOrder = "ASC"
	}

	return &model.Filter{
		Search:    search,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}
}

// ParsePaginationFromQuery parses pagination parameters from echo context query parameters
func ParsePaginationFromQuery(c echo.Context) *model.Pagination {
	pageParam := c.QueryParam("page")
	limitParam := c.QueryParam("limit")

	var (
		page  int = 1
		limit int = 10
	)

	if pageParam != "" {
		if p, err := strconv.Atoi(pageParam); err == nil && p > 0 {
			page = p
		}
	}

	if limitParam != "" {
		if l, err := strconv.Atoi(limitParam); err == nil && l > 0 {
			limit = l
		}
	}

	return &model.Pagination{
		Page:  page,
		Limit: limit,
	}
}

// BuildSearchWhereClause builds a WHERE clause for search with multiple fields
// searchFields should be a slice of column names to search in
func BuildSearchWhereClause(search string, searchFields []string) string {
	if search == "" || len(searchFields) == 0 {
		return ""
	}

	searchEscaped := strings.ReplaceAll(search, "'", "''")
	var conditions []string
	for _, field := range searchFields {
		conditions = append(conditions, fmt.Sprintf("%s ILIKE '%%%s%%'", field, searchEscaped))
	}

	return " WHERE (" + strings.Join(conditions, " OR ") + ")"
}

// BuildOrderByClause builds an ORDER BY clause with validation
// allowedFields is a map of allowed field names to their actual column names (can include table aliases)
// defaultField is the default field to sort by if sortBy is not in allowedFields
func BuildOrderByClause(filter *model.Filter, allowedFields map[string]string, defaultField string) string {
	sortBy := filter.SortBy
	sortOrder := filter.SortOrder

	// Validate sortBy
	actualField, ok := allowedFields[sortBy]
	if !ok {
		actualField = defaultField
	}

	// Validate sortOrder
	if sortOrder != "ASC" && sortOrder != "DESC" {
		sortOrder = "ASC"
	}

	return fmt.Sprintf(" ORDER BY %s %s", actualField, sortOrder)
}
