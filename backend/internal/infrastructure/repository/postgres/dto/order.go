package dto

import "time"

type OrderModel struct {
	OrderUID          string    `db:"order_uid"`
	TrackNumber       string    `db:"track_number"`
	Entry             string    `db:"entry"`
	Locale            string    `db:"locale"`
	InternalSignature string    `db:"internal_signature"`
	CustomerID        string    `db:"customer_id"`
	DeliveryService   string    `db:"delivery_service"`
	ShardKey          string    `db:"shardkey"`
	SmID              int       `db:"sm_id"`
	DateCreated       time.Time `db:"date_created"`
	OofShard          string    `db:"oof_shard"`

	DeliveryName    string `db:"delivery_name"`
	DeliveryPhone   string `db:"delivery_phone"`
	DeliveryZip     string `db:"delivery_zip"`
	DeliveryCity    string `db:"delivery_city"`
	DeliveryAddress string `db:"delivery_address"`
	DeliveryRegion  string `db:"delivery_region"`
	DeliveryEmail   string `db:"delivery_email"`

	PaymentTransaction  string `db:"payment_transaction"`
	PaymentRequestID    string `db:"payment_request_id"`
	PaymentCurrency     string `db:"payment_currency"`
	PaymentProvider     string `db:"payment_provider"`
	PaymentAmount       int    `db:"payment_amount"`
	PaymentDT           int64  `db:"payment_dt"`
	PaymentBank         string `db:"payment_bank"`
	PaymentDeliveryCost int    `db:"payment_delivery_cost"`
	PaymentGoodsTotal   int    `db:"payment_goods_total"`
	PaymentCustomFee    int    `db:"payment_custom_fee"`
}
