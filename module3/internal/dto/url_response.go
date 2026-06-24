package dto

type UrlResponse struct {
	Id    int64  `json:"id,omitempty"`
	Url   string `json:"url,omitempty"`
	Alias string `json:"alias,omitempty"`
}
