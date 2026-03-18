package projectmodels

type CreateProjectRequestBody struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	Environment    string `json:"environment"`
	GenerateApiKey bool   `json:"generate_api_key"`
}

type UpdateProjectRequestBody struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Environment string `json:"environment"`
}
