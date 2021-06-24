package endpoint

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/pehks1980/go_gb_be1_kurs/web-link/internal/pkg/model"
)

// интерфейс очередного сервиса также имеет put get - для работы с файлохранилищем
// у "драйвера хранилища" методы
type linkSvc interface {
	Get(uid, key string) (model.DataEl, error)
	Put(uid, key string, value model.DataEl) error
	Del(uid, key string) error
	List(uid string) ([]string, error)
	GetUn(shortlink string) (model.DataEl, error)
}

// RegisterPublicHTTP - регистрация роутинга путей типа urls.py для обработки сервером
func RegisterPublicHTTP(linkSvc linkSvc) *mux.Router {
	// mux golrilla почему он? не знаю, - прикольное название, простота работы..
	r := mux.NewRouter()
	// JWT authorization
	r.HandleFunc("/user/auth", postAuth(linkSvc)).Methods(http.MethodPost)
	r.HandleFunc("/token/refresh", postTokenRefresh(linkSvc)).Methods(http.MethodPost)
	// main function
	r.HandleFunc("/shortopen/{shortlink}", getShortOpen(linkSvc)).Methods(http.MethodGet)
	r.HandleFunc("/shortstat/{shortlink}", getShortStat(linkSvc)).Methods(http.MethodGet)
	// links crud
	r.HandleFunc("/links", postToLink(linkSvc)).Methods(http.MethodPost)
	r.HandleFunc("/links/all", getFromLink(linkSvc)).Methods(http.MethodGet)
	r.HandleFunc("/links/{shortlink}", putToLink(linkSvc)).Methods(http.MethodPut)
	r.HandleFunc("/links/{shortlink}", delFromLink(linkSvc)).Methods(http.MethodDelete)

	r.Use(LoggingMiddleware)
	r.Use(JWTCheckMiddleware)
	return r
}

// postTokenRefresh - get new pair of jwt tokens when access token is expired
func postTokenRefresh(svc linkSvc) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {
		type TokenAnswer struct {
			Access  string `json:"accessToken"`
			Refresh string `json:"refreshToken"`
		}

		props, _ := request.Context().Value("props").(jwt.MapClaims)
		//fmt.Println(props["uid"])
		UID := fmt.Sprintf("%v", props["uid"])
		Issuer := fmt.Sprintf("%v", props["iss"])

		if Issuer != "weblink_refresh" {
			ResponseApiError(w, 7, http.StatusBadRequest)
			return
		}

		token_access, _ := GenJWTWithClaims(UID, 0)
		token_refresh, _ := GenJWTWithClaims(UID, 1)

		var jsonTokens = TokenAnswer{
			Access:  token_access,
			Refresh: token_refresh,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(jsonTokens)

		if err != nil {
			return
		}

	}
}

// postAuth - autheticate and give authorization token
func postAuth(svc linkSvc) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {

		type TokenAnswer struct {
			Access  string `json:"accessToken"`
			Refresh string `json:"refreshToken"`
		}

		type PostJsonRq struct {
			Uid string `json:"uid"`
		}

		var jsonPostRq = PostJsonRq{}

		err := json.NewDecoder(request.Body).Decode(&jsonPostRq)
		if err != nil {
			ResponseApiError(w, 400, http.StatusBadRequest)
			return
		}
		// get uid
		if jsonPostRq.Uid == "" {
			ResponseApiError(w, 8, http.StatusBadRequest)
			return
		}

		token_access, _ := GenJWTWithClaims(jsonPostRq.Uid, 0)
		token_refresh, _ := GenJWTWithClaims(jsonPostRq.Uid, 1)

		var jsonTokens = TokenAnswer{
			Access:  token_access,
			Refresh: token_refresh,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(jsonTokens)
		if err != nil {
			return
		}

		return

	}
}

// delFromLink deletes link from api storage by shortlink
func delFromLink(linkSvc linkSvc) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {

		UID, storageKey, res := ValidateRequestShortLink(request, linkSvc)
		if !res {
			ResponseApiError(w, 4, http.StatusBadRequest)
			return
		}

		//found key, delete it
		err := linkSvc.Del(UID, storageKey)
		if err != nil {
			ResponseApiError(w, 10, http.StatusBadRequest)
			return
		}

	}
}

