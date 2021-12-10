package common


type GetListParams struct {
	Page int `json:"page"`
	PageSize int `json:"page_size"`
}

type ReturnGetList struct {
	Id int `json:"id"`
	Cover string `json:"cover"`
	Title string `json:"title"`
}

type Select struct {
	Id int `json:"id"`
}