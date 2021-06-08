package service

import (
	"log"

	web_link "github.com/pehks1980/go_gb_be1_kurs/web-link"
)

func (s *Service) Get(req *web_link.GetValueReq) (*web_link.GetValueResp, error) {
	value, err := s.repo.Get(req.Key)
	if err != nil {
		log.Printf("service/Get: get from repo err: %v", err)
		return nil, err
	}

	return &web_broker.GetValueResp{Value: value}, nil
}
