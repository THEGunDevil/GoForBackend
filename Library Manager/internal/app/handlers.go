package app

import (
	"LibraryManager/internal/manager"
	"LibraryManager/internal/models"
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var (
	Ctx     = context.Background()
	Dbpool  *pgxpool.Pool
	Queries *manager.Queries
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL environment variable is not set")
	}
	Dbpool, err := pgxpool.New(Ctx, dbUrl)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	Queries = manager.New(Dbpool)
}

func FilesReload() {
	ClearJSONFiles()
	allBooks, err := Queries.ListBooks(Ctx)
	if err != nil {
		fmt.Println("Unable to fetch books:", err)
		return
	}
	allBorrows, err := Queries.ListBorrows(Ctx)
	if err != nil {
		fmt.Println("Unable to fetch borrows:", err)
		return
	}
	allUsers, err := Queries.ListUsers(Ctx)
	if err != nil {
		fmt.Println("Unable to fetch books:", err)
		return
	}
	books := LoadBooks()
	borrows := LoadBorrows()
	users := LoadUsers()

	for _, b := range allBooks {
		books = append(books, ManagerBookToModel(b))
	}
	for _, b := range allBorrows {
		borrows = append(borrows, ManagerBorrowsToModel(b))
	}
	for _, u := range allUsers {
		users = append(users, ManagerUsersToModel(u))
	}
	SaveBooks(books)
	SaveBorrows(borrows)
	SaveUsers(users)
	fmt.Println("Files are reloaded from database into JSON.")
}
func AddBook() {
	reader := bufio.NewReader(os.Stdin)

	for {

		fmt.Println("Enter book's title: ")
		title, _ := reader.ReadString('\n')
		title = strings.TrimSpace(title)
		if title == "exit" {
			return
		}
		if title == "" {
			return
		}
		fmt.Println("Enter book's author: ")
		author, _ := reader.ReadString('\n')
		author = strings.TrimSpace(author)
		if title == "exit" {
			return
		}
		if title == "" {
			return
		}
		fmt.Println("Enter publication year: ")
		year, _ := reader.ReadString('\n')
		year = strings.TrimSpace(year)

		yearInt, err := strconv.Atoi(year)
		if title == "" {
			return
		}
		if err != nil {
			fmt.Println("The year must be an integer. Please try again.")
			continue
		}

		b, err := Queries.CreateBook(Ctx, manager.CreateBookParams{
			Title:  title,
			Author: author,
			Year:   pgtype.Int4{Int32: int32(yearInt), Valid: true},
		})
		if err != nil {
			fmt.Printf("Failed to create book in DB: %v\n", err)
			continue
		}
		FilesReload()
		fmt.Println("Book added successfully!", b)
	}
}

func AddUser() {
	reader := bufio.NewReader(os.Stdin)
	for {

		fmt.Println("Enter users name: ")
		name, _ := reader.ReadString('\n')
		name = strings.TrimSpace(name)
		if name == "exit" {
			return
		}
		if name == "" {
			return
		}
		fmt.Println("Enter users email: ")
		email, _ := reader.ReadString('\n')
		email = strings.TrimSpace(email)
		if email == "exit" {
			return
		}
		if email == "" {
			return
		}
		u, err := Queries.CreateUser(Ctx, manager.CreateUserParams{
			Name:  name,
			Email: pgtype.Text{String: email, Valid: true},
		})
		if err != nil {
			fmt.Println("failed to create user: ", err)
		}
		FilesReload()
		fmt.Println("User added successfully!", u)
	}

}

