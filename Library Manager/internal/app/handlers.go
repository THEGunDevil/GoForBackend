package app

import (
	"LibraryManager/internal/manager"
	"LibraryManager/internal/models"
	"bufio"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
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

func AddBook() {
	reader := bufio.NewReader(os.Stdin)

	for {

		fmt.Println("Enter book's title: ")
		title, _ := reader.ReadString('\n')
		title = strings.TrimSpace(title)

		fmt.Println("Enter book's author: ")
		author, _ := reader.ReadString('\n')
		author = strings.TrimSpace(author)

		fmt.Println("Enter publication year: ")
		year, _ := reader.ReadString('\n')
		year = strings.TrimSpace(year)

		yearInt, err := strconv.Atoi(year)
		if err != nil {
			fmt.Println("The year must be an integer. Please try again.")
			continue
		}

		// Insert into Postgres
		// books := LoadBooks()
		b, err := Queries.CreateBook(Ctx, manager.CreateBookParams{
			Title:  title,
			Author: author,
			Year:   pgtype.Int4{Int32: int32(yearInt), Valid: true},
		})
		if err != nil {
			fmt.Printf("Failed to create book in DB: %v\n", err)
			continue
		}

		// Save locally in JSON
		book := models.Book{
			ID:     b.ID, // Use DB-generated ID
			Title:  b.Title,
			Author: b.Author,
			Year:   int(b.Year.Int32)}

		books := LoadBooks()
		books = append(books, book)
		SaveBooks(books)

		fmt.Printf("Book added successfully! ID: %d, Title: %s\n", b.ID, b.Title)
	}
}

func AddUser() {
	reader := bufio.NewReader(os.Stdin)
	for {

		fmt.Println("Enter users name: ")
		name, _ := reader.ReadString('\n')
		name = strings.TrimSpace(name)
		fmt.Println("Enter users email: ")
		email, _ := reader.ReadString('\n')
		email = strings.TrimSpace(email)
		u, err := Queries.CreateUser(Ctx, manager.CreateUserParams{
			Name:  name,
			Email: pgtype.Text{String: email, Valid: true},
		})
		if err != nil {
			log.Fatalf("failed to create user: %v", err)
		}
		users := LoadUsers()
		user := models.Users{
			ID:    u.ID, // use the DB-generated ID
			Name:  u.Name,
			Email: u.Email.String,
		}
		users = append(users, user)
		SaveUsers(users)
		fmt.Println(u)
		fmt.Println("User added successfully!")
	}

}

func DeleteBook() {
	fmt.Println("Do you want to delete the book with or ID or Name?")
	fmt.Println("1: With ID,")
	fmt.Println("2: With Name")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	choice, err := strconv.Atoi(input)
	if err != nil {
		log.Fatalln("invalid Input", err)
	}
	switch choice {
	case 1:
		fmt.Println("Enter Book ID: ")
		reader = bufio.NewReader(os.Stdin)
		input, _ = reader.ReadString('\n')
		input = strings.TrimSpace(input)
		id, err := strconv.Atoi(input)
		if err != nil {
			log.Fatalln("invalid Input", err)
		}
		_, err = Queries.DeleteBookWithId(Ctx, int32(id))
		if err != nil {
			log.Fatalln("failed to delete book.", err)
		}
		fmt.Println("Successfully Deleted a book!")
	case 2:
		fmt.Println("Enter Book title: ")
		reader = bufio.NewReader(os.Stdin)
		input, _ = reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if err != nil {
			log.Fatalln("invalid Input", err)
		}
		_, err = Queries.DeleteBookWithTitle(Ctx, input)
		if err != nil {
			log.Fatalln("failed to delete book.", err)
		}
		fmt.Println("Successfully Deleted a book!")

	}

}
func ShowAllBooks() {
	books, err := Queries.ListBooks(Ctx)
	if err != nil {
		fmt.Println("Unable to fetch books:", err)
	}
	fmt.Println(books)
}
func ClearBooks() {
	SaveBooks([]models.Book{})
}

func FilesReload() {
	ClearBooks()
	allBooks, err := Queries.ListBooks(Ctx)
	if err != nil {
		fmt.Println("Unable to fetch books:", err)
		return
	}
	books := LoadBooks()

	for _, b := range allBooks {
		books = append(books, ManagerBookToModel(b))
	}
	SaveBooks(books)
	fmt.Println("Books reloaded from database into JSON.")
}

func Borrow() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter your user ID: ")
	userInput, _ := reader.ReadString('\n')
	userInput = strings.TrimSpace(userInput)
	userID, err := strconv.Atoi(userInput)
	if err != nil {
		log.Fatalln("Invalid user ID", err)
	}

	fmt.Println("Enter book ID to borrow: ")
	bookInput, _ := reader.ReadString('\n')
	bookInput = strings.TrimSpace(bookInput)
	bookID, err := strconv.Atoi(bookInput)
	if err != nil {
		log.Fatalln("Invalid book ID", err)
	}

	// Mark the book as borrowed in DB and get updated book info
	ub, err := Queries.UpdateBook(Ctx, manager.UpdateBookParams{
		ID:         int32(bookID),
		IsBorrowed: pgtype.Bool{Bool: true, Valid: true},
	})
	if err != nil {
		fmt.Println("Error Updating book:", err)
		return
	}

	brw, err := Queries.CreateBorrow(Ctx, manager.CreateBorrowParams{
		UserID:     pgtype.Int4{Int32: int32(userID), Valid: true},
		BookID:     pgtype.Int4{Int32: int32(bookID), Valid: true},
		BorrowedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
		ReturnedAt: pgtype.Timestamp{Valid: false},
	})

	if err != nil {
		fmt.Println("couldn't create a borrow into postgres. ", err)
	}

	borrow := models.Borrows{
		ID:         brw.ID, // or use DB-generated ID
		UserId:     brw.UserID.Int32,
		BookId:     brw.BookID.Int32,
		BorrowedAt: brw.BorrowedAt.Time.String(),
		ReturnedAt: "",
	}

	// Save to JSON or insert into Borrows table
	books := LoadBooks()
	borrows := LoadBorrows()
	for i, b := range books {
		if b.ID == ub.ID {
			books[i] = models.Book{
				ID:         ub.ID,
				Title:      ub.Title,
				Author:     ub.Author,
				Year:       int(ub.Year.Int32),
				IsBorrowed: ub.IsBorrowed.Bool,
			}
			break
		}
	}
	borrows = append(borrows, borrow)
	SaveBooks(books)
	SaveBorrows(borrows)

	fmt.Printf("Book borrowed successfully: %+v\n", ub)
}
