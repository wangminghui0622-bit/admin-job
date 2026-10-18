// Package dberrorhandler maps gorm errors to unified admin-zero business errors.
package dberrorhandler

import (
	"errors"
	"strings"

	"github.com/wangminghui0622-bit/admin-common/errorx"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

// DefaultError converts a gorm error into an admin-zero error.
func DefaultError(logger logx.Logger, err error, detail any) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		logger.Errorw(err.Error(), logx.Field("detail", detail))
		return errorx.New(errorx.InvalidArgument, "target not exist")
	case isConstraintError(err):
		logger.Errorw(err.Error(), logx.Field("detail", detail))
		return errorx.New(errorx.InvalidArgument, "constraint error")
	default:
		logger.Errorw("database error", logx.Field("detail", err.Error()))
		return errorx.New(errorx.ServerError, "database error")
	}
}

func isConstraintError(err error) bool {
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "duplicate") ||
		strings.Contains(msg, "constraint") ||
		strings.Contains(msg, "cannot be null") ||
		strings.Contains(msg, "doesn't have a default value")
}
