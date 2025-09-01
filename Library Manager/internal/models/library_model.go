package models

type Book struct {
	ID         int32  `json:"id"`
	Title      string `json:"title"`
	Author     string `json:"author"`
	Year       int    `json:"year"`
	IsBorrowed bool   `json:"is_borrowed"`
}

type Users struct {
	ID    int32  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Borrows struct {
	ID         int32  `json:"id"`
	UserId     int32 `json:"title"`
	BookId     int32 `json:"author"`
	BorrowedAt string `json:"borrowed_at"`
	ReturnedAt string `json:"returned_at"`
}


