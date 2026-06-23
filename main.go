// Package main starts the web server and hooks up the database.
package main

import (
	"accounts"
	"blogs"
	"comments"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// load the .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, relying on system env variables")
	}

	// env variables
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Africa/Nairobi", host, user, password, dbName, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	pgDB, err := db.DB()

	pgDB.SetMaxIdleConns(10)
	pgDB.SetMaxOpenConns(100)
	pgDB.SetConnMaxLifetime(time.Hour)

	if err != nil {
		// If it fails to connect, stop the application immediately
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	if err := db.AutoMigrate(&accounts.Users{}, &blogs.Blogs{}, &blogs.BlogImages{}, &comments.Comments{}); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	// Configure global log settings once at application startup
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Create the router
	r := mux.NewRouter()

	// configure cors
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	// Wrap your router with the middleware
	wrappedMux := LoggingMiddleware(handler)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintln(w, "Welcome to the Golang Server: "+r.URL.Path)
	})

	// initialize the service layer with the db instance
	userService := accounts.UserService{Db: db}
	userController := accounts.UserController{Service: userService}

	blogService := blogs.BlogServices{Db: db}
	blogController := blogs.BlogController{Service: blogService}

	commentService := comments.CommentService{Db: db}
	commentController := comments.CommentController{Service: commentService}

	r.HandleFunc("/api/v1/users/register", userController.RegisterUsers).Methods("POST")
	r.HandleFunc("/api/v1/users/login", userController.Login).Methods("POST")

	// create protected routes
	protected := r.PathPrefix("/api/v1").Subrouter()
	protected.Use(AuthenticationMiddleware)

	protected.HandleFunc("/blog/create", blogController.CreateBlog).Methods("POST")
	protected.HandleFunc("/blog/edit/{blogID}", blogController.EditBlog).Methods("PUT")
	protected.HandleFunc("/comment/{blogId}/blog", commentController.CreateComment).Methods("POST")

	fmt.Println("Server starting on port :8080...")

	if err := http.ListenAndServe(":8080", wrappedMux); err != nil {
		log.Fatal("there was a problem starting the server")
	}
}
