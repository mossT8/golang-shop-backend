package model

import "tannar.moss/backend/internal/model"

type Response struct {
	User           model.UserResponse   `json:"user,omitempty"`
	Users          []model.UserResponse `json:"users,omitempty"`
	PagingResponse model.PagingResponse `json:"paging_response,omitempty"`
}
