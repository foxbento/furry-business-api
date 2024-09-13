package models

type Business struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Link        string `json:"link"`
	Type        string `json:"type"`
	Country     string `json:"country"`
	State       string `json:"state"`
	NSFW        bool   `json:"nsfw"`
	Overview    string `json:"overview"`
	Gendered    string `json:"gendered"`
	Conventions bool   `json:"conventions"`
	Notes       string `json:"notes"`
}