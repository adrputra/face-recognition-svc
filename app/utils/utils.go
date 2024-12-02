package utils

import (
	"context"
	"errors"
	"face-recognition-svc/app/model"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/microcosm-cc/bluemonday"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
)

func LogError(c echo.Context, err error, stack []byte) error {
	logrus.Printf("Panic recovered: %v\nStack trace:\n%s\n", err, stack)
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

	if t, ok := md["username"]; ok {
		metaData.Username = sanitizer(t[0])
	}

	if t, ok := md["role_id"]; ok {
		metaData.RoleID = sanitizer(t[0])
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
