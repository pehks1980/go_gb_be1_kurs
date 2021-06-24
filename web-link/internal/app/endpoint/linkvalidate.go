package endpoint

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"net/http"
)

// GetUserStorageKeys - get all keys for this user in repo
func GetUserStorageKeys(request *http.Request, linkSvc linkSvc) ([]string, string, error) {
	props, _ := request.Context().Value("props").(jwt.MapClaims)
	//fmt.Println(props["uid"])
	UID := fmt.Sprintf("%v", props["uid"])
	storageKeys, err := linkSvc.List(UID)
	return storageKeys, UID, err
}

// ValidateRequestShortLink - валидация shortlink + токен авторизации
// ошибка валидации вылетает клиенту сама а результат - true ok
// Эту ф можно применить там где есть этот парам
func ValidateRequestShortLink(request *http.Request, linkSvc linkSvc) (string, string, bool) {

	storageKeys, UID, err := GetUserStorageKeys(request, linkSvc)
	if err != nil {
		return "", "", false
	}

	params := mux.Vars(request)
	shortUrl := params["shortlink"]

	for _, storageKey := range storageKeys {
		if storageKey == shortUrl {
			return UID, storageKey, true
		}
	}

	return "", "", false
}

