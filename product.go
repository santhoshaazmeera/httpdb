package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// Product model
type Product struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

var products []Product
var db *sql.DB

func createdb() {
	connectiondetails := "user=santhosha dbname=mydb_http password=santhosha@123 "
	vdb, err := sql.Open("postgres", connectiondetails)
	db = vdb
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("connected to database !")
	// defer db.Close()

	//return db

}

func GetProducts(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(products)
	rows, err := db.Query("select * from httpdata")
	if err != nil {
		fmt.Println(err, "error while quering")
	}
	defer rows.Close()
	var data []string
	for rows.Next() {
		var id string
		var name string
		var price float64
		err := rows.Scan(&id, &name, &price)

		data = append(data, fmt.Sprintf("id: %s \nname : %s\nprice :%f", id, name, price))
		if err != nil {
			fmt.Println(err, "at scanning")
		} else {
			w.Write([]byte("retrived data succesfully..\n"))
		}
	}

	fmt.Println(data)

}

func AddProduct(w http.ResponseWriter, r *http.Request) {
	var newProduct Product

	fmt.Printf("\n%+v\n", r.Body)
	err := json.NewDecoder(r.Body).Decode(&newProduct)
	if err != nil {
		fmt.Println(err, "errrorr!")
	}
	//fmt.Println(newProduct)
	// products = append(products, newProduct)

	// json.NewEncoder(w).Encode(newProduct)
	// // for _, p := range products {
	// // 	fmt.Println(p)
	// // }
	// body, err := io.ReadAll(r.Body)
	// if err != nil {
	// 	http.Error(w, "Error while reading r body", http.StatusInternalServerError)
	// 	return
	// }
	// fmt.Println("printing the  Body:", string(body))

	// err = json.Unmarshal(body, &newProduct)
	// if err != nil {
	// 	fmt.Println(err,"error while unmarshalling")
	// 	http.Error(w, "Error decoding JSON data", http.StatusBadRequest)
	// 	return
	// }
	//rows :=
	//fmt.Println(newProduct.ID, newProduct.Price, newProduct.Name)
	_, err = db.Exec(`insert into httpdata(id,name,price) values($1,$2,$3)`, newProduct.ID, newProduct.Name, newProduct.Price)
	if err != nil {
		log.Fatal("error while inserting", err)
	} else {
		fmt.Println("data inserted ")
		w.Write([]byte("data posted OR inserted succesfully..\n"))
	}

}

func deleting(w http.ResponseWriter, r *http.Request) {
	prdct := mux.Vars(r)
	productId := prdct["id"]
	// for i, pr := range products {

	// 	if pr.ID == productId {
	// 		products = append(products[:i], products[i+1:]...)
	// 		//w.WriteHeader(http.StatusNoContent)
	// 	}
	// }
	_, err := db.Exec("delete from httpdata where id=$1 ", productId)
	if err != nil {
		fmt.Println(err, "error at execution")
	} else {
		fmt.Print("the specified data has been deleted")
		w.Write([]byte("the data deleted of specified field ."))
	}
}

func updating(w http.ResponseWriter, r *http.Request) {
	var newupdate Product
	err := json.NewDecoder(r.Body).Decode(&newupdate)
	//fmt.Println(newupdate,"after json")
	if err != nil {
		fmt.Println(err)
	}
	//newupdate.ID = "9000"
	//fmt.Println(newupdate)
	productvariables := mux.Vars(r)
	//fmt.Println(productvariables, "productvariables")
	prodctuadate := productvariables["id"]
	// for i, p := range products {
	// 	if p.ID == prodctuadate {
	// 		products[i] = newupdate
	// 		json.NewEncoder(w).Encode(newupdate)
	// 	}

	// }
	//fmt.Println(prodctuadate)

	_, err = db.Exec("update httpdata set id = $1  where id=$2", newupdate.ID, prodctuadate)
	if err != nil {
		log.Fatal("the error is because of :", err.Error())
	} else {
		fmt.Println("updated successfully!")
		w.Write([]byte("data updated succesfully..\n"))
	}

}

func main() {
	createdb()
	// defer db.Close()

	// products = append(products, Product{ID: "1", Name: "Product A", Price: 19.99})
	// products = append(products, Product{ID: "2", Name: "Product B", Price: 29.99})

	router := mux.NewRouter()

	router.HandleFunc("/products", GetProducts).Methods("GET")
	router.HandleFunc("/products", AddProduct).Methods("POST")
	router.HandleFunc("/products/{id}", deleting).Methods("DELETE")
	router.HandleFunc("/products/{id}", updating).Methods("PUT")

	http.ListenAndServe(":8080", router)
}
