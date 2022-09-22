package dto

type MasterCreateApplicationMessage struct {
	ApplicationName string `json:"applicationName" validate:"required"`
}
