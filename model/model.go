package model

type RequestBody struct {
	UserID int    `json:"userId" validate:"required,min=1"`
	ID     int    `json:"id" validate:"omitempty,min=0"`
	Title  string `json:"title" validate:"required,min=1"`
	Body   string `json:"body" validate:"required,min=1"`
}

type ResponseBody struct {
	UserID int    `json:"userId"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type ResponseArrayBody struct {
	Posts []ResponseBody `json:"posts"`
}
