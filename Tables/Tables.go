package Tables

import "time"

type Orders struct {
	Order_id      int
	Customer_name string
	Ordered_at    time.Time
	Item          []Items
}

type Items struct {
	Item_id     int
	Item_code   int
	Description string
	Quantity    int
	OrderId     int
}
