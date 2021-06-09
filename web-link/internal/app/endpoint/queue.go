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

		storageKeys, err := queueSvc.List()
		if err != nil {
			http.Error(w, "Unable to get List", http.StatusBadRequest)
			return
		}

		params := mux.Vars(request)
		shortUrl := params["shortlink"]
		for _, key := range storageKeys{
			if key == shortUrl {
				//found key, delete it
				err := queueSvc.Del(shortUrl)
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

		switch contentType {
		case "application/json":
			storageKeys, err := queueSvc.List()
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
					//looks ok, update storage
					err = queueSvc.Put(element.Shorturl, element)
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
