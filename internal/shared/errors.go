package shared

import (
	"errors"
	"net/http"
)

var (
	ErrSystemTemplateNotPermitted = errors.New("system templates cannot be created via API")
	ErrRequiredFieldName          = errors.New("name is required")
	ErrRequiredFieldChannel       = errors.New("channel is required")
	ErrRequiredFieldBody          = errors.New("body is required")
	ErrRequiredFieldSubject       = errors.New("subject is required for email channel")
	ErrTemplateNotFound           = errors.New("template not found")
	ErrDuplicateTemplateRecord    = errors.New("duplicate template record")
	ErrInvalidRecipient           = errors.New("invalid recipient, please check the format")
	ErrInvalidTemplateKeyValue    = errors.New("invalid template_key_value, please check the format")
	ErrRecordNotFound             = errors.New("record not found")
)

func ErrorHttpMapper(err error) int {
	switch err {
	case ErrRequiredFieldBody, ErrRequiredFieldSubject,
		ErrInvalidRecipient, ErrInvalidTemplateKeyValue,
		ErrRequiredFieldChannel, ErrRequiredFieldName, ErrTemplateNotFound:
		return http.StatusBadRequest
	case ErrSystemTemplateNotPermitted:
		return http.StatusForbidden
	case ErrRecordNotFound:
		return http.StatusNotFound
	case ErrDuplicateTemplateRecord:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
