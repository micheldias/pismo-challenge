package model

import "strings"

const (
	BadRequestErrorCode     = "CODE-400"
	InternalServerErrorCode = "CODE-500"
	EntityNorFoundErrorCode = "CODE-404"
)

type EntityNorFoundError struct {
	Reason string
}

func (e EntityNorFoundError) Error() string {
	if e.Reason != "" {
		return e.Reason
	}
	return "entity not found"
}

type MissingRequiredInformation struct {
	Errors []error
	Reason string
}

func (b MissingRequiredInformation) Error() string {
	if len(b.Errors) > 0 {
		var arrOfErrors []string
		for _, err := range b.Errors {
			arrOfErrors = append(arrOfErrors, err.Error())
		}
		return strings.Join(arrOfErrors, ";")
	}
	if b.Reason != "" {
		return b.Reason
	}
	return "missing required information"
}
