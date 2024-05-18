-- +goose Up
-- +goose StatementBegin
INSERT INTO delivery_service (delivery_name)
VALUES ('meest');

INSERT INTO delivery (zip, city, address, region, delivery_service_id)
VALUES ('2639809', 'Kiryat Mozkin', 'Ploshad Mira 15', 'Kraiot', 1);

INSERT INTO bank (bank_name)
VALUES ('alpha');

INSERT INTO currency (currency_name)
VALUES ('USD');

INSERT INTO payment (transaction, request_id, provider, amount, payment_dt, delivery_cost, goods_total, custom_fee, bank_id, currency_id)
VALUES ('b563feb7b2b84b6test', '', 'wbpay', 1817, to_timestamp(1637907727), 1500, 317, 0, 1, 1);

INSERT INTO items (price, rid, name, size, nm_id, brand)
VALUES (453, 'ab4219087a764ae0btest', 'Mascaras', 0, 2389212, 'Vivienne Sabo');

INSERT INTO items (price, rid, name, size, nm_id, brand)
VALUES (500, 'cd4219087a764ae0btest', 'Lipstick', 0, 19438381, 'Vivienne Sabo');

INSERT INTO status(name_status)
VALUES (202);

INSERT INTO customer(name, phone, email)
VALUES ('Test Testov', '+9720000000', 'test@gmail.com');

INSERT INTO orders (order_uid, track_number, chrt_id, entry, locale, internal_signature, shardkey, sm_id, date_created, off_shard, payment_id, status_id, delivery_id, customer_id)
VALUES ('b563feb7b2b84b6test', 'WBILMTESTTRACK', 9934930, 'WBIL', 'en', '', '9', 99, to_timestamp(1637907739), '1', 1, 1, 1, 1);

INSERT INTO order_items (sale, total_price, item_id, order_id)
VALUES (30, 317, 2, 1);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
