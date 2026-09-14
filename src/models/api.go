package models

type ValidationErrorResponse struct {
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors"`
}

type PersonRequest struct {
	Name    string `json:"name"`
	Age     int32  `json:"age"`
	Address string `json:"address"`
	Work    string `json:"work"`
}

type PersonResponse struct {
	ID      int32  `json:"id"`
	Name    string `json:"name"`
	Age     int32  `json:"age"`
	Address string `json:"address"`
	Work    string `json:"work"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
