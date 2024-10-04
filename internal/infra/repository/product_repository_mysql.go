package repository

import (
	"database/sql"
	"fmt"
	"log"
	"productSystem/internal/entity"
)

type ProductRepositoryMysql struct {
	DB *sql.DB
}

func NewProductRepositoryMysql(db *sql.DB) *ProductRepositoryMysql {
	return &ProductRepositoryMysql{DB: db}
}

func (r *ProductRepositoryMysql) Create(product *entity.Product) error {
	_, err := r.DB.Exec("Insert into products (id, name, price) values(?,?,?)", product.ID, product.Name, product.Price)

	if err != nil {
		// Usando fmt.Printf para imprimir o erro de forma legível
		fmt.Printf("Erro ao inserir produto: %v\n", err)
		// Ou usando log.Printf para registrar o erro
		log.Printf("Erro ao inserir produto: %v", err)
		return err
	}

	return nil
}

func (r *ProductRepositoryMysql) FindAll() ([]*entity.Product, error) {
	rows, err := r.DB.Query("select id, name, price from products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*entity.Product
	for rows.Next() {
		var product entity.Product
		err = rows.Scan(&product.ID, &product.Name, &product.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, &product)

	}

	return products, nil
}
