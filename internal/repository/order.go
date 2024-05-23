package repository

import (
	"context"
	"database/sql"
	"errors"
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
		db: db,
	}
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
   oi.total_price AS order_items_total_price,
   cu.external_customer_id AS customer_external_customer_id
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

	if len(orderDB) == 0 {
		return nil, errors.New("some error")
	}

	var items []*dto.Item

	for _, orderInfo := range orderDB {
		orderInfo := orderInfo
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

		size := strconv.Itoa(orderInfo.Size)

		items = append(items, &dto.Item{
			ChrtID:      &chrtID,
			TrackNumber: &orderInfo.TrackNumber,
			Price:       &orderInfo.Price,
			Rid:         &orderInfo.RID,
			Name:        &orderInfo.ItemName,
			Sale:        &orderInfo.Sale,
			Size:        &size,
			TotalPrice:  &orderInfo.TotalPrice,
			NmID:        &nmID,
			Brand:       &orderInfo.Brand,
			Status:      &orderInfo.NameStatus,
		})
	}

	orderInfo := new(dto.Order)
	if len(orderDB) != 0 {
		shardKey := strconv.Itoa(orderDB[0].ShardKey)
		offShard := strconv.Itoa(orderDB[0].OffShard)

		orderInfo = &dto.Order{
			OrderUID:    &orderDB[0].OrderUID,
			TrackNumber: &orderDB[0].TrackNumber,
			Entry:       &orderDB[0].Entry,
			Delivery: dto.DeliveryInfo{
				Name:    &orderDB[0].CustomerName,
				Phone:   &orderDB[0].CustomerPhone,
				Zip:     &orderDB[0].ZIP,
				City:    &orderDB[0].City,
				Address: &orderDB[0].Address,
				Region:  &orderDB[0].Region,
				Email:   &orderDB[0].CustomerEmail,
			},
			Payment: dto.PaymentInfo{
				Transaction:  &orderDB[0].Transaction,
				RequestID:    &orderDB[0].RequestID,
				Currency:     &orderDB[0].CurrencyName,
				Provider:     &orderDB[0].Provider,
				Amount:       &orderDB[0].Amount,
				PaymentDt:    &orderDB[0].PaymentDT,
				Bank:         &orderDB[0].BankName,
				DeliveryCost: &orderDB[0].DeliveryCost,
				GoodsTotal:   &orderDB[0].GoodsTotal,
				CustomFee:    &orderDB[0].CustomFee,
			},
			Items:           items,
			Locale:          &orderDB[0].Locale,
			InternalSign:    &orderDB[0].InternalSignature,
			CustomerID:      &orderDB[0].CustomerID,
			DeliveryService: &orderDB[0].DeliveryName,
			ShardKey:        &shardKey,
			SmID:            &orderDB[0].SmID,
			DateCreated:     &orderDB[0].DateCreated,
			OofShard:        &offShard,
		}
	}

	return orderInfo, nil
}

