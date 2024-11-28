package service

type CreateShortURLFromJSONReq struct {
	URL string `json:"url"`
}

type CreateShortURLFromJSONRes struct {
	Result string `json:"result"`
}
