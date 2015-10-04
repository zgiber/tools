package main

import "log"

type Service struct{}

func (service *Service) Answer(args interface{}, resp *int) error {
	log.Println(args)
	*resp = 42
	return nil
}
