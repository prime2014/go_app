package main

import (
	"accounts"
	"blogs"
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
	err := godotenv.Load()

	// env variables
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	db_name := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Africa/Nairobi", host, user, password, db_name, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	pgDB, err := db.DB()

	pgDB.SetMaxIdleConns(10)
	pgDB.SetMaxOpenConns(100)
	pgDB.SetConnMaxLifetime(time.Hour)

	if err != nil {
		// If it fails to connect, stop the application immediately
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	db.AutoMigrate(&accounts.Users{})
	db.AutoMigrate(&blogs.Blogs{})
	db.AutoMigrate(&blogs.BlogImages{})

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
	// wrappedMux = AuthenticationMiddleware(wrappedMux)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome to the Golang Server")
	})

	// initialize the service layer with the db instance
	userService := accounts.UserService{Db: db}
	userController := accounts.UserController{Service: userService}

	blogService := blogs.BlogServices{Db: db}
	blogController := blogs.BlogController{Service: blogService}

	r.HandleFunc("/api/v1/users/register", userController.RegisterUsers).Methods("POST")
	r.HandleFunc("/api/v1/users/login", userController.Login).Methods("POST")
	r.HandleFunc("/api/v1/blog/create", blogController.CreateBlog).Methods("POST")

	fmt.Println("Server starting on port :8080...")

	if err := http.ListenAndServe(":8080", wrappedMux); err != nil {
		log.Fatal("there was a problem starting the server")
	}
}
