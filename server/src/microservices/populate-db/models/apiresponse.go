package models

type APIResponse struct {
	Status  string `json:"status"`
	Results struct {
		Books []Book `json:"books"`
	} `json:"results"`
}
