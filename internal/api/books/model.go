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

type CreateBookResponseDto struct {
	Id string `json:"id"`
}
