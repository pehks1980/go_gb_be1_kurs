package endpoint

import (
	"encoding/json"
	"net/http"
)

// json error api
type Errors struct {
	Errors []ErrorEl `json:"errors"`
}

// Error json array
type ErrorEl struct {
	Code    uint64  `json:"code"`
	Message string `json:"message"`
}

var (
	APIErrorList = map[uint64]string {
		1: "Token has expired time",    //# ошибка поиска токена в headers
		2: "Unknown token",                 //# не авторизован на сервере
		3: "Access token is lost",         // # не передан токен в headers
		4: "The shortlink name must be provided.",
		5: "The shortlink with the specified name already exists",
		6: "The shortlink with the specified name does not exist",
		7: "Please provide refresh token, or authenticate again",
		8: "No uid (user id), please set uid",
		9: "Unknown content type",
		10: "Internal repo problem",
		400: "Bad request",
		401: "Unauthorized",
		402: "Payment required",
		403: "Forbidden",
		404: "Not found",
		405: "Method not allowed",
	}
)

func ResponseApiError(w http.ResponseWriter, code uint64, status int) {
	var errorsjson = Errors{}
	errorel := ErrorEl{Code: code, Message: APIErrorList[code]}
	errorsjson.Errors = append(errorsjson.Errors, errorel)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(errorsjson)
	w.WriteHeader(status)
}
