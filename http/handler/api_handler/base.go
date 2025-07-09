package api_handler

import "github.com/jiangrx816/wechat/service/api_service"

func NewApiHandler() *ApiHandler {
	return &ApiHandler{
		service: api_service.NewApiService(),
	}
}

type ApiHandler struct {
	service *api_service.ApiService
}
