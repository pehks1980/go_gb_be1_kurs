package service

import (
	"log"

	"github.com/pehks1980/go_gb_be1_kurs/web-link/internal/pkg/model"
	web_broker "github.com/pehks1980/go_gb_be1_kurs/web-link/pkg/web-broker"
)

func (s *Service) Put(req *web_broker.PutValueReq) error {
	if err := s.repo.Put(&model.PutValue{
		Key:   req.Key,
		Value: req.Value,
	}); err != nil {
		log.Printf("service/Put: put repo err: %v", err)
		return err
	}

	return nil
}
