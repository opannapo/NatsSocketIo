package dto

type MasterSendEmailMessage struct {
	Email        string      `json:"email" validate:"required,email"`
	TemplateType string      `json:"templateType" validate:"required"`
	TemplateData interface{} `json:"templateData" validate:"required"`
}

type TemplateInterface interface {
	TemplateChangePassword | TemplateForgotPassword
}

type TemplateChangePassword struct {
	AppUrl            string `json:"appUrl" validate:"required"`
	AccountFirstName  string `json:"accountFirstName" validate:"required"`
	AccountMiddleName string `json:"accountMiddleName"`
	AccountLastName   string `json:"accountLastName"`
}

type TemplateForgotPassword struct {
	AppUrl               string `json:"appUrl" validate:"required"`
	AccountFirstName     string `json:"accountFirstName" validate:"required"`
	AccountMiddleName    string `json:"accountMiddleName"`
	AccountLastName      string `json:"accountLastName"`
	PasswordVerification string `json:"passwordVerification" validate:"required"`
}
