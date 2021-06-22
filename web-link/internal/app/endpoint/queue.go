package endpoint

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/pehks1980/go_gb_be1_kurs/web-link/internal/pkg/model"
)

// интерфейс очередного сервиса также имеет put get - для работы с файлохранилищем
// файлохранилище это судя по всему словарь (обьект json в строковом виде)
// у "драйвера хранилища" методы
type queueSvc interface {
	Get(uid, key string) (model.DataEl, error)
	Put(uid, key string, value model.DataEl) error
	Del(uid, key string) error
	List(uid string) ([]string, error)
	GetUn(shortlink string) (model.DataEl, error)
}

// регистрация роутинга путей типа urls.py для обработки сервером
func RegisterPublicHTTP(queueSvc queueSvc) *mux.Router {
	// mux golrilla почему он? не знаю, - прикольное название, простота работы..
	r := mux.NewRouter()
	// JWT authorization
	r.HandleFunc("/user/auth", postAuth(queueSvc)).Methods(http.MethodPost)
	r.HandleFunc("/token/refresh", postTokenRefresh(queueSvc)).Methods(http.MethodPost)
	// main function
	r.HandleFunc("/shortopen/{shortlink}", getShortOpen(queueSvc)).Methods(http.MethodGet)
	r.HandleFunc("/shortstat/{shortlink}", getShortStat(queueSvc)).Methods(http.MethodGet)
	// links crud
	r.HandleFunc("/links", postToQueue(queueSvc)).Methods(http.MethodPost)
	r.HandleFunc("/links/all", getFromQueue(queueSvc)).Methods(http.MethodGet)
	r.HandleFunc("/links/{shortlink}", putToQueue(queueSvc)).Methods(http.MethodPut)
	r.HandleFunc("/links/{shortlink}", delFromQueue(queueSvc)).Methods(http.MethodDelete)

	r.Use(LoggingMiddleware)
	r.Use(JWTCheckMiddleware)
	return r
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		//вызов следующего хендлера в цепочке
		next.ServeHTTP(w, r)

		log.Printf("request: %s %s - %v\n",
			r.Method,
			r.URL.EscapedPath(),
			time.Since(start),
		)
	})
}

func GenJWTWithClaims(uidtext string, token_type int) (string, error) {
	mySigningKey := []byte("AllYourBase")

	type MyCustomClaims struct {
		Uid string `json:"uid"`
		jwt.StandardClaims
	}
	// type 0  access token is valid for 24 hours
	var time_expiry = time.Now().Add(time.Hour * 24).Unix()
	var issuer = "weblink_access"

	if token_type == 1 {
		// refresh token type 1 is valid for 5 days
		time_expiry = time.Now().Add(time.Hour * 24 * 5).Unix()
		issuer = "weblink_refresh"
	}

	// Create the Claims
	claims := MyCustomClaims{
		uidtext,
		jwt.StandardClaims{
			ExpiresAt: time_expiry, // access token will expire in 24h after creating
			Issuer:    issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}
	fmt.Printf("%v %v", ss, err)
	return ss, nil
	//Output: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJleHAiOjE1MDAwLCJpc3MiOiJ0ZXN0In0.HE7fK0xOQwFEr4WDgRWj4teRPZ6i3GLwD5YCm6Pwu_c <nil>
}

func JWTCheckMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.RequestURI == "/user/auth" {
			//bypass jwt check when authenticating
			next.ServeHTTP(w, r)
			return
		}

		re := regexp.MustCompile(`/shortopen/`)
		res := re.FindStringSubmatch(r.RequestURI)
		if len(res) != 0 {
			//bypass jwt check when authenticating
			next.ServeHTTP(w, r)
			return
		}

		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")

		if len(authHeader) != 2 {
			ResponseApiError(w, 2, http.StatusUnauthorized)
			return
		} else {
			jwtToken := authHeader[1]

			//jwtToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJleHAiOjE1MDAwLCJpc3MiOiJ0ZXN0In0.HE7fK0xOQwFEr4WDgRWj4teRPZ6i3GLwD5YCm6Pwu_c"
			//jwtToken, err := GenJWTWithClaims("ass")

			token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}
				SECRETKEY := "AllYourBase"
				return []byte(SECRETKEY), nil
			})

			if token.Valid {
				//fmt.Println("You look nice today")
				if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
					ctx := context.WithValue(r.Context(), "props", claims)
					// Access context values in handlers like this
					// props, _ := r.Context().Value("props").(jwt.MapClaims)
					if r.RequestURI != "/token/refresh" {
						// allow access to all API nodes with access token
						iss := fmt.Sprintf("%v", claims["iss"])
						if iss == "weblink_access" {
							next.ServeHTTP(w, r.WithContext(ctx))
							return
						}
					} else {
						//allow only refresh tokens to go to /token/refresh endpoint
						//check type of token iss should be weblink_refresh
						iss := fmt.Sprintf("%v", claims["iss"])
						if iss == "weblink_refresh" {
							next.ServeHTTP(w, r.WithContext(ctx))
							return
						}
						ResponseApiError(w, 7, http.StatusUnauthorized)
						return
					}

				} else {
					log.Printf("%v \n", err)
					ResponseApiError(w, 2, http.StatusUnauthorized)
					return
				}

			} else if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
					log.Printf("Token is either expired or not active yet %v", err)
					ResponseApiError(w, 1, http.StatusUnauthorized)
					return
				}
			}
		}
		ResponseApiError(w, 3, http.StatusUnauthorized)
	})
}

