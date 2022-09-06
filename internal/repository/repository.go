package repository

type Books interface {
	GetAllBooks() ([]string, error)
}

type Repository struct {
	Books
}
