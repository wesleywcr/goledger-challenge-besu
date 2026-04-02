package api

import (
	"app/internal/service"
)

type Handler struct {
	service *service.Service
}

type setValueRequest struct {
	Value uint64 `json:"value"`
}

func NewHandler(serviceLayer *service.Service) *Handler {
	return &Handler{service: serviceLayer}
}
