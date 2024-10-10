package products

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type ProductsRepository interface {
	FetchProduct(id uuid.UUID) (*DetailedProduct, error)
	FetchAllProducts() ([]*DetailedProduct, error)
	AddProduct(id uuid.UUID, name string, desc string, price int64, stock int64) error

	UpdateInventoryAdd(id uuid.UUID, num int64) error
	UpdateInventorySubtract(id uuid.UUID, num int64) error
	UpdateInventoryDelete(id uuid.UUID) error
}

type productsRepository struct {
	db *sql.DB
}

func NewProductsRepository(db *sql.DB) ProductsRepository {
	return &productsRepository{
		db: db,
	}
}

const (
	fetchProductQuery     = "SELECT name, description, price, stock, created_at FROM products WHERE id = $1"
	fetchAllProductsQuery = "SELECT id, name, description, price, stock, created_at FROM products"
	addProductQuery       = "INSERT INTO products(id, name, description, price, stock, created_at) VALUES($1, $2, $3, $4, $5, $6)"

	updateInventoryLockQuery     = "SELECT stock FROM products WHERE id = $1 FOR UPDATE"
	updateInventoryAddQuery      = "UPDATE products SET stock = stock + $1 WHERE id = $2"
	updateInventorySubtractQuery = "UPDATE products SET stock = stock - $1 WHERE id = $2"
	updateInventoryDeleteQuery   = "DELETE FROM products WHERE id = $1"
)

func (r *productsRepository) FetchProduct(id uuid.UUID) (*DetailedProduct, error) {
	product := &DetailedProduct{Id: id}
	err := r.db.QueryRow(fetchProductQuery, id).Scan(
		&product.Name,
		&product.Desc,
		&product.Price,
		&product.Stock,
		&product.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (r *productsRepository) FetchAllProducts() ([]*DetailedProduct, error) {
	products := []*DetailedProduct{}
	rows, err := r.db.Query(fetchAllProductsQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		product := new(DetailedProduct)
		err = rows.Scan(
			&product.Id,
			&product.Name,
			&product.Desc,
			&product.Price,
			&product.Stock,
			&product.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *productsRepository) AddProduct(id uuid.UUID, name string, desc string, price int64, stock int64) error {
	_, err := r.db.Exec(
		addProductQuery,
		id,
		name,
		desc,
		price,
		stock,
		time.Now().UTC(),
	)
	return err
}

func (r *productsRepository) UpdateInventoryAdd(id uuid.UUID, num int64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// Lock row
	_, err = tx.Exec(updateInventoryLockQuery, id)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	// Exec update
	_, err = tx.Exec(
		updateInventoryAddQuery,
		num,
		id,
	)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	err = tx.Commit()
	return err
}

func (r *productsRepository) UpdateInventorySubtract(id uuid.UUID, num int64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// Lock row
	_, err = tx.Exec(updateInventoryLockQuery, id)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	// Exec update
	_, err = tx.Exec(
		updateInventorySubtractQuery,
		num,
		id,
	)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	err = tx.Commit()
	return err
}

func (r *productsRepository) UpdateInventoryDelete(id uuid.UUID) error {
	_, err := r.db.Exec(
		updateInventoryDeleteQuery,
		id,
	)
	return err
}