func postTokenRefresh(svc queueSvc) func(http.ResponseWriter, *http.Request) {
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
			//http.Error(w, "Please provide refresh token, or authenticate again", http.StatusBadRequest)
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

// autheticate and give authorization token
func postAuth(svc queueSvc) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {

		type TokenAnswer struct {
			Access  string `json:"accessToken"`
			Refresh string `json:"refreshToken"`
		}

		type PostJsonRq struct {
			Uid string `json:"uid"`
		}

		contentType := request.Header.Get("Content-Type")

		switch contentType {
		case "application/json":
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

		default:
			ResponseApiError(w, 8, http.StatusBadRequest)
			return

		}

	}
}

// del
func delFromQueue(queueSvc queueSvc) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		//w.Header().Set("Content-Type", "application/json")
		if UID, storageKey, res := validateRequestShortLink(w, request, queueSvc); res == true {
			//found key, delete it
			err := queueSvc.Del(UID, storageKey)
			if err != nil {
				ResponseApiError(w, 10, http.StatusBadGateway)
			}
			w.WriteHeader(http.StatusOK)
			return
		}
		return

	}
}

// put
func putToQueue(queueSvc queueSvc) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		var element = model.DataEl{}
		contentType := request.Header.Get("Content-Type")
		w.Header().Set("Content-Type", "application/json")

		switch contentType {
		case "application/json":

			if UID, _, res := validateRequestShortLink(w, request, queueSvc); res == true {
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
				err = queueSvc.Put(UID, element.Shorturl, element)
				if err != nil {
					ResponseApiError(w, 10, http.StatusBadGateway)
					//http.Error(w, "Cannot write to repo", http.StatusBadRequest)
				}
				// form answer json
				err = json.NewEncoder(w).Encode(element)
				if err != nil {
					return
				}
				return
			}

		default:
			ResponseApiError(w, 9, http.StatusBadRequest)
			return
		}

	}
}

// вьюха для post
func postToQueue(queueSvc queueSvc) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		var element = model.DataEl{}
		contentType := request.Header.Get("Content-Type")
		w.Header().Set("Content-Type", "application/json")

		switch contentType {
		case "application/json":
			storageKeys, UID, err := getUserStorageKeys(request, queueSvc)
			if err != nil {
				ResponseApiError(w, 10, http.StatusBadGateway)
				return
				//http.Error(w, "Cannot read List of keys from repo", http.StatusBadRequest)
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
					//http.Error(w, "This shortlink already exists", http.StatusBadRequest)
					return
				}
			}
			element.UID = UID
			element.Active = 1
			err = queueSvc.Put(UID, element.Shorturl, element)
			if err != nil {
				ResponseApiError(w, 10, http.StatusBadGateway)
				//http.Error(w, "Cannot write to repo", http.StatusBadRequest)
			}
			w.WriteHeader(http.StatusCreated) // this has to be the first write!!!
			err = json.NewEncoder(w).Encode(element)
			if err != nil {
				return
			}
			return

		default:
			ResponseApiError(w, 9, http.StatusBadRequest)
			return
		}

	}
}

