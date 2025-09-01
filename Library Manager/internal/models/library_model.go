package models

type Book struct {
	ID         int32  `json:"id"`
	Title      string `json:"title"`
	Author     string `json:"author"`
	Year       int32    `json:"year"`
	IsBorrowed bool   `json:"is_borrowed"`
}

type Users struct {
	ID    int32  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Borrows struct {
	ID         int32  `json:"id"`
	UserId     int32 `json:"user_id"`
	BookId     int32 `json:"book_id"`
	BorrowedAt string `json:"borrowed_at"`
	ReturnedAt string `json:"returned_at"`
}


