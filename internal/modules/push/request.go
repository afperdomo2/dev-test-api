package push

type RegisterTokenRequest struct {
	Token    string `json:"token" binding:"required,max=255"`
	Platform string `json:"platform" binding:"required,oneof=ios android"`
}

type UnregisterTokenRequest struct {
	Token string `json:"token" binding:"required,max=255"`
}
