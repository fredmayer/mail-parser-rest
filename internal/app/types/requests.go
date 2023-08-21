package types

type PageRequest struct {
	Page int `query:"page"`
}

type AttachmentRequest struct {
	Mime string `json:"mime"`
	Name string `json:"name"`
}
