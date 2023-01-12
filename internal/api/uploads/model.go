package uploads

import "time"

type Status uint

const (
	Created Status = iota
	InProgress
	Done
)

type Upload struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Filename  string    `json:"filename" bson:"filename"`
	BookID    string    `json:"bookId" bson:"bookId,omitempty"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	Status    Status    `json:"status" bson:"status"`
}

func newUpload(dto CreateUploadDTO) Upload {
	return Upload{
		Filename:  dto.Filename,
		BookID:    dto.BookID,
		CreatedAt: time.Now(),
		Status:    Created,
	}
}

type CreateUploadDTO struct {
	Filename string `json:"filename"`
	BookID   string `json:"bookId"`
}

type CreateUploadResponseDTO struct {
	ID string `json:"id"`
}

type ResponseDTO struct {
	Filename string `json:"filename"`
}