// putToLink updates link from api storage
func putToLink(linkSvc linkSvc) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		var element = model.DataEl{}
		w.Header().Set("Content-Type", "application/json")

		UID, _, res := ValidateRequestShortLink(request, linkSvc)
		if !res {
			ResponseApiError(w, 9, http.StatusBadRequest)
			return
		}
		//found key, work with body
		err := json.NewDecoder(request.Body).Decode(&element)
		if err != nil {
			ResponseApiError(w, 9, http.StatusBadRequest)
			return
		}
		element.Datetime = time.Now()
		element.UID = UID
		element.Active = 1
		//looks ok, update storage
		err = linkSvc.Put(UID, element.Shorturl, element)
		if err != nil {
			ResponseApiError(w, 9, http.StatusBadRequest)
		}
		// form answer json
		err = json.NewEncoder(w).Encode(element)
		if err != nil {
			return
		}
		return
	}
}

// postToLink - creates new item in api storage
func postToLink(linkSvc linkSvc) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		var element = model.DataEl{}
		w.Header().Set("Content-Type", "application/json")

		storageKeys, UID, err := GetUserStorageKeys(request, linkSvc)
		if err != nil {
			ResponseApiError(w, 10, http.StatusBadRequest)
			return
		}

		err = json.NewDecoder(request.Body).Decode(&element)
		if err != nil {
			ResponseApiError(w, 9, http.StatusBadRequest)
			return
		}
		// check if we have key
		if element.Shorturl == "" {
			ResponseApiError(w, 11, http.StatusBadRequest)
			return
		}

		element.Datetime = time.Now()
		// check if this key already exists
		for _, storageKey := range storageKeys {
			if storageKey == element.Shorturl {
				ResponseApiError(w, 5, http.StatusBadRequest)
				return
			}
		}
		element.UID = UID
		element.Active = 1
		err = linkSvc.Put(UID, element.Shorturl, element)
		if err != nil {
			ResponseApiError(w, 10, http.StatusBadGateway)
		}
		w.WriteHeader(http.StatusCreated) // this has to be the first write!!!
		err = json.NewEncoder(w).Encode(element)
		if err != nil {
			return
		}
		return

	}
}

// getFromLink - get links list in json
func getFromLink(linkSvc linkSvc) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var datajson = model.Data{}

		storageKeys, UID, err := GetUserStorageKeys(request, linkSvc)
		if err != nil {
			ResponseApiError(w, 10, http.StatusBadRequest)
			return
		}

		for _, storageKey := range storageKeys {
			getElement, errfor := linkSvc.Get(UID, storageKey)
			if errfor != nil {
				ResponseApiError(w, 10, http.StatusBadGateway)
				return
				//http.Error(w, "Cannot read from repo", http.StatusBadRequest)
			}

			datajson.Data = append(datajson.Data, getElement)
		}

		err = json.NewEncoder(w).Encode(datajson)
		if err != nil {
			return
		}

	}
}

// getShortStat - get one link from api
func getShortStat(linkSvc linkSvc) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var datajson = model.Data{}
		// check user authorization, get user UID, get key (for this user, check if key exists)
		// if res - yes then do the action  - give string from repo as json
		UID, storageKey, res := ValidateRequestShortLink(request, linkSvc)
		if !res {
			ResponseApiError(w, 11, http.StatusBadRequest)
			return
		}

		getElement, err := linkSvc.Get(UID, storageKey)
		if err != nil {
			ResponseApiError(w, 10, http.StatusBadRequest)
			return
		}

		datajson.Data = append(datajson.Data, getElement)
		err = json.NewEncoder(w).Encode(datajson)
		if err != nil {
			return
		}
	}
}

// getShortOpen - get link opened (unonimously)
func getShortOpen(linkSvc linkSvc) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		// get data
		// update data
		// redir to real link

		params := mux.Vars(request)
		shortUrl := params["shortlink"]
		// GetUn retreives link and updates redir count
		getElement, err := linkSvc.GetUn(shortUrl)
		if err != nil {
			ResponseApiError(w, 10, http.StatusBadRequest)
			return
		}

		log.Printf("opening user %s link  %s (short is %s) redirs(++) %d \n", getElement.UID, getElement.URL, getElement.Shorturl, getElement.Redirs)
		http.Redirect(w, request, getElement.URL, http.StatusSeeOther)
		//<a href="/shortopen/www.mail.ru">See Other</a>.
		return

	}
}
