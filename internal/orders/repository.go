package orders

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type OrdersRepository interface {
	FetchOrder(id uuid.UUID) (*DetailedOrder, error)
	FetchAllOrders() ([]*Order, error)
	CreateOrder(id uuid.UUID, userId uuid.UUID, totalQuantity int, priceTotal int, products []*DetailedOrderProduct) error
}

type ordersRepository struct {
	db *sql.DB
}

func NewOrdersRepository(db *sql.DB) OrdersRepository {
	return &ordersRepository{
		db: db,
	}
}

const (
	fetchOrderQuery         = "SELECT user_id, total_quantity, price_total, created_at FROM orders WHERE id = $1"
	fetchOrderProductsQuery = "SELECT product_id, price, quantity FROM order_products WHERE order_id = $1"

	fetchAllOrdersQuery = "SELECT id, user_id, total_quantity, price_total, created_at FROM orders"

	createOrderQuery        = "INSERT INTO orders(id, user_id, total_quantity, price_total, created_at) VALUES($1, $2, $3, $4, $5)"
	createOrderProductQuery = "INSERT INTO order_products(id, order_id, product_id, price, quantity) VALUES($1, $2, $3, $4, $5)"
)

func (r *ordersRepository) FetchOrder(id uuid.UUID) (*DetailedOrder, error) {
	order := DetailedOrder{Id: id}

	err := r.db.QueryRow(fetchOrderQuery, id).Scan(
		&order.UserId,
		&order.TotalQuantity,
		&order.TotalPrice,
		&order.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(fetchOrderProductsQuery, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		product := new(DetailedOrderProduct)
		err = rows.Scan(
			&product.ProductId,
			&product.Price,
			&product.Quantity,
		)
		if err != nil {
			return nil, err
		}
		order.Products = append(order.Products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *ordersRepository) FetchAllOrders() ([]*Order, error) {
	orders := []*Order{}
	rows, err := r.db.Query(fetchAllOrdersQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		order := new(Order)
		err = rows.Scan(
			&order.Id,
			&order.UserId,
			&order.TotalQuantity,
			&order.TotalPrice,
			&order.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *ordersRepository) CreateOrder(orderId uuid.UUID, userId uuid.UUID, totalQuantity int, priceTotal int, products []*DetailedOrderProduct) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		createOrderQuery,
		orderId,
		userId,
		totalQuantity,
		priceTotal,
		time.Now().UTC(),
	)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	for _, v := range products {
		rowId, err := uuid.NewV7()
		if err != nil {
			_ = tx.Rollback()
			return err
		}

		_, err = tx.Exec(
			createOrderProductQuery,
			rowId,
			orderId,
			v.ProductId,
			v.Price,
			v.Quantity,
		)
		if err != nil {
			_ = tx.Rollback()
			return err
		}

	}
	err = tx.Commit()
	return err
}
