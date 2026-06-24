package dto

type UrlUpdateRequest struct {
	Id    int64  `json:"id"`
	Url   string `json:"url"`
	Alias string `json:"alias"`
}
