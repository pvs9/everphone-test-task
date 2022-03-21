package request

import "mime/multipart"

type GiftUpdateRequest struct {
	Name string   `form:"name" json:"name" binding:"required"`
	Tags []string `form:"categories" json:"categories" binding:"required"`
}

type GiftDatasetRequest struct {
	Dataset *multipart.FileHeader `form:"dataset" binding:"required"`
}

type GiftDatasetEntity struct {
	Name string   `json:"name"`
	Tags []string `json:"categories"`
}
