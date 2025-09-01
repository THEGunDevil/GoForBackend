-- name: CreateBook :one
INSERT INTO books (title, author, year) VALUES ($1, $2, $3)
RETURNING id, title, author, year, is_borrowed;

-- name: CreateUser :one
INSERT INTO users (name, email) VALUES ($1, $2)
RETURNING id, name, email;

-- name: CreateBorrow :one
INSERT INTO borrows (user_id, book_id, borrowed_at, returned_at) VALUES ($1, $2, $3, $4)
RETURNING id, user_id, book_id, borrowed_at, returned_at;

-- name: ListBooks :many
SELECT * FROM books ORDER BY id;
-- name: ListBorrows :many
SELECT * FROM borrows ORDER BY id;
-- name: ListUsers :many
SELECT * FROM users ORDER BY id;

-- name: UpdateBook :one
UPDATE books
SET 
    title       = COALESCE($2, title),
    author      = COALESCE($3, author),
    year        = COALESCE($4, year),
    is_borrowed = COALESCE($5, is_borrowed)
WHERE id = $1
RETURNING id, title, author, year, is_borrowed;




-- name: ReturnBook :exec
UPDATE books SET is_borrowed = false WHERE id = $1;

-- name: DeleteBookWithId :execrows
DELETE FROM books WHERE id = $1;

-- name: DeleteBookWithTitle :execrows
DELETE FROM books WHERE title = $1;

-- name: DeleteBorrow :one
DELETE FROM borrows
WHERE book_id = $1 AND user_id = $2
RETURNING *;


-- name: SearchBooks :many
SELECT id, title, author, year, is_borrowed
FROM books
WHERE title ILIKE '%' || $1 || '%' OR author ILIKE '%' || $1 || '%';
