package http_response

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/tigapilarmandiri/perkakas/common/constant"
	"github.com/tigapilarmandiri/perkakas/common/util"
)

type HttpResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Meta    *Meta  `json:"meta,omitempty"`
	Debug   *debug `json:"debug,omitempty"`
}

type Meta struct {
	Page      int `json:"page"`
	TotalPage int `json:"total_page"`
	TotalData int `json:"total_data"`
}

type debug struct {
	Error        bool `json:"error"`
	ErrorMessage any  `json:"error_message"`
}

func SendFromNatsResponse(w http.ResponseWriter, statusCode int, message string, data []byte, meta *Meta, debugMessage []byte) {
	response := HttpResponse{
		Status:  statusCode,
		Message: message,
		Data:    json.RawMessage(data),
		Meta:    meta,
	}
	if len(debugMessage) > 0 {
		if debugMessage[0] == '{' || debugMessage[0] == '[' {
			response.Debug = &debug{
				Error:        true,
				ErrorMessage: json.RawMessage(debugMessage),
			}
		} else {
			response.Debug = &debug{
				Error:        true,
				ErrorMessage: string(debugMessage),
			}
		}
	}

	b, err := json.Marshal(response)
	if err != nil {
		util.Log.Error().Msg(err.Error())
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	w.Write(b)
}

func SendSuccess(w http.ResponseWriter, data any, meta *Meta, errorMessage any) {
	successResponse := HttpResponse{
		Status:  200,
		Message: constant.MSG_SUCCESS,
		Data:    data,
		Meta:    meta,
		Debug:   validateErrorMessage(errorMessage),
	}
	res, err := json.Marshal(successResponse)
	if err != nil {
		util.Log.Error().Msg(err.Error())
		SendForbiddenResponse(w, nil)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func SendForbiddenResponse(w http.ResponseWriter, errorMessage any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusForbidden)

	forbiddenResponse := HttpResponse{
		Status:  http.StatusForbidden,
		Message: constant.MSG_FORBIDDEN_ACCESS,
		Data:    nil,
		Debug:   validateErrorMessage(errorMessage),
	}

	res, _ := json.Marshal(forbiddenResponse)
	w.Write(res)
}

func SendNotFoundResponse(w http.ResponseWriter, errorMessage any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)

	notFoundResponse := HttpResponse{
		Status:  http.StatusNotFound,
		Message: constant.MSG_NOT_FOUND,
		Data:    nil,
		Debug:   validateErrorMessage(errorMessage),
	}

	res, _ := json.Marshal(notFoundResponse)
	w.Write(res)
}

func SendRedirectResponse(w http.ResponseWriter, errorMessage any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusPermanentRedirect)

	forbiddenResponse := HttpResponse{
		Status:  http.StatusForbidden,
		Message: constant.MSG_PERMANENTLY_REDIRECT,
		Data:    nil,
		Debug:   validateErrorMessage(errorMessage),
	}

	res, _ := json.Marshal(forbiddenResponse)
	w.Write(res)
}

func validateErrorMessage(errorMessage any) *debug {
	if errorMessage == nil {
		return nil
	}

	msgs, ok := errorMessage.(validator.ValidationErrors)
	if ok {
		return &debug{
			Error:        true,
			ErrorMessage: validatorErrors(msgs),
		}
	}

	err, ok := errorMessage.(error)
	if ok {
		errorMessage = err.Error()
	}

	return &debug{
		Error:        true,
		ErrorMessage: errorMessage,
	}
}

func validatorErrors(errs validator.ValidationErrors) []string {
	var res []string

	for _, v := range errs {
		translate := v.Translate(util.Trans)
		res = append(res, translate)
	}
	return res
}
