package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	_ "github.com/lib/pq"
)

type BlogPost struct {
	ID      	string `json:"id"`
	Title   	string `json:"title"`
	Content 	string `json:"content"`
	CreatedAt 	string `json:"created_at"`
	UpdatedAt 	string `json:"updated_at"`
	Author		string `json:"author"`
}

var db *sql.DB

const (
	dbSQL  = "postgres"
	dbUser = "hagavi"
	dbPass = "heartgatavirus"
	dbName = "hagavi_blog"
	blogPostTable = "blog_posting"
	port   = ":8080"
)

func createBlogPostHandler(w http.ResponseWriter, r *http.Request) {
	var newPost BlogPost
	err := json.NewDecoder(r.Body).Decode(&newPost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	createdTime := time.Now()
	sqlStatementCreate := fmt.Sprintf(`INSERT INTO %s (title, content, created_at, updated_at, author) VALUES ($1, $2, $3, $4, $5)`, blogPostTable)
	_, err = db.Exec(sqlStatementCreate, newPost.Title, newPost.Content, createdTime, createdTime, newPost.Author)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(&newPost)
}

func getBlogPostHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Post ID Not Found", http.StatusNotFound)
		return
	}
	var post BlogPost
	sqlStatementGet := fmt.Sprintf(`SELECT id, title, content, created_at, updated_at, author FROM %s WHERE id = $1`, blogPostTable)
	rowResult := db.QueryRow(sqlStatementGet, id)
	errorResult := rowResult.Scan(&post.ID, &post.Title, &post.Content,&post.CreatedAt, &post.UpdatedAt, &post.Author)
	switch errorResult {
	case sql.ErrNoRows:
		http.Error(w, errorResult.Error(), http.StatusInternalServerError)
	case nil:
		json.NewEncoder(w).Encode(post)
	default:
		http.Error(w, errorResult.Error(), http.StatusInternalServerError)
	}
	return
}

func updateBlogPostHandler(w http.ResponseWriter, r *http.Request) {
	var updatedPost BlogPost
	err := json.NewDecoder(r.Body).Decode(&updatedPost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	updatedTime := time.Now()
	sqlStatementUpdated := fmt.Sprintf(`UPDATE %s SET title=$1, content=$2, created_at=$3, updated_at=$4, author=$5 WHERE id=$6`, blogPostTable)
	_, err = db.Exec(sqlStatementUpdated, updatedPost.Title, updatedPost.Content, updatedPost.CreatedAt, updatedTime, updatedPost.Author, updatedPost.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(&updatedPost)
}

func deleteBlogPostHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing Post ID Parameter", http.StatusBadRequest)
		return
	}
	sqlStatementDelete := fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, blogPostTable)
	_, err := db.Exec(sqlStatementDelete, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func getAllBlogPostHandler(w http.ResponseWriter, r *http.Request) {
	sqlStatementGetAll := fmt.Sprintf(`SELECT id, title, content, created_at, updated_at, author FROM %s`, blogPostTable)
	rows, err := db.Query(sqlStatementGetAll)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	blog_posting := []BlogPost{}
	for rows.Next() {
		var post BlogPost
		err := rows.Scan(&post.ID, &post.Title, &post.Content,&post.CreatedAt, &post.UpdatedAt, &post.Author)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		blog_posting = append(blog_posting, post)
	}
	json.NewEncoder(w).Encode(blog_posting)
}

func main() {
	var err error
	sqlOpeningStatement := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbUser, dbPass, dbName)
	db, err = sql.Open(dbSQL, sqlOpeningStatement)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	http.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		switch r.Method {
		case "POST":
			createBlogPostHandler(w, r)
		case "GET":
			if id == "" {
				getAllBlogPostHandler(w, r)
			} else {
				getBlogPostHandler(w, r)
			}
		case "PUT":
			updateBlogPostHandler(w, r)
		case "DELETE":
			deleteBlogPostHandler(w, r)
		default:
			http.Error(w, "Method Undefined", http.StatusMethodNotAllowed)
		}
	})
	fmt.Println("Server listening on port " + port)
	http.ListenAndServe(port, nil)
}