func DeleteBook() {
	fmt.Println("Do you want to delete the book with or ID or Name?")
	fmt.Println("1: With ID,")
	fmt.Println("2: With Name")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "exit" {
		return
	}
	if input == "" {
		return
	}
	choice, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("invalid Input", err)
	}
	switch choice {
	case 1:
		fmt.Println("Enter Book ID: ")
		reader = bufio.NewReader(os.Stdin)
		input, _ = reader.ReadString('\n')
		input = strings.TrimSpace(input)
		id, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("invalid Input", err)
		}
		_, err = Queries.DeleteBookWithId(Ctx, int32(id))
		if err != nil {
			fmt.Println("failed to delete book.", err)
		}
		fmt.Println("Successfully Deleted a book!")
	case 2:
		fmt.Println("Enter Book title: ")
		reader = bufio.NewReader(os.Stdin)
		input, _ = reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if err != nil {
			fmt.Println("invalid Input", err)
		}
		_, err = Queries.DeleteBookWithTitle(Ctx, input)
		if err != nil {
			fmt.Println("failed to delete book.", err)
		}
		fmt.Println("Successfully Deleted a book!")

	}
	FilesReload()
}
func ShowAllBooks() {
	books, err := Queries.ListBooks(Ctx)
	if err != nil {
		fmt.Println("Unable to fetch books:", err)
	}
	fmt.Println(books)
}
func ShowAllUser() {
	users, err := Queries.ListUsers(Ctx)
	if err != nil {
		fmt.Println("Unable to fetch users:", err)
	}
	fmt.Println(users)
}
func ClearJSONFiles() {
	SaveBooks([]models.Book{})
	SaveBorrows([]models.Borrows{})
	SaveUsers([]models.Users{})

}

func Borrow() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter your user ID: ")
	userInput, _ := reader.ReadString('\n')
	userInput = strings.TrimSpace(userInput)
	if userInput == "exit" {
		return
	}
	if userInput == "" {
		return
	}
	userID, err := strconv.Atoi(userInput)
	if err != nil {
		fmt.Println("Invalid user ID", err)
	}

	fmt.Println("Enter book ID to borrow: ")
	bookInput, _ := reader.ReadString('\n')
	bookInput = strings.TrimSpace(bookInput)
	if bookInput == "exit" {
		return
	}
	if bookInput == "" {
		return
	}
	bookID, err := strconv.Atoi(bookInput)
	if err != nil {
		fmt.Println("Invalid book ID", err)
	}

	// Step 1: Mark the book as borrowed
	_, err = Queries.UpdateBook(Ctx, manager.UpdateBookParams{
		ID:         int32(bookID),
		IsBorrowed: pgtype.Bool{Bool: true, Valid: true},
	})
	if err != nil {
		fmt.Println("Error updating book:", err)
		return
	}

	// Step 2: Create borrow record
	brw, err := Queries.CreateBorrow(Ctx, manager.CreateBorrowParams{
		UserID:     pgtype.Int4{Int32: int32(userID), Valid: true},
		BookID:     pgtype.Int4{Int32: int32(bookID), Valid: true},
		BorrowedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
		ReturnedAt: pgtype.Timestamp{Valid: false},
	})
	if err != nil {
		fmt.Println("Couldn't create a borrow record:", err)
		return
	}
	FilesReload()
	fmt.Println("Book borrowed successfully!", brw)
}

func ReturnBook() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter book id: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "exit" {
		return
	}
	if input == "" {
		return
	}
	bookId, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Invalid input.", err)
		return
	}
	fmt.Println("Enter user id: ")
	Input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(Input)
	if Input == "exit" {
		return
	}
	if input == "" {
		return
	}
	userId, err := strconv.Atoi(input)

	if err != nil {
		fmt.Println("Invalid input.", err)
		return
	}

	b, err := Queries.UpdateBook(Ctx, manager.UpdateBookParams{
		ID:         int32(bookId),
		IsBorrowed: pgtype.Bool{Bool: false, Valid: true},
	})
	if err != nil {
		fmt.Println("Error updating book:", err)
		return
	}

	_, err = Queries.DeleteBorrow(Ctx, manager.DeleteBorrowParams{
		BookID: pgtype.Int4{Int32: int32(bookId), Valid: true},
		UserID: pgtype.Int4{Int32: int32(userId), Valid: true},
	})

	if err != nil {
		fmt.Println("There was an error updating borrows db.", err)
	}
	FilesReload()

	fmt.Println("Book returned successfully!", b)
}