// вьюха для get
func getFromQueue(queueSvc queueSvc) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {

		w.Header().Set("Content-Type", "application/json")
		var datajson = model.Data{}

		var err error

		storageKeys, UID, err := getUserStorageKeys(request, queueSvc)
		if err != nil {
			ResponseApiError(w, 10, http.StatusBadGateway)
			return
			///http.Error(w, "Cannot read List of keys from repo", http.StatusBadRequest)
		}

		for _, storageKey := range storageKeys {

			getElement, err := queueSvc.Get(UID, storageKey)
			if err != nil {
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
		//log.Println(getElement)
	}
}

func getShortStat(queueSvc queueSvc) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {

		// check user authorization, get user UID, get key (for this user, check if key exists)
		// if res - yes then do the action  - give string from repo as json
		if UID, storageKey, res := validateRequestShortLink(w, request, queueSvc); res == true {
			w.Header().Set("Content-Type", "application/json")
			var datajson = model.Data{}

			getElement, err := queueSvc.Get(UID, storageKey)
			if err != nil {
				ResponseApiError(w, 10, http.StatusBadGateway)
				//http.Error(w, "Cannot read from repo", http.StatusBadRequest)
				return
			}
			datajson.Data = append(datajson.Data, getElement)
			err = json.NewEncoder(w).Encode(datajson)
			if err != nil {
				return
			}

		}

	}
}

// getUserStorageKeys - get all keys for this user in repo
func getUserStorageKeys(request *http.Request, queueSvc queueSvc) ([]string, string, error) {
	props, _ := request.Context().Value("props").(jwt.MapClaims)
	//fmt.Println(props["uid"])
	UID := fmt.Sprintf("%v", props["uid"])

	storageKeys, err := queueSvc.List(UID)

	return storageKeys, UID, err
}

// validateRequestShortLink - валидация shortlink + токен авторизации
// ошибка валидации вылетает клиенту сама а результат - true ok
// Эту ф можно применить там где есть этот парам
func validateRequestShortLink(w http.ResponseWriter, request *http.Request, queueSvc queueSvc) (string, string, bool) {

	storageKeys, UID, err := getUserStorageKeys(request, queueSvc)

	if err != nil {
		// real error
		// http.Error(w, "Cannot read List of keys from repo", http.StatusBadRequest)
		// api error
		ResponseApiError(w, 10, http.StatusBadGateway)
		//http.Error(w, "Cannot find link", http.StatusBadRequest)
		return "", "", false
	}

	params := mux.Vars(request)
	shortUrl := params["shortlink"]

	for _, storageKey := range storageKeys {
		if storageKey == shortUrl {
			return UID, storageKey, true
		}
	}

	ResponseApiError(w, 6, http.StatusBadRequest)
	return "", "", false
}

func getShortOpen(queueSvc queueSvc) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {

		// check user authorization, get user UID, get key (for this user, check if key exists)
		// if res - yes then do the action  - redir
		//if UID, storageKey, res := validateRequestShortLink(w, request, queueSvc); res == true {
		// get data
		// update data
		// redir to real link

		//make link opened without authorization
		// we make uid = '', shortlink += : so it will find it if it exists (shortlink must be really unique)
		// otherwise wrong link will be opened and updated

		params := mux.Vars(request)
		shortUrl := params["shortlink"]

		getElement, err := queueSvc.GetUn(shortUrl)
		if err != nil {
			ResponseApiError(w, 10, http.StatusBadGateway)
			return
			//http.Error(w, "Cannot read from repo", http.StatusBadRequest)
		}
		getElement.Redirs++
		UID := getElement.UID

		err = queueSvc.Put(UID, getElement.Shorturl, getElement)
		if err != nil {
			ResponseApiError(w, 10, http.StatusBadGateway)
			return
			//http.Error(w, "Cannot read from repo", http.StatusBadRequest)
		}
		log.Printf("opening user %s link  %s (short is %s) redirs(++) %d \n", getElement.UID, getElement.URL, getElement.Shorturl, getElement.Redirs)
		http.Redirect(w, request, getElement.URL, http.StatusSeeOther)
		//return // todo redir here linke this? <a href="/shortopen/www.mail.ru">See Other</a>.
		return
		//}
		//ResponseApiError(w,400,http.StatusBadRequest )
	}

}
