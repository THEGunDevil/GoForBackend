package main

import (
	"LibraryManager/internal/app"
	"fmt"
)

func main() {
	fmt.Println("$$$$$$$ Welcome to the Library!!$$$$$$$")
	fmt.Println("$$$$$$$ Please select a service from down below!!$$$$$$$")
	fmt.Println("1. Add a Book")
	fmt.Println("2. Delete a Book")
	fmt.Println("3. Borrow a Book")
	fmt.Println("4. Show all Books")
	fmt.Println("5. Return a borrowed Book")
	fmt.Println("6. Search a book")
	fmt.Println("7. Add a new user")
	fmt.Println("8. Show all user")
	fmt.Println("9. Filter a user by email or name")
	fmt.Println("0. exit")
	app.SelectService()
	app.FilesReload()

}
