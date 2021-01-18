package main

import (
	"context"
	"fmt"
	"gallery/server/database"
	"gallery/server/domain"
	"gallery/server/handler"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
}

func main() {

	var (
		databaseUser     = os.Getenv("DATABASE_USER")
		databaseName     = os.Getenv("DATABASE_NAME")
		databaseHost     = os.Getenv("DATABASE_HOST")
		databasePort     = os.Getenv("DATABASE_PORT")
		databasePassword = os.Getenv("DATABASE_PASSWORD")
	)

	dbConn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", databaseHost, databasePort, databaseUser, databaseName, databasePassword)

	conn := database.OpenDB(dbConn)

	appAddr := ":" + os.Getenv("APP_PORT")

	r := gin.Default()

	imgDomain := domain.NewImageService(conn)
	userDomain := domain.NewUserService(conn)

	h := handler.NewHandlerService(imgDomain, userDomain)

	err := userDomain.SeedUsers()
	if err != nil {
		log.Fatal("Error seeding users: ", err)
	}

	err = imgDomain.SeedImages()
	if err != nil {
		log.Fatal("Error seeding images: ", err)
	}

	r.POST("/login", h.Login)
	r.DELETE("/images/:imageId", h.DeleteImage)
	r.DELETE("/bulk_delete/images", h.BulkDeleteImage)

	//Starting and Shutting down Server
	srv := &http.Server{
		Addr:    appAddr,
		Handler: r,
	}

	go func() {
		//service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	//Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Println("Server exiting")

}