func (o *OrderRepo) SaveOrder(ctx context.Context, order *dto.Order) error {
	txOptions := &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	}

	tx, err := o.db.BeginTxx(ctx, txOptions)
	if err != nil {
		log.Println("Failed to begin transaction:", err)
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			log.Println("Transaction rolled back due to error:", err)
		} else {
			err = tx.Commit()
			if err != nil {
				log.Println("Failed to commit transaction:", err)
			}
		}
	}()

	query := `INSERT INTO bank (bank_name)
				VALUES ($1)
				ON CONFLICT (bank_name) DO UPDATE SET bank_name=EXCLUDED.bank_name RETURNING id`
	var bankID int
	if err = tx.GetContext(ctx, &bankID, query, order.Payment.Bank); err != nil {
		log.Println("Error inserting into bank:", err)
		return err
	}

	query = `INSERT INTO currency (currency_name) 
				VALUES ($1)
				ON CONFLICT (currency_name) DO UPDATE SET currency_name=EXCLUDED.currency_name RETURNING id`

	var currencyID int
	if err = tx.GetContext(ctx, &currencyID, query, order.Payment.Currency); err != nil {
		log.Println("Error inserting into currency:", err)
		return err
	}

	query = `INSERT INTO payment (transaction, request_id, provider, amount,  payment_dt, delivery_cost, goods_total, custom_fee, bank_id, currency_id) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
				ON CONFLICT (transaction) DO UPDATE SET transaction = EXCLUDED.transaction RETURNING id `

	var paymentID int
	if err := tx.GetContext(ctx, &paymentID, query, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDt, order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee, bankID, currencyID); err != nil {
		log.Println("Error inserting into payment:", err)
		return err
	}

	query = `INSERT INTO delivery_service (delivery_name) 
				VALUES ($1)
				ON CONFLICT (delivery_name) DO UPDATE SET delivery_name = EXCLUDED.delivery_name RETURNING id  `
	var deliveryServiceID int
	if err := tx.GetContext(ctx, &deliveryServiceID, query, order.DeliveryService); err != nil {
		log.Println("Error inserting into delivery_service:", err)
		return err
	}

	query = `INSERT INTO delivery (zip, city, address, region, delivery_service_id) 
				VALUES ($1, $2, $3, $4, $5)
				 ON CONFLICT (zip, city, address, region, delivery_service_id) DO UPDATE SET zip = EXCLUDED.zip RETURNING id`

	var deliveryID int
	if err := tx.GetContext(ctx, &deliveryID, query, order.Delivery.Zip, order.Delivery.City, order.Delivery.Address, order.Delivery.Region, deliveryServiceID); err != nil {
		log.Println("Error inserting into delivery:", err)
		return err
	}

	query = `INSERT INTO customer (name, phone, email, external_customer_id) 
				VALUES ($1, $2, $3, $4)
				ON CONFLICT (external_customer_id) DO UPDATE SET external_customer_id = EXCLUDED.external_customer_id RETURNING id`

	var customerID int
	if err := tx.GetContext(ctx, &customerID, query, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Email, order.CustomerID); err != nil {
		log.Println("Error inserting into customer:", err)
		return err
	}

	var (
		chrtID     *int
		status     *int
		totalPrice *int
		price      *int
	)

	if len(order.Items) > 0 {
		chrtID = order.Items[0].ChrtID
		status = order.Items[0].Status
		price = order.Items[0].Price
		totalPrice = order.Items[0].TotalPrice
	}

	query = `INSERT INTO status (name_status) 
					VALUES ($1)
					ON CONFLICT (name_status) DO UPDATE SET name_status = EXCLUDED.name_status RETURNING id`
	var statusID int
	if err := tx.GetContext(ctx, &statusID, query, status); err != nil {
		log.Println("Error inserting into status:", err)
		return err
	}

	query = `INSERT INTO orders (order_uid, track_number, chrt_id, entry, locale, internal_signature, shardkey, sm_id, date_created, off_shard, payment_id, status_id, delivery_id, customer_id)
				VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
				ON CONFLICT (order_uid) DO UPDATE SET order_uid  = EXCLUDED.order_uid RETURNING id`

	var orderID int
	if err := tx.GetContext(ctx, &orderID, query, order.OrderUID, order.TrackNumber, chrtID, order.Entry, order.Locale, order.InternalSign, order.ShardKey, order.SmID, order.DateCreated, order.OofShard, paymentID, statusID, deliveryID, customerID); err != nil {
		log.Println("Error inserting into order:", err)
		return err
	}

	for i := range order.Items {
		query = `INSERT INTO items (price, rid, name, size, nm_id, brand) 
					VALUES ($1, $2, $3, $4, $5, $6)
					ON CONFLICT (rid) DO UPDATE SET rid  = EXCLUDED.rid RETURNING id`
		var itemID int
		if err := tx.GetContext(ctx, &itemID, query, order.Items[i].Price, order.Items[i].Rid, order.Items[i].Name, order.Items[i].Size, order.Items[i].NmID, order.Items[i].Brand); err != nil {
			log.Println("Error inserting into items:", err)
			return err
		}

		query = `INSERT INTO order_items (sale, total_price, item_id, order_id)
				VALUES($1, $2, $3, $4) RETURNING id`
		var orderItemsID int
		if err := tx.GetContext(ctx, &orderItemsID, query, totalPrice, price, itemID, orderID); err != nil {
			log.Println("Error inserting into order_items:", err)
			return err
		}
	}

	return nil
}
