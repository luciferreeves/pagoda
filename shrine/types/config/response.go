package config

type ConfigResponse struct {
	MaxFileSize    int64 `json:"max_file_size"`
	MaxAttachments int   `json:"max_attachments"`
}