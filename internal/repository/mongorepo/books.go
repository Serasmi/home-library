package mongorepo

type Books struct {
}

func NewBooks() *Books {
	return &Books{}
}

func (b Books) GetAllBooks() ([]string, error) {
	return []string{"Book 1", "Book 2", "Book 3"}, nil
}
