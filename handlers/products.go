package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

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

// ServeHTTP implements the go http.handler interface
// interface
func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// handle the request for a list of products
	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
		return
	}

	// Handle update
	if r.Method == http.MethodPost {
		p.addProduct(rw, r)
		return
	}

	//
	if r.Method == http.MethodPut {
		// expect ID in URL
		reg := regexp.MustCompile(`/([0-9]+)`)
		g := reg.FindAllStringSubmatch(r.URL.Path, -1)

		if len(g) != 1 {
			p.l.Println("Invalid ID - more than one ID")
			http.Error(rw, "Invalid ID", http.StatusBadRequest)
			return
		}

		if len(g[0]) != 2 {
			p.l.Println("Invalid ID - more than one capture group")
			http.Error(rw, "Invalid ID", http.StatusBadRequest)
			return
		}

		idString := g[0][1]
		id, err := strconv.Atoi(idString)

		if err != nil {
			p.l.Println("Unable to convert to string", idString)
			http.Error(rw, "Unable to convert the string", http.StatusNotAcceptable)
		}

		p.updateProduct(id, rw, r)

	}

	// catch all
	// if no method is satisfied return an error
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

// getProducts returns the products from the data store
func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {

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

func (p *Products) addProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle Post Product", r.URL.Path)

	product := &data.Product{}
	err := product.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to parse Json data", http.StatusBadRequest)
	}

	data.AddProduct(product)
}

func (p *Products) updateProduct(id int, rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT Product", r.URL.Path)

	product := &data.Product{}
	p.l.Println("New Product =========", product)
	err := product.FromJSON(r.Body)

	if err != nil {
		http.Error(rw, "Unable to parse Json data", http.StatusBadRequest)
	}

	err = data.UpdateProduct(id, product)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product Not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product Not found", http.StatusInternalServerError)
		return
	}
}
