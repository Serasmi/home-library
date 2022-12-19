package upload

import "time"

type Status uint

const (
	Created Status = iota
	InProgress
	Done
)

type Meta struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Filename  string    `json:"filename" bson:"filename"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	Status    Status    `json:"status" bson:"status"`
}

func newMeta(dto CreateMetaDTO) Meta {
	return Meta{
		Filename:  dto.Filename,
		CreatedAt: time.Now(),
		Status:    Created,
	}
}

type CreateMetaDTO struct {
	Filename string `json:"filename"`
}

type CreateMetaResponseDTO struct {
	ID string `json:"id"`
}
