package books

type Book struct {
	Id     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Read   bool   `json:"read"`
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

func updateBook(dto UpdateBookDto) Book {
	return Book{
		Id:     dto.Id,
		Title:  dto.Title,
		Author: dto.Author,
		Read:   dto.Read,
	}
}

type UpdateBookDto struct {
	Id     string `json:"id"`
	Title  string `json:"title,omitempty"`
	Author string `json:"author,omitempty"`
	Read   bool   `json:"read,omitempty"`
}

type CreateBookResponseDto struct {
	Id string `json:"id"`
}
