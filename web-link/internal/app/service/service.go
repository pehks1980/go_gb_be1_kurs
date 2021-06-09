package service

import (
	"github.com/pehks1980/go_gb_be1_kurs/web-link/internal/pkg/model"
	"log"
)

// repo имеет тип интерфейс (2 метода)
type repo interface {
	Get(key string) (model.DataEl, error)
	Put(key string, value model.DataEl) error
	Del(key string) error
	List() ([]string, error)
}

// service имеет тип структура
// содержит член repo
type Service struct {
	repo repo
}

// конструктор Service
// возвращает указатель на структуру с интерфейсом
// что в repo подставим та структура и будет - главное методы должны иметь одинаковую сигнатуру и имя!
func New(repo repo) *Service {
	return &Service{repo: repo}
}

func (s *Service) Put(key string, value model.DataEl) error {
	if err := s.repo.Put(key, value); err != nil {
		log.Printf("service/Put: put repo err: %v", err)
		return err
	}

	return nil
}

func (s *Service) Get(key string) (model.DataEl, error) {
	value, err := s.repo.Get(key)
	if err != nil {
		log.Printf("service/Get: get from repo err: %v", err)
		return model.DataEl{}, err
	}

	return value, nil
}

func (s *Service) Del(key string) error {
	if err := s.repo.Del(key); err != nil {
		log.Printf("service/Del: del repo err: %v", err)
		return err
	}

	return nil
}

func (s *Service) List() ([]string, error) {
	items, err := s.repo.List()
	if err != nil {
		log.Printf("service/List: get from repo err: %v", err)
		return nil, err
	}

	return items, nil
}
