package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"projectOrder/internal/dto"
	"strconv"
)

type OrderRepo struct {
	db *sqlx.DB
}

func NewOrderRepo(db *sqlx.DB) *OrderRepo {
	return &OrderRepo{
		db: db}
}

func (o *OrderRepo) GetOrderById(orderID string) (*dto.Order, error) {
	var orderDB []dto.OrderDB

	query := fmt.Sprintf(`SELECT 
   ds.delivery_name AS delivery_name,
   d.zip AS delivery_zip, 
   d.city AS delivery_city, 
   d.address AS delivery_address, 
   d.region AS delivery_region,
   b.bank_name AS bank_name,
   c.currency_name AS currency_name,
   p.transaction AS payment_transaction,
   p.request_id AS payment_request, 
   p.provider AS payment_provider, 
   p.amount AS payment_amount, 
   p.payment_dt AS payment_payment_dt,
   p.delivery_cost AS payment_delivery_cost, 
   p.goods_total AS payment_goods_total,
   p.custom_fee AS payment_custom_fee,
   i.price AS items_price, 
   i.rid AS items_rid, 
   i.name AS items_name, 
   i.size AS items_size,
   i.nm_id AS items_nm_id, 
   i.brand AS items_brand,
   s.name_status AS status_name,
   cu.name AS customer_name, 
   cu.phone AS customer_phone, 
   cu.email AS customer_email,
   o.order_uid AS orders_uid, 
   o.track_number AS orders_track_number, 
   o.chrt_id AS orders_chrt_id,
   o.entry AS orders_entry, 
   o.locale AS orders_locale, 
   o.internal_signature AS orders_internal_signature,
   o.shardkey AS orders_shardkey, 
   o.sm_id AS orders_sm_id, 
   o.date_created AS orders_date_created, 
   o.off_shard AS orders_off_shard,
   oi.sale AS order_items_sale, 
   oi.total_price AS order_items_total_price
FROM order_items oi
        INNER JOIN orders o ON oi.order_id = o.id
        INNER JOIN items i ON oi.item_id = i.id
        INNER JOIN delivery d ON o.delivery_id = d.id
        INNER JOIN delivery_service ds ON d.delivery_service_id = ds.id
        INNER JOIN payment p ON o.payment_id = p.id
        INNER JOIN bank b ON p.bank_id = b.id
        INNER JOIN currency c ON p.currency_id = c.id
        INNER JOIN status s ON o.status_id = s.id
        INNER JOIN customer cu ON o.customer_id = cu.id
       WHERE o.order_uid = $1`)

	err := o.db.Select(&orderDB, query, orderID)

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	var items []*dto.Item

	for _, orderInfo := range orderDB {
		chrtID, err := strconv.Atoi(orderInfo.ChrtID)

		if err != nil {
			log.Println(err.Error())
			return nil, err
		}

		nmID, err := strconv.Atoi(orderInfo.NmID)

		if err != nil {
			log.Println(err.Error())
			return nil, err
		}

		items = append(items, &dto.Item{
			ChrtID:      chrtID,
			TrackNumber: orderInfo.TrackNumber,
			Price:       orderInfo.Price,
			Rid:         orderInfo.RID,
			Name:        orderInfo.ItemName,
			Sale:        orderInfo.Sale,
			Size:        strconv.Itoa(orderInfo.Size),
			TotalPrice:  orderInfo.TotalPrice,
			NmID:        nmID,
			Brand:       orderInfo.Brand,
			Status:      orderInfo.NameStatus,
		})
	}

	orderInfo := new(dto.Order)
	if len(orderDB) != 0 {
		orderInfo = &dto.Order{
			OrderUID:    orderDB[0].OrderUID,
			TrackNumber: orderDB[0].TrackNumber,
			Entry:       orderDB[0].Entry,
			Delivery: dto.DeliveryInfo{
				Name:    orderDB[0].DeliveryName,
				Phone:   orderDB[0].CustomerPhone,
				Zip:     orderDB[0].ZIP,
				City:    orderDB[0].City,
				Address: orderDB[0].Address,
				Region:  orderDB[0].Region,
				Email:   orderDB[0].CustomerEmail,
			},
			Payment: dto.PaymentInfo{
				Transaction:  orderDB[0].Transaction,
				RequestID:    orderDB[0].RequestID,
				Currency:     orderDB[0].CurrencyName,
				Provider:     orderDB[0].Provider,
				Amount:       orderDB[0].Amount,
				PaymentDt:    orderDB[0].PaymentDT.Unix(),
				Bank:         orderDB[0].BankName,
				DeliveryCost: orderDB[0].DeliveryCost,
				GoodsTotal:   orderDB[0].GoodsTotal,
				CustomFee:    orderDB[0].CustomFee,
			},
			Items:           items,
			Locale:          orderDB[0].Locale,
			InternalSign:    orderDB[0].InternalSignature,
			DeliveryService: orderDB[0].DeliveryName,
			ShardKey:        strconv.Itoa(orderDB[0].ShardKey),
			SmID:            orderDB[0].SmID,
			DateCreated:     orderDB[0].DateCreated,
			OofShard:        strconv.Itoa(orderDB[0].OffShard),
		}
	}

	return orderInfo, nil
}
