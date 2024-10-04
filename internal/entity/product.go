package entity

import "github.com/google/uuid"

type ProductRepository interface {
	Create(product *Product) error //se der erro retorna erro, se não der o erro vem vázio
	FindAll() ([]*Product, error)  //Devolve uma lista de produtos ou um erro
}

type Product struct {
	ID    string
	Name  string
	Price float64
}

func NewProduct(name string, price float64) *Product {
	return &Product{
		ID:    uuid.New().String(),
		Name:  name,
		Price: price,
	}
}
