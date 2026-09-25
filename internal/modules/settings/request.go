package settings

type UpdateSettingsRequest struct {
	Settings []UpdateSettingRequest `json:"settings" binding:"required,min=1,dive"`
}

type UpdateSettingRequest struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"`
}
