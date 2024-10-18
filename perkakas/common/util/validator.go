package util

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"sync"

	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	id_translations "github.com/go-playground/validator/v10/translations/id"
)

var RoleUserPermission = map[string][]string{
	"POLDA":  {"UMUM", "POLDA", "POLRES", "POLSEK"},
	"POLRES": {"UMUM", "POLRES", "POLSEK"},
	"POLSEK": {"UMUM", "POLSEK"},
}

var (
	internalValidator *validator.Validate
	onceValidator     sync.Once
	Trans             ut.Translator
)

var (
	ErrJsonInvalid = errors.New("json invalid")
	ErrValidation  = errors.New("error validation")
)

func init() {
	onceValidator.Do(func() {
		id_locales := id.New()
		universalTranslator := ut.New(id_locales, id_locales)
		Trans, _ = universalTranslator.GetTranslator("id")
		internalValidator = validator.New()
		err := id_translations.RegisterDefaultTranslations(internalValidator, Trans)
		if err != nil {
			Log.Fatal().Msg(err.Error())
		}
	})
}

func ValidateStruct(ctx context.Context, i any) error {
	return internalValidator.StructCtx(ctx, i)
}

type helperService struct {
	errs []string
}

func (hs helperService) Error() string {
	return strings.Join(hs.errs, ", ")
}

func ServiceValidateStruct(ctx context.Context, i any) error {
	var hs helperService

	err := internalValidator.StructCtx(ctx, i)
	if err == nil {
		return nil
	}
	errs := err.(validator.ValidationErrors)

	for _, v := range errs {
		translate := v.Translate(Trans)
		hs.errs = append(hs.errs, translate)
	}
	return hs
}

type ErrorStruct struct {
	Status  int
	Message string
	Err     error
}

func ValidateAndUnmarshal[T any](ctx context.Context, data []byte, model *T) *ErrorStruct {
	if err := json.Unmarshal(data, model); err != nil {
		return &ErrorStruct{
			Status:  http.StatusUnauthorized,
			Message: ErrJsonInvalid.Error(),
			Err:     ErrJsonInvalid,
		}
	}

	if err := ServiceValidateStruct(ctx, *model); err != nil {
		return &ErrorStruct{
			Status:  http.StatusBadRequest,
			Message: ErrValidation.Error(),
			Err:     err,
		}
	}

	return nil
}
