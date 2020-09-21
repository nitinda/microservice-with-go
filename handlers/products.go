package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/nitinda/microservice-with-go/data"
)

// Products is a http.Handler
type Products struct {
	l *log.Logger
}

// NewProducts creates a products handler with the given logger
func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// GetProducts returns the products from the data store
func (p Products) GetProducts(rw http.ResponseWriter, r *http.Request) {

	p.l.Println("Handle Get Products", r.URL.Path)

	// fetch the products from the datastore
	lp := data.GetProducts()

	// serialize the list to JSON
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to parse Json data", http.StatusInternalServerError)
		return
	}
}

func (p Products) AddProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle Post Product", r.URL.Path)

	product := r.Context().Value(KeyProduct{}).(data.Product)
	data.AddProduct(&product)
}

func (p *Products) UpdateProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT Product", r.URL.Path)

	vars := mux.Vars(r)
	id, er := strconv.Atoi(vars["id"])
	if er != nil {
		p.l.Println("Unable to convert the integer to string")
		http.Error(rw, "Unable to convert the integer to string", http.StatusBadRequest)
		return
	}

	product := r.Context().Value(KeyProduct{}).(data.Product)
	err := data.UpdateProduct(id, &product)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product Not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product Not found", http.StatusInternalServerError)
		return
	}
}

type KeyProduct struct{}

func (p Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		product := data.Product{}
		p.l.Println("Procunt ====== ", product)
		err := product.FromJSON(r.Body)
		if err != nil {
			p.l.Println("[ERROR] deserializing product", err)
			http.Error(rw, "Unable to parse Json data", http.StatusBadRequest)
			return
		}

		// add the product to the context
		ctx := context.WithValue(r.Context(), KeyProduct{}, product)
		req := r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, req)
	})
}
