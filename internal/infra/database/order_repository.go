package database

import (
	"database/sql"
	"github.com/AlexandreLima658/go-intensivo-jul/internal/entity"
)

type OrderRepository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{
		DB: db,
	}
}

func (r *OrderRepository) Save(order *entity.Order) error {
	_, err := r.DB.Exec("INSERT INTO orders(id, price, tax, final_price) VALUES ($1, $2, $3, $4)",
		order.ID, order.Price, order.Tax, order.FinalPrice)
	if err != nil {
		return err
	}
	return nil
}

func (r *OrderRepository) GetTotalTransactions() (int, error) {
	var total int
	err := r.DB.QueryRow("SELECT COUNT(*) FROM orders").Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}
