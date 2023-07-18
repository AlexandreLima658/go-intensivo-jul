package main

import (
	"database/sql"
	"fmt"
	"github.com/AlexandreLima658/go-intensivo-jul/internal/infra/database"
	"github.com/AlexandreLima658/go-intensivo-jul/internal/usecase"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	password := os.Getenv("DB_PASSWORD")

	stringConnection := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", host, port, user, dbname, password)
	db, err := sql.Open("postgres", stringConnection)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	orderRepository := database.NewRepository(db)
	uc := usecase.NewCalculateFinalPrice(orderRepository)
	input := usecase.OrderInput{
		ID:    "1",
		Price: 25.0,
		Tax:   1.0,
	}
	output, err := uc.Execute(input)
	if err != nil {
		panic(err)
	}
	fmt.Println(output)
}
