package endpoint

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/pehks1980/go_gb_be1_kurs/web-link/internal/pkg/model"
	"net/http"
	"time"
)

// интерфейс очередного сервиса также имеет put get - для работы с файлохранилищем
// файлохранилище это судя по всему словарь (обьект json в строковом виде)
// у "драйвера хранилища" методы
type queueSvc interface {
	Get(key string) (model.DataEl, error)
	Put(key string, value model.DataEl) error
	Del(key string) error
	List() ([]string, error)
}

// регистрация роутинга путей типа urls.py для обработки сервером
func RegisterPublicHTTP(queueSvc queueSvc) *mux.Router {
	// mux golrilla почему он? не знаю, - прикольное название, простота работы..
	r := mux.NewRouter()
	//links crud
	r.HandleFunc("/links", postToQueue(queueSvc)).Methods(http.MethodPost)
	r.HandleFunc("/links/all", getFromQueue(queueSvc)).Methods(http.MethodGet)
	r.HandleFunc("/links/{shortlink}", putToQueue(queueSvc)).Methods(http.MethodPut)
	r.HandleFunc("/links/{shortlink}", delFromQueue(queueSvc)).Methods(http.MethodDelete)

	return r
}

//del
func delFromQueue(queueSvc queueSvc) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		// TODO: parse req and call queueSvc.Put(...)
		// когда пришел запрос http request
		// он тут обрабатывается - т.е это типа вьюха контроллер django
		// метод PUT значит это апдейт данные присланные request сохраняются в файлохранилище queueSvc.Put
		// респонз - ок или нот ok

		params := mux.Vars(request)
		shortUrl := params["shortlink"]
		// todo validate shortUrl with regex
		err := queueSvc.Del(shortUrl)
		if err != nil {
			http.Error(w, "Cannot write to repo", http.StatusBadRequest)
		}

	}
}

//put
func putToQueue(queueSvc queueSvc) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		// TODO: parse req and call queueSvc.Put(...)
		// когда пришел запрос http request
		// он тут обрабатывается - т.е это типа вьюха контроллер django
		// метод PUT значит это апдейт данные присланные request сохраняются в файлохранилище queueSvc.Put
		// респонз - ок или нот ok
		var element = model.DataEl{}
		contentType := request.Header.Get("Content-Type")

		switch contentType {
		case "application/json":
			err := json.NewDecoder(request.Body).Decode(&element)
			if err != nil {
				http.Error(w, "Unable to unmarshal JSON", http.StatusBadRequest)
				return
			}
			element.Datetime = time.Now()
			err = queueSvc.Put(element.Shorturl, element)
			if err != nil {
				http.Error(w, "Cannot write to repo", http.StatusBadRequest)
			}

		default:
			http.Error(w, "Unknown content type", http.StatusBadRequest)
			return
		}

		params := mux.Vars(request)
		shortUrl := params["shortlink"]

		err := queueSvc.Put(shortUrl, element)
		if err != nil {
			http.Error(w, "Cannot write to repo", http.StatusBadRequest)
		}

	}
}

// вьюха для post
func postToQueue(queueSvc queueSvc) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		// TODO: parse req and call queueSvc.Put(...)
		// когда пришел запрос http request
		// он тут обрабатывается - т.е это типа вьюха контроллер django
		// метод PUT значит это апдейт данные присланные request сохраняются в файлохранилище queueSvc.Put
		// респонз - ок или нот ok
		var element = model.DataEl{}
		contentType := request.Header.Get("Content-Type")

		switch contentType {
		case "application/json":
			err := json.NewDecoder(request.Body).Decode(&element)
			if err != nil {
				http.Error(w, "Unable to unmarshal JSON", http.StatusBadRequest)
				return
			}
			element.Datetime = time.Now()
			err = queueSvc.Put(element.Shorturl, element)
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
		// TODO: parse req and call queueSvc.Get(...)
		//вьюха для GET запроса - какие то данные в POST уходят
		// респонз - ок или нот ok

		w.Header().Set("Content-Type", "application/json")
		var datajson = model.Data{}

		storageKeys, err := queueSvc.List()
		if err != nil {
			http.Error(w, "Cannot read List of keys from repo", http.StatusBadRequest)
		}

		for _, storageKey := range storageKeys {
			getElement, err := queueSvc.Get(storageKey)
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
