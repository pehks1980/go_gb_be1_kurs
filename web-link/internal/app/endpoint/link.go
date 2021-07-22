package endpoint

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/pehks1980/go_gb_be1_kurs/web-link/internal/pkg/model"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

//  linkSvc - интерфейс сервиса со стороны http
//  также имеет put get del crud - для работы с файлохранилищем
//  list - list all links for uid user
//  GetUn - open link for redir and add 1 to redir count
type linkSvc interface {
	Get(uid, key string, su bool) (model.DataEl, error)
	Put(uid, key string, value model.DataEl, su bool) error
	Del(uid, key string, su bool) error
	List(uid string) ([]string, error)
	GetUn(shortlink string) (string, error)
	PutUser(value model.User) (string, error)
	DelUser(name, email string) error
	GetUser(uid string) (model.User, error)
	WhoAmI() uint64
	PayUser(uidA, uidB, amount string) error
	FindSuperUser() (string, error)
	GetAll() (model.Data, error)
}

// RegisterPublicHTTP - регистрация роутинга путей типа urls.py для обработки сервером
func RegisterPublicHTTP(linkSvc linkSvc) *mux.Router {
	// mux golrilla почему он? не знаю, - прикольное название, простота работы..
	r := mux.NewRouter()
	// JWT authorization
	r.HandleFunc("/user/auth", postAuth(linkSvc)).Methods(http.MethodPost)
	r.HandleFunc("/token/refresh", postTokenRefresh(linkSvc)).Methods(http.MethodPost)
	// Main function
	r.HandleFunc("/shortopen/{shortlink}", getShortOpen(linkSvc)).Methods(http.MethodGet)
	r.HandleFunc("/shortstat/{shortlink}", getShortStat(linkSvc)).Methods(http.MethodGet)
	// Links crud
	r.HandleFunc("/links", postToLink(linkSvc)).Methods(http.MethodPost)
	r.HandleFunc("/links/all", getFromLink(linkSvc)).Methods(http.MethodGet)
	r.HandleFunc("/links/{shortlink}", putToLink(linkSvc)).Methods(http.MethodPut)
	r.HandleFunc("/links/{shortlink}", delFromLink(linkSvc)).Methods(http.MethodDelete)
	// MiddleWare first goes JWT second goes Logging
	r.Use(JWTCheckMiddleware)
	r.Use(LoggingMiddleware)
	return r
}

// postTokenRefresh - get new pair of jwt tokens when access token is expired
func postTokenRefresh(svc linkSvc) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {
		type TokenAnswer struct {
			Access  string `json:"accessToken"`
			Refresh string `json:"refreshToken"`
		}

		props, _ := request.Context().Value(ctxKey{}).(jwt.MapClaims)

		UID := fmt.Sprintf("%v", props["uid"])

		/*
			Issuer := fmt.Sprintf("%v", props["iss"])

			if Issuer != "weblink_refresh" {
				ResponseAPIError(w, 7, http.StatusBadRequest)
				return
			}
		*/
		tokenAccess, _ := GenJWTWithClaims(UID, 0)
		tokenRefresh, _ := GenJWTWithClaims(UID, 1)

		var jsonTokens = TokenAnswer{
			Access:  tokenAccess,
			Refresh: tokenRefresh,
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
		//json header check
		contentType := request.Header.Get("Content-Type")
		if contentType != "application/json" {
			ResponseAPIError(w, 9, http.StatusBadRequest)
			return
		}

		type TokenAnswer struct {
			Access  string `json:"accessToken"`
			Refresh string `json:"refreshToken"`
		}

		type PostJSONRq struct {
			UID string `json:"uid"`
		}

		var jsonPostRq = PostJSONRq{}

		err := json.NewDecoder(request.Body).Decode(&jsonPostRq)
		if err != nil {
			ResponseAPIError(w, 400, http.StatusBadRequest)
			return
		}
		// get uid
		if jsonPostRq.UID == "" {
			ResponseAPIError(w, 8, http.StatusBadRequest)
			return
		}

		UID := jsonPostRq.UID

		checkif := svc.WhoAmI()

		if checkif == 1 {

			//check if user with UID exists
			user, err4 := svc.GetUser(UID)
			if err4 != nil {
				fmt.Printf("can't check user in db %s \n", UID)
				ResponseAPIError(w, 10, http.StatusBadRequest)
				return
			}

			if user.Role == "" {
				user = model.User{
					Name:    jsonPostRq.UID,
					Passwd:  "123",
					Email:   "L@u.ca",
					Balance: "100.00",
					Role:    "USER",
				}

			}

			var err1 error
			UID, err1 = svc.PutUser(user)

			if err1 != nil {
				ResponseAPIError(w, 10, http.StatusBadRequest)
				return
			}

			fmt.Printf("pg added user. %s \n", user.Name)
		}

		tokenAccess, _ := GenJWTWithClaims(UID, 0)
		tokenRefresh, _ := GenJWTWithClaims(UID, 1)

		var jsonTokens = TokenAnswer{
			Access:  tokenAccess,
			Refresh: tokenRefresh,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(jsonTokens)
		if err != nil {
			ResponseAPIError(w, 10, http.StatusBadRequest)
			return
		}

		return

	}
}

// delFromLink deletes link from api storage by shortlink
func delFromLink(linkSvc linkSvc) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {

		checkif := linkSvc.WhoAmI()

		var usefulUID, storageKey string
		var res bool
		usefulUID, storageKey, res = ValidateRequestShortLink(request, linkSvc)
		var flag bool = false

		if checkif == 1 {
			props, _ := request.Context().Value(ctxKey{}).(jwt.MapClaims)
			//fmt.Println(props["uid"])
			UID := fmt.Sprintf("%v", props["uid"])

			suid, _ := linkSvc.FindSuperUser()
			if suid == UID {
				// superuser deletes other user record here
				params := mux.Vars(request)
				storageKey = params["shortlink"]
				flag = true
			}
		}

		if !res && !flag {
			ResponseAPIError(w, 9, http.StatusBadRequest)
			return
		}

		//found key, delete it
		err := linkSvc.Del(usefulUID, storageKey, false)
		if err != nil {
			ResponseAPIError(w, 10, http.StatusBadRequest)
			return
		}

	}
}

