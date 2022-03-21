package request

import "mime/multipart"

type EmployeeUpdateRequest struct {
	Name string   `form:"name" json:"name" binding:"required"`
	Tags []string `form:"interests" json:"interests" binding:"required"`
}

type EmployeeDatasetRequest struct {
	Dataset *multipart.FileHeader `form:"dataset" binding:"required"`
}

type EmployeeDatasetEntity struct {
	Name string   `json:"name"`
	Tags []string `json:"interests"`
}
