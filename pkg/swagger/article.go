package swagger

type ArticleRequest struct {
	Title string `json:"title" validate:"required" example:"Title"`
	Desc  string `json:"desc" validate:"required" example:"Description"`
}