// putToLink updates link from api storage
func putToLink(linkSvc linkSvc) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		//json header check
		contentType := request.Header.Get("Content-Type")
		if contentType != "application/json" {
			ResponseAPIError(w, 9, http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		checkif := linkSvc.WhoAmI()

		var usefulUID string
		var res bool
		usefulUID, _, res = ValidateRequestShortLink(request, linkSvc)
		var flag bool = false

		if checkif == 1 {
			props, _ := request.Context().Value(ctxKey{}).(jwt.MapClaims)
			//fmt.Println(props["uid"])
			UID := fmt.Sprintf("%v", props["uid"])
			params := mux.Vars(request)
			shortURL := params["shortlink"]

			suid, _ := linkSvc.FindSuperUser()
			if suid == UID {
				// superuser updates other user record here
				// get uid of that user
				dbElem, _ := linkSvc.Get(UID, shortURL, true)
				usefulUID = dbElem.UID
				flag = true
			}
		}

		if !res && !flag {
			ResponseAPIError(w, 9, http.StatusBadRequest)
			return
		}

		var element = model.DataEl{}
		//found key, work with body
		err := json.NewDecoder(request.Body).Decode(&element)
		if err != nil {
			ResponseAPIError(w, 9, http.StatusBadRequest)
			return
		}
		element.Datetime = time.Now()
		element.UID = usefulUID
		element.Active = 1
		//looks ok, update storage
		err = linkSvc.Put(usefulUID, element.Shorturl, element, false)
		if err != nil {
			ResponseAPIError(w, 9, http.StatusBadRequest)
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
		//json header check
		contentType := request.Header.Get("Content-Type")
		if contentType != "application/json" {
			ResponseAPIError(w, 9, http.StatusBadRequest)
			return
		}

		storageKeys, UID, err := GetUserStorageKeys(request, linkSvc)
		if err != nil {
			ResponseAPIError(w, 10, http.StatusBadRequest)
			return
		}

		checkif := linkSvc.WhoAmI()
		//db version supports payments for adding links
		if checkif == 1 {
			//make payment of 50.00 for the user account who uploaded link from SU account

			user, err2 := linkSvc.GetUser(UID)
			if err2 != nil {
				log.Printf("Coild not get user profile, err: %v\n", err2)
			}

			if user.Role == "CREATOR" || user.Role == "SUPERUSER" {
				//find payer - su
				suid, err1 := linkSvc.FindSuperUser()
				if err1 != nil {
					log.Printf("Could not find suid.. sorry, payment cannot be done.. err: %v\n", err1)
				}

				err1 = linkSvc.PayUser(suid, UID, "50.00")
				if err1 != nil {
					log.Printf("Payment error, payment to cannot be done.. err: %v\n", err1)
				}
			} else {
				log.Printf("user is not CREATOR, cannot add link and no payment available\n")
				ResponseAPIError(w, 401, http.StatusBadRequest)
				return
			}

		}

		var element = model.DataEl{}

		w.Header().Set("Content-Type", "application/json")

		err = json.NewDecoder(request.Body).Decode(&element)
		if err != nil {
			ResponseAPIError(w, 9, http.StatusBadRequest)
			return
		}
		// check if we have key
		if element.Shorturl == "" {
			ResponseAPIError(w, 11, http.StatusBadRequest)
			return
		}

		element.Datetime = time.Now()
		// check if this key already exists
		for _, storageKey := range storageKeys {
			if storageKey == element.Shorturl {
				ResponseAPIError(w, 5, http.StatusBadRequest)
				return
			}
		}
		element.UID = UID
		element.Active = 1
		err = linkSvc.Put(UID, element.Shorturl, element, false)
		if err != nil {
			ResponseAPIError(w, 10, http.StatusBadRequest)
			return
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

		// logic for file - get links of the user
		storageKeys, UID, err := GetUserStorageKeys(request, linkSvc)
		if err != nil {
			ResponseAPIError(w, 10, http.StatusBadRequest)
			return
		}

		for _, storageKey := range storageKeys {
			getElement, errfor := linkSvc.Get(UID, storageKey, false)
			if errfor != nil {
				ResponseAPIError(w, 10, http.StatusBadRequest)
				return
				//http.Error(w, "Cannot read from repo", http.StatusBadRequest)
			}

			datajson.Data = append(datajson.Data, getElement)
		}

		// sort by date asc
		sort.Slice(datajson.Data, func(i, j int) bool {
			return datajson.Data[i].Datetime.Before(datajson.Data[j].Datetime)
		})

		checkif := linkSvc.WhoAmI()

		if checkif == 1 {
			//logic for pg if user role = user list all links in database

			user, err1 := linkSvc.GetUser(UID)
			if err1 != nil {
				log.Printf("could not get user profie")
				// todo insert reply with error
				return
			}
			if user.Role == "USER" || user.Role == "SUPERUSER" {
				//
				sqlData, err3 := linkSvc.GetAll()
				if err3 != nil {
					//to do insert reply with error
					return
				}
				err2 := json.NewEncoder(w).Encode(sqlData)
				if err2 != nil {
					//to do insert reply with error
					return
				}
				//send http reply
				return
			}
		}

		err = json.NewEncoder(w).Encode(datajson)
		if err != nil {
			//to do insert reply with error
			return
		}

	}
}

// getShortStat - get one link from api
func getShortStat(linkSvc linkSvc) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		checkif := linkSvc.WhoAmI()

		if checkif == 1 {
			props, _ := request.Context().Value(ctxKey{}).(jwt.MapClaims)
			//fmt.Println(props["uid"])
			UID := fmt.Sprintf("%v", props["uid"])
			suid, _ := linkSvc.FindSuperUser()
			if suid == UID {
				//get link info of other user
				params := mux.Vars(request)
				shortURL := params["shortlink"]
				//ignore uid
				getElement, err := linkSvc.Get(UID, shortURL, true)
				if err != nil {
					ResponseAPIError(w, 10, http.StatusBadRequest)
					return
				}

				var datajson = model.Data{}
				datajson.Data = append(datajson.Data, getElement)
				err = json.NewEncoder(w).Encode(datajson)
				if err != nil {
					return
				}
				//finish and reply
				return

			}

		}
		// check user authorization, get user UID, get key (for this user, check if key exists)
		// if res - yes then do the action  - give string from repo as json
		UID, storageKey, res := ValidateRequestShortLink(request, linkSvc)
		if !res {
			ResponseAPIError(w, 11, http.StatusBadRequest)
			return
		}

		getElement, err := linkSvc.Get(UID, storageKey, false)
		if err != nil {
			ResponseAPIError(w, 10, http.StatusBadRequest)
			return
		}

		var datajson = model.Data{}
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
		shortURL := params["shortlink"]
		// GetUn retreives link and updates redir count
		URL, err := linkSvc.GetUn(shortURL)
		if err != nil {
			ResponseAPIError(w, 10, http.StatusBadRequest)
			return
		}
		if URL == "" {
			ResponseAPIError(w, 404, http.StatusBadRequest)
			return
		}

		checkif := linkSvc.WhoAmI()
		//db version supports payments for opening links
		if checkif == 1 {
			//make payment of 10.00 for the superuser account from USER who opened link
			//get UID from token

			amount := "10.0"
			props, _ := request.Context().Value(ctxKey{}).(jwt.MapClaims)
			//fmt.Println(props["uid"])
			UID := fmt.Sprintf("%v", props["uid"])

			user, err2 := linkSvc.GetUser(UID)
			if err2 != nil {
				log.Printf("Could not get user profile, err: %v\n", err2)
			}

			if user.Role == "USER" {
				//todo check if balance is not less then amount
				balance, _ := strconv.ParseFloat(user.Balance, 64)
				flamount, _ := strconv.ParseFloat(amount, 64)
				if balance < 0 || balance < flamount {
					ResponseAPIError(w, 402, http.StatusBadRequest)
					return
				}
				//find payer - su
				suid, err1 := linkSvc.FindSuperUser()
				if err1 != nil {
					log.Printf("Could not find suid.. sorry, payment cannot be done.. err: %v\n", err1)
				}

				err1 = linkSvc.PayUser(UID, suid, amount)
				if err1 != nil {
					log.Printf("Payment error, payment to cannot be done.. err: %v\n", err1)
				}
			} else {
				log.Printf("user is not CREATOR, no payment available\n")
			}
		}

		//log.Printf("opening user %s link  %s (short is %s) redirs(++) %d \n", getElement.UID, getElement.URL, getElement.Shorturl, getElement.Redirs)
		http.Redirect(w, request, URL, http.StatusSeeOther)
		//<a href="/shortopen/www.mail.ru">See Other</a>.
		return

	}
}
