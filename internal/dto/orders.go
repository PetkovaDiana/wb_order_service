package dto

import "time"

type DeliveryInfo struct {
	Name    string `json:"name,omitempty"`
	Phone   string `json:"phone,omitempty"`
	Zip     string `json:"zip,omitempty"`
	City    string `json:"city,omitempty"`
	Address string `json:"address,omitempty"`
	Region  string `json:"region,omitempty"`
	Email   string `json:"email,omitempty"`
}

type PaymentInfo struct {
	Transaction  string `json:"transaction,omitempty"`
	RequestID    string `json:"request_id,omitempty"`
	Currency     string `json:"currency,omitempty"`
	Provider     string `json:"provider,omitempty"`
	Amount       int    `json:"amount,omitempty"`
	PaymentDt    int64  `json:"payment_dt,omitempty"`
	Bank         string `json:"bank,omitempty"`
	DeliveryCost int    `json:"delivery_cost,omitempty"`
	GoodsTotal   int    `json:"goods_total,omitempty"`
	CustomFee    int    `json:"custom_fee,omitempty"`
}

type Item struct {
	ChrtID      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmID        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

type Order struct {
	OrderUID        string       `json:"order_uid"`
	TrackNumber     string       `json:"track_number"`
	Entry           string       `json:"entry"`
	Delivery        DeliveryInfo `json:"delivery"`
	Payment         PaymentInfo  `json:"payment"`
	Items           []*Item      `json:"items"`
	Locale          string       `json:"locale"`
	InternalSign    string       `json:"internal_signature"`
	CustomerID      string       `json:"customer_id"`
	DeliveryService string       `json:"delivery_service"`
	ShardKey        string       `json:"shardkey"`
	SmID            int          `json:"sm_id"`
	DateCreated     time.Time    `json:"date_created"`
	OofShard        string       `json:"oof_shard"`
}

type DeliveryDB struct {
	DeliveryName string `db:"delivery_name"`
	ZIP          string `db:"delivery_zip"`
	City         string `db:"delivery_city"`
	Address      string `db:"delivery_address"`
	Region       string `db:"delivery_region"`
}

type BankDB struct {
	BankName     string `db:"bank_name"`
	CurrencyName string `db:"currency_name"`
}

type OrderDB struct {
	DeliveryDB
	BankDB
	Transaction       string    `db:"payment_transaction"`
	RequestID         string    `db:"payment_request"`
	Provider          string    `db:"payment_provider"`
	Amount            int       `db:"payment_amount"`
	PaymentDT         time.Time `db:"payment_payment_dt"`
	DeliveryCost      int       `db:"payment_delivery_cost"`
	GoodsTotal        int       `db:"payment_goods_total"`
	CustomFee         int       `db:"payment_custom_fee"`
	Price             int       `db:"items_price"`
	RID               string    `db:"items_rid"`
	ItemName          string    `db:"items_name"`
	Size              int       `db:"items_size"`
	NmID              string    `db:"items_nm_id"`
	Brand             string    `db:"items_brand"`
	NameStatus        int       `db:"status_name"`
	CustomerName      string    `db:"customer_name"`
	CustomerPhone     string    `db:"customer_phone"`
	CustomerEmail     string    `db:"customer_email"`
	OrderUID          string    `db:"orders_uid"`
	TrackNumber       string    `db:"orders_track_number"`
	ChrtID            string    `db:"orders_chrt_id"`
	Entry             string    `db:"orders_entry"`
	Locale            string    `db:"orders_locale"`
	InternalSignature string    `db:"orders_internal_signature"`
	ShardKey          int       `db:"orders_shardkey"`
	SmID              int       `db:"orders_sm_id"`
	DateCreated       time.Time `db:"orders_date_created"`
	OffShard          int       `db:"orders_off_shard"`
	Sale              int       `db:"order_items_sale"`
	TotalPrice        int       `db:"order_items_total_price"`
}
