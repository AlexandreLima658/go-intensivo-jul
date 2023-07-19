package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/AlexandreLima658/go-intensivo-jul/internal/infra/database"
	"github.com/AlexandreLima658/go-intensivo-jul/internal/usecase"
	"github.com/AlexandreLima658/go-intensivo-jul/pkg/rabbitmq"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
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
	ch, err := rabbitmq.OpenChannel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()
	msgRabbitmqChannel := make(chan amqp.Delivery)
	go rabbitmq.Consume(ch, msgRabbitmqChannel)
	rabbitmqWorker(msgRabbitmqChannel, uc)
}

func rabbitmqWorker(msgChan chan amqp.Delivery, uc *usecase.CalculateFinalPrice) {
	fmt.Println("Starting RabbitMQ Worker")
	for msg := range msgChan {
		var input usecase.OrderInput
		err := json.Unmarshal(msg.Body, &input)
		if err != nil {
			panic(err)

		}
		output, err := uc.Execute(input)
		if err != nil {
			panic(err)

		}
		fmt.Println("Mensagem processada e salva no banco ", output)
	}
}
