package services

import (
	"shrine/config"
	configTypes "shrine/types/config"
)

func GetConfig() configTypes.ConfigResponse {
	return configTypes.ConfigResponse{
		MaxFileSize:    config.Storage.MaxFileSize,
		MaxAttachments: config.Storage.MaxAttachments,
	}
}