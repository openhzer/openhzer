package captcha

type reqCaptcha struct {
	Data string `json:"data"`
	Id   string `json:"id"`
}

type Model struct {
	Id    string `json:"id" binding:"required"`
	Value string `json:"value" binding:"required"`
}
