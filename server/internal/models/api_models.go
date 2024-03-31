package models

// APIResponse defined by NYT API response, list of books returned
type APIResponse struct {
	Status  string `json:"status"`
	Results struct {
		Books []Book `json:"books"`
	} `json:"results"`
}
