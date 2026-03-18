package orgmodels

type CreateOrgRequestBody struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Plan        string `json:"plan"`
}

type UpdateOrgRequestBody struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	IconUrl     string `json:"icon_url"`
}
