package endpoint

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/pehks1980/go_gb_be1_kurs/web-link/internal/pkg/model"
	"log"
	"net/http"
	"strings"
	"time"
)

// интерфейс очередного сервиса также имеет put get - для работы с файлохранилищем
// файлохранилище это судя по всему словарь (обьект json в строковом виде)
// у "драйвера хранилища" методы
type queueSvc interface {
	Get(uid, key string) (model.DataEl, error)
	Put(uid, key string, value model.DataEl) error
	Del(uid, key string) error
	List(uid string) ([]string, error)
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

		if r.RequestURI == "/user/auth"{
			//bypass jwt check when authenticating
			next.ServeHTTP(w, r)
			return
		}

		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")

		if len(authHeader) != 2 {
			fmt.Println("Malformed token")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Malformed Token"))
		} else {
			jwtToken := authHeader[1]

			//jwtToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJleHAiOjE1MDAwLCJpc3MiOiJ0ZXN0In0.HE7fK0xOQwFEr4WDgRWj4teRPZ6i3GLwD5YCm6Pwu_c"
			//jwtToken, err := GenJWTWithClaims("ass")

			token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}
				//SECRETKEY := "MYjwtKEY"

				SECRETKEY := "AllYourBase"
				return []byte(SECRETKEY), nil
			})

			if token.Valid {
				fmt.Println("You look nice today")
				if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
					ctx := context.WithValue(r.Context(), "props", claims)
					// Access context values in handlers like this
					// props, _ := r.Context().Value("props").(jwt.MapClaims)
					if r.RequestURI != "/token/refresh" {
						next.ServeHTTP(w, r.WithContext(ctx))
						return
					} else {
						//allow only refresh tokens to go to /token/refresh endpoint
						//check type of token iss should be weblink_refresh
						iss := fmt.Sprintf("%v", claims["iss"])
						if iss == "weblink_refresh"{
							next.ServeHTTP(w, r.WithContext(ctx))
							return
						}
					}

				} else {
					fmt.Println(err)
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("Unauthorized"))
				}

			} else if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&jwt.ValidationErrorMalformed != 0 {
					fmt.Println("That's not even a token")
				} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
					// Token is either expired or not active yet
					fmt.Println("Timing is everything")
				} else {
					fmt.Println("Couldn't handle this token:", err)
				}
			} else {
				fmt.Println("Couldn't handle this token:", err)
			}

			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
		}
	})
}

func postTokenRefresh(svc queueSvc) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {
		type TokenAnswer struct {
			Access string `json:"accessToken"`
			Refresh string `json:"refreshToken"`
		}

		props, _ := request.Context().Value("props").(jwt.MapClaims)
		//fmt.Println(props["uid"])
		UID := fmt.Sprintf("%v", props["uid"])
		Issuer := fmt.Sprintf("%v", props["iss"])

		if Issuer != "weblink_refresh" {
			http.Error(w, "Please provide refresh token, or authenticate again", http.StatusBadRequest)
			return
		}

		token_access, _:= GenJWTWithClaims(UID,0)
		token_refresh, _:= GenJWTWithClaims(UID,1)

		var jsonTokens = TokenAnswer {
			Access:  token_access,
			Refresh: token_refresh,
		}

		w.Header().Set("Content-Type", "application/json")
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
			Access string `json:"accessToken"`
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
					http.Error(w, "Unable to unmarshal JSON", http.StatusBadRequest)
					return
				}
				// get uid
				if jsonPostRq.Uid == "" {
					http.Error(w, "uid nil, please set uid!!!", http.StatusBadRequest)
					return
				}

				token_access, _:= GenJWTWithClaims(jsonPostRq.Uid,0)
				token_refresh, _:= GenJWTWithClaims(jsonPostRq.Uid,1)

				var jsonTokens = TokenAnswer {
					Access:  token_access,
					Refresh: token_refresh,
				}

				w.Header().Set("Content-Type", "application/json")
				err = json.NewEncoder(w).Encode(jsonTokens)
				if err != nil {
					return
				}
				return

			default:
				http.Error(w, "Unknown content type", http.StatusBadRequest)
				return

		}


	}
}

// del
func delFromQueue(queueSvc queueSvc) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		//get uid from JWT token
		props, _ := request.Context().Value("props").(jwt.MapClaims)
		//fmt.Println(props["uid"])
		UID := fmt.Sprintf("%v", props["uid"])

		storageKeys, err := queueSvc.List(UID)
		if err != nil {
			http.Error(w, "Unable to get List", http.StatusBadRequest)
			return
		}

		params := mux.Vars(request)
		shortUrl := params["shortlink"]
		for _, key := range storageKeys{
			if key == shortUrl {
				//found key, delete it
				err := queueSvc.Del(UID, shortUrl)
				if err != nil {
					http.Error(w, "Cannot write to repo", http.StatusBadRequest)
				}

				return
			}
		}
		// key not found
		http.Error(w, "key is not exist", http.StatusBadRequest)
		return

	}
}

