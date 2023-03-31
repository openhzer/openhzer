package middleware

type AdminLoad struct {
	Username string `json:"username"`
}

type WechatLoad struct {
	UUID string `json:"uuid"`
}
