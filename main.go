package main

import (
	"Assignment2/Tables"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	host	= "localhost"
	port	= 5432
	user	= "postgres"
	password	= "admin"
	dbname		= "assignment2"
)

var (
	db *sql.DB
	err error
)

var PORT = ":8088"

var orders = map[int]Tables.Orders{
	// 1: {
	// 	Order_id: 1,
	// 	Customer_name: "Edrick",
	// 	Item: []Tables.Items{
	// 		{
	// 			Item_id: 1,
	// 			Item_code: 10,
	// 			Description: "pensil",
	// 			Quantity: 2,
	// 			OrderId: 1,
	// 		},
	// 	},
	// },
}

func main(){

	r := mux.NewRouter()
	r.HandleFunc("/orders", userHandler)
	r.HandleFunc("/orders/{Id}", userHandler)
	
	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8088",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	//connect database
	psqlInfo := fmt.Sprintf("host= %s port= %d user= %s "+" password= %s dbname= %s sslmode=disable", host, port, user, password, dbname) 
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil{
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("successfully Connect to Database")
	// insertOrders()

	log.Fatal(srv.ListenAndServe())
	
}

func userHandler(w http.ResponseWriter, r *http.Request){
	param := mux.Vars(r)
	id := param["Id"]

	if r.Method == "GET"{
		results := []Tables.Orders{}
		results = getData()
		fmt.Println(results)
		if id != "" {
			fmt.Println(id)
			if idInt, err := strconv.Atoi(id); err != nil{
				return
			} else{
				jsonData, _ :=json.Marshal(results[idInt])
				w.Header().Add("Content-Type", "application/json")
				w.Write(jsonData)
			}
		} else {
			jsonData, _ := json.Marshal(&results)
			w.Header().Add("Content-Type", "application/json")
			w.Write(jsonData)
		}
	}
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
    	var newOrders Tables.Orders
    	err := decoder.Decode(&newOrders)
		if err != nil {
			panic(err)
		}
		orders[int(newOrders.Order_id)] = newOrders
		//users = append(users, newUsers)
		// fmt.Println(orders)
		sliceOrders := []Tables.Orders{
		}
		for _, value := range orders {
			sliceOrders = append(sliceOrders, value)
		}
		// json.NewEncoder(w).Encode(sliceOrders)
		// return
		jsonData, _ := json.Marshal(&sliceOrders)
		w.Header().Add("Content-Type", "application/json")
		w.Write(jsonData)
		fmt.Printf("%+v\n",newOrders)
		insertOrders(newOrders)

		// var instOrder
	} 
	if r.Method == "PUT" {
		decoder := json.NewDecoder(r.Body)
    	var newOrders Tables.Orders
    	err := decoder.Decode(&newOrders)
		if err != nil {
			panic(err)
		}
		if id != "" {
			fmt.Println(id)
			if idInt, err := strconv.Atoi(id); err != nil{
				return
			} else{
				fmt.Println(idInt)
				newOrders.Order_id = idInt
				fmt.Printf("%+v\n",newOrders)
				updateOrders(newOrders)
				for i:=0;i<len(newOrders.Item);i++ {
					updateItems(newOrders.Item[i])
				}
			}
		}
	} 
	if r.Method == "DELETE" {
		if id != "" {
			if	index, err := strconv.Atoi(id); err != nil{
				return
			}else{
				delOrder := Tables.Orders{}
				delItems := Tables.Items{}
				delOrder.Order_id = index
				delItems.OrderId = index

				deleteItems(delItems.OrderId)
				deleteOrder(delOrder.Order_id)
			}
		}
	}
}

func insertOrders(neworder Tables.Orders) {
	// var insert = Tables.Orders{}

	sqlStatement := `
	INSERT INTO orders (customer_name, ordered_at)
	VALUES ($1, $2)
	returning order_id
	`

	rows, err := db.Query(sqlStatement, neworder.Customer_name, neworder.Ordered_at)
	if err != nil {
		panic (err)
	}

	for rows.Next() {
		err = rows.Scan(&neworder.Order_id)

		if err != nil {
			panic(err)
		}
	}
	fmt.Println(neworder.Order_id)
	fmt.Printf("New order data : %+v\n", neworder)
	// for _,v := range neworder.Item {
	// 	insertItems(v)
	// }
	for i:=0;i<len(neworder.Item);i++ {
		neworder.Item[i].OrderId = neworder.Order_id
		insertItems(neworder.Item[i])
	}
}

func insertItems(newitem Tables.Items) {
	// var insert = Tables.Orders{}

	sqlStatement := `
	INSERT INTO items (item_code, description, quantity, orderid)
	VALUES ($1, $2, $3, $4)
	`

	_, err := db.Exec(sqlStatement, newitem.Item_code, newitem.Description, newitem.Quantity, newitem.OrderId)
	if err != nil {
		panic (err)
	}
}

func getData() []Tables.Orders{
	resultsOrders := []Tables.Orders{}
	sqlStatement := `SELECT * from items JOIN orders ON (items.orderid = orders.order_id)`

	rows, err := db.Query(sqlStatement)

	if err != nil {
		panic(err)
	}
	defer rows.Close() 

	for rows.Next() {
		orders := Tables.Orders{}
		items := Tables.Items{}
		sliceItems := []Tables.Items{}

		err = rows.Scan(&items.Item_id, &items.Item_code, &items.Description, &items.Quantity, &items.OrderId, &orders.Order_id, &orders.Customer_name, &orders.Ordered_at)

		if err != nil {
			panic(err)
		}
		
		sliceItems = append(sliceItems, items)
		orders.Item = sliceItems
		resultsOrders = append(resultsOrders, orders)
	}

	fmt.Println("orders :", resultsOrders)
	return resultsOrders
}

func updateOrders(updateorders Tables.Orders){
	sqlStatement := `
	UPDATE orders 
	SET customer_name = $1,
    	ordered_at = $2
	WHERE order_id = $3;
	`
	_, err := db.Exec(sqlStatement, updateorders.Customer_name, updateorders.Ordered_at, updateorders.Order_id)
	if err != nil {
		panic (err)
	}
}

func updateItems(updateitems Tables.Items){
	sqlStatement := `
	UPDATE items 
	SET item_code = $1,
    	description = $2,
    	quantity = $3
	WHERE item_id = $4;
	`
	_, err := db.Exec(sqlStatement, updateitems.Item_code, updateitems.Description, updateitems.Quantity, updateitems.Item_id)
	if err != nil {
		panic (err)
	}
}

func deleteItems(orderId int){
	sqlStatement :=`
	DELETE FROM items
	WHERE orderid = $1
	`

	_, err := db.Exec(sqlStatement, orderId)
	if err != nil {
		panic (err)
	}
}

func deleteOrder(orderId int){
	sqlStatement :=`
	DELETE FROM orders
	WHERE order_id = $1
	`

	_, err := db.Exec(sqlStatement, orderId)
	if err != nil {
		panic (err)
	}
}
