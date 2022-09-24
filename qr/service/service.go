package service

var Service = NewService()

type IService interface{}

func NewService() IService {
	return &service{}
}

type service struct{}
