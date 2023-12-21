package config

type SmsConfig struct {
	AppId      string `json:"app_id"`
	SecretId   string `json:"secret_id"`
	SecretKey  string `json:"secret_key"`
	Sign       string `json:"sign"`
	TemplateId string `json:"template_id"`
}
