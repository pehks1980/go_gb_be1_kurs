package endpoint

import (
	"github.com/gorilla/mux"
	"github.com/pehks1980/go_gb_be1_kurs/web-link/internal/pkg/model"
	"net/http"
)
// интерфейс очередного сервиса также имеет put get - для работы с файлохранилищем
// файлохранилище это судя по всему словарь (обьект json в строковом виде)
// у "драйвера хранилища" методы
type queueSvc interface {
	Get(key string) (model.DataEl, error)
	Put(key string, value model.DataEl) error
	Del(key string) error
}

// регистрация роутинга путей типа urls.py для обработки сервером 
func RegisterPublicHTTP(queueSvc queueSvc) *mux.Router {
	// 
	r := mux.NewRouter()
	// HandleFunc registers a new route with a matcher for the URL path
	// путь - обработчик(интерфейс очередной)
	// func (r *Router) HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *Route

	// вызывается r.HandleFunc и r.Methods(http.MethodPut)
	// в документации вызов r.Methods(http.MethodPut) не обязателен
	r.HandleFunc("/{queue}", putToQueue(queueSvc) ).Methods(http.MethodPut)
	r.HandleFunc("/{queue}", getFromQueue(queueSvc) ).Methods(http.MethodGet)

	return r
}
// вьюха для put
func putToQueue(queueSvc queueSvc) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		// TODO: parse req and call queueSvc.Put(...)
		// когда пришел запрос http request
		// он тут обрабатывается - т.е это типа вьюха контроллер django
		// метод PUT значит это апдейт данные присланные request сохраняются в файлохранилище queueSvc.Put
		// респонз - ок или нот ok
	}
}
// вьюха для get
func getFromQueue(queueSvc queueSvc) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		// TODO: parse req and call queueSvc.Get(...)
		//вьюха для GET запроса - какие то данные в POST уходят
		// респонз - ок или нот ok
	}
}
