package books

type Book struct {
	Id     string `json:"id" bson:"_id,omitempty"`
	Title  string `json:"title" bson:"title,omitempty"`
	Author string `json:"author" bson:"author,omitempty"`
	Read   bool   `json:"read" bson:"read,omitempty"`
}

func newBook(dto CreateBookDto) Book {
	return Book{
		Title:  dto.Title,
		Author: dto.Author,
		Read:   false,
	}
}

type CreateBookDto struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}

type UpdateBookDto struct {
	Id     string `json:"id" bson:"_id,omitempty"`
	Title  string `json:"title,omitempty" bson:"title,omitempty"`
	Author string `json:"author,omitempty" bson:"author,omitempty"`
	Read   bool   `json:"read,omitempty" bson:"read,omitempty"`
}

type CreateBookResponseDto struct {
	Id string `json:"id"`
}
