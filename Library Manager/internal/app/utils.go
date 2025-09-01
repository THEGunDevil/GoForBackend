package app

import (
	"LibraryManager/internal/manager"
	"LibraryManager/internal/models"
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const BooksFile = "books.json"
const UsersFile = "users.json"
const BorrowsFile = "borrows.json"

func SelectService() {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	service, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Invalid input")
	}
	switch service {
	case 1:
		AddBook()
	case 2:
		DeleteBook()
	case 3:
		Borrow()
	case 4:
		ShowAllBooks()
	// case 5:
	// 	ReturnBook()
	// case 6:
	// 	SearchBook()
	case 7:
		AddUser()
	// case 8:
	// 	ShowAllUser()
	// case 9:
	// 	Filter()
	case 0:
		return
	default:
		fmt.Println("Invalid service")
	}
}

func LoadBooks() []models.Book {
	file, err := os.ReadFile(BooksFile)
	if err != nil {
		return []models.Book{}
	}
	var books []models.Book
	err = json.Unmarshal(file, &books)
	if err != nil {
		fmt.Println("Error reading books file:", err)
		return []models.Book{}
	}
	return books
}
func LoadUsers() []models.Users {
	file, err := os.ReadFile(UsersFile)
	if err != nil {
		return []models.Users{}
	}
	var users []models.Users
	err = json.Unmarshal(file, &users)
	if err != nil {
		fmt.Println("Error reading users file:", err)
		return []models.Users{}
	}
	return users
}
func LoadBorrows() []models.Borrows {
	file, err := os.ReadFile(BorrowsFile)
	if err != nil {
		return []models.Borrows{}
	}
	var borrows []models.Borrows
	err = json.Unmarshal(file, &borrows)
	if err != nil {
		fmt.Println("Error reading borrows file:", err)
		return []models.Borrows{}
	}
	return borrows
}

func SaveBooks(books []models.Book) {
	data, err := json.Marshal(books)
	if err != nil {
		fmt.Println("Error marshalling boooks:", err)
		return
	}
	err = os.WriteFile(BooksFile, data, 0644)
	if err != nil {
		fmt.Println("Error saving books:", err)
	}
}
func SaveUsers(users []models.Users) {
	data, err := json.Marshal(users)
	if err != nil {
		fmt.Println("Error marshalling users:", err)
		return
	}
	err = os.WriteFile(UsersFile, data, 0644)
	if err != nil {
		fmt.Println("Error saving users:", err)
	}
}
func SaveBorrows(borrows []models.Borrows) {
	data, err := json.Marshal(borrows)
	if err != nil {
		fmt.Println("Error marshalling borrows:", err)
		return
	}
	err = os.WriteFile(BorrowsFile, data, 0644)
	if err != nil {
		fmt.Println("Error saving borrows:", err)
	}
}
func ManagerBookToModel(b manager.Book) models.Book {
	return models.Book{
		ID:         b.ID,
		Title:      b.Title,
		Author:     b.Author,
		Year:       int(b.Year.Int32),
		IsBorrowed: b.IsBorrowed.Bool, // pgtype.Bool → bool
	}
}