//put
func putToQueue(queueSvc queueSvc) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		var element = model.DataEl{}
		contentType := request.Header.Get("Content-Type")

		//get uid from JWT token
		props, _ := request.Context().Value("props").(jwt.MapClaims)
		//fmt.Println(props["uid"])
		UID := fmt.Sprintf("%v", props["uid"])

		switch contentType {
		case "application/json":
			storageKeys, err := queueSvc.List(UID)
			if err != nil {
				http.Error(w, "Unable to get List", http.StatusBadRequest)
				return
			}

			params := mux.Vars(request)
			shortUrl := params["shortlink"]
			for _, key := range storageKeys{
				if key == shortUrl {
					//found key, work with body
					err = json.NewDecoder(request.Body).Decode(&element)
					if err != nil {
						http.Error(w, "Unable to unmarshal JSON", http.StatusBadRequest)
						return
					}
					element.Datetime = time.Now()
					element.UID = UID
					//looks ok, update storage
					err = queueSvc.Put(UID, element.Shorturl, element)
					if err != nil {
						http.Error(w, "Cannot write to repo", http.StatusBadRequest)
					}
					return
				}
			}
			http.Error(w, "key is not exist", http.StatusBadRequest)
			return

		default:
			http.Error(w, "Unknown content type", http.StatusBadRequest)
			return
		}

	}
}

// вьюха для post
func postToQueue(queueSvc queueSvc) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		var element = model.DataEl{}
		contentType := request.Header.Get("Content-Type")

		//get uid from JWT token
		props, _ := request.Context().Value("props").(jwt.MapClaims)
		//fmt.Println(props["uid"])
		UID := fmt.Sprintf("%v", props["uid"])

		storageKeys, err := queueSvc.List(UID)
		if err != nil {
			http.Error(w, "Cannot read List of keys from repo", http.StatusBadRequest)
		}

		switch contentType {
		case "application/json":

			err := json.NewDecoder(request.Body).Decode(&element)
			if err != nil {
				http.Error(w, "Unable to unmarshal JSON", http.StatusBadRequest)
				return
			}
			element.Datetime = time.Now()
			// check if this key already exists
			for _, storageKey := range storageKeys {
				if storageKey == element.Shorturl {
					http.Error(w, "This shortlink already exists", http.StatusBadRequest)
					return
				}
			}
			element.UID = UID
			err = queueSvc.Put(UID, element.Shorturl, element)
			if err != nil {
				http.Error(w, "Cannot write to repo", http.StatusBadRequest)
			}

		default:
			http.Error(w, "Unknown content type", http.StatusBadRequest)
			return
		}

	}
}

// вьюха для get
func getFromQueue(queueSvc queueSvc) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {

		w.Header().Set("Content-Type", "application/json")
		var datajson = model.Data{}
		//get uid from JWT token
		props, _ := request.Context().Value("props").(jwt.MapClaims)
		//fmt.Println(props["uid"])
		UID := fmt.Sprintf("%v", props["uid"])

		storageKeys, err := queueSvc.List(UID)
		if err != nil {
			http.Error(w, "Cannot read List of keys from repo", http.StatusBadRequest)
		}

		for _, storageKey := range storageKeys {
			getElement, err := queueSvc.Get(UID, storageKey)
			if err != nil {
				http.Error(w, "Cannot read from repo", http.StatusBadRequest)
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
		w.Header().Set("Content-Type", "application/json")
		var datajson = model.Data{}
		//get uid from JWT token
		props, _ := request.Context().Value("props").(jwt.MapClaims)
		//fmt.Println(props["uid"])
		UID := fmt.Sprintf("%v", props["uid"])

		storageKeys, err := queueSvc.List(UID)
		if err != nil {
			http.Error(w, "Cannot read List of keys from repo", http.StatusBadRequest)
		}
		params := mux.Vars(request)
		shortUrl := params["shortlink"]
		for _, storageKey := range storageKeys {
			if storageKey == shortUrl {
				getElement, err := queueSvc.Get(UID, storageKey)
				if err != nil {
					http.Error(w, "Cannot read from repo", http.StatusBadRequest)
					return
				}
				datajson.Data = append(datajson.Data, getElement)
				err = json.NewEncoder(w).Encode(datajson)
				if err != nil {
					return
				}
				return
			}
		}

		http.Error(w, "Cannot find link", http.StatusBadRequest)
		//log.Println(getElement)
	}
}

func getShortOpen(queueSvc queueSvc) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {

		w.Header().Set("Content-Type", "application/json")
		//get uid from JWT token
		props, _ := request.Context().Value("props").(jwt.MapClaims)
		//fmt.Println(props["uid"])
		UID := fmt.Sprintf("%v", props["uid"])

		storageKeys, err := queueSvc.List(UID)
		if err != nil {
			http.Error(w, "Cannot read List of keys from repo", http.StatusBadRequest)
		}
		params := mux.Vars(request)
		shortUrl := params["shortlink"]

		for _, storageKey := range storageKeys {
			if storageKey == shortUrl {
				// found key
				// get data
				// update data
				// redir to real link
				getElement, err := queueSvc.Get(UID, storageKey)
				if err != nil {
					http.Error(w, "Cannot read from repo", http.StatusBadRequest)
				}
				getElement.Redirs++
				err = queueSvc.Put(UID, storageKey, getElement)
				log.Printf("opening link %s (short is %s) redirs(++) %d \n", getElement.URL, getElement.Shorturl, getElement.Redirs)
				return // todo redir here
			}

		}
		http.Error(w, "Cannot find link", http.StatusBadRequest)

	}
}