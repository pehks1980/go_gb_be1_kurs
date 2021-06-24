package endpoint

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

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

					// now check for json type in header
					// all specific queries like shortstat etc
					// must have it in request!!!
					contentType := r.Header.Get("Content-Type")
					if contentType != "application/json" {
						ResponseApiError(w, 8, http.StatusBadRequest)
						return
					}

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