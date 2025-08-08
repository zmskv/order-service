package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/zmskv/order-service/internal/domain/order/entity"
	"github.com/zmskv/order-service/internal/infrastructure/repository/postgres/dto"
	"go.uber.org/zap"
)

var ErrOrderNotFound = errors.New("order not found")

type orderRepository struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewOrderRepository(db *sqlx.DB, log *zap.Logger) *orderRepository {
	return &orderRepository{db: db, logger: log.Named("order_repo")}
}

func (r *orderRepository) Save(ctx context.Context, order entity.Order) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		r.logger.Error("failed to begin transaction", zap.Error(err), zap.String("order_uid", order.OrderUID))
		return err
	}
	defer tx.Rollback()

	orderModel := dto.ToOrderModel(order)
	orderItemModels := dto.ToItemModels(order.OrderUID, order.Items)

	_, err = tx.NamedExecContext(ctx, `
		INSERT INTO orders (
			order_uid, track_number, entry, locale, internal_signature,
			customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard,
			delivery_name, delivery_phone, delivery_zip, delivery_city, delivery_address,
			delivery_region, delivery_email,
			payment_transaction, payment_request_id, payment_currency, payment_provider,
			payment_amount, payment_dt, payment_bank, payment_delivery_cost,
			payment_goods_total, payment_custom_fee
		) VALUES (
			:order_uid, :track_number, :entry, :locale, :internal_signature,
			:customer_id, :delivery_service, :shardkey, :sm_id, :date_created, :oof_shard,
			:delivery_name, :delivery_phone, :delivery_zip, :delivery_city, :delivery_address,
			:delivery_region, :delivery_email,
			:payment_transaction, :payment_request_id, :payment_currency, :payment_provider,
			:payment_amount, :payment_dt, :payment_bank, :payment_delivery_cost,
			:payment_goods_total, :payment_custom_fee
		)
	`, orderModel)
	if err != nil {
		r.logger.Error("failed to insert order", zap.Error(err), zap.String("order_uid", order.OrderUID))
		tx.Rollback()
		return err
	}

	for _, item := range orderItemModels {
		_, err := tx.NamedExecContext(ctx, `
			INSERT INTO items (
				order_uid, chrt_id, track_number, price, rid,
				name, sale, size, total_price, nm_id, brand, status
			) VALUES (
				:order_uid, :chrt_id, :track_number, :price, :rid,
				:name, :sale, :size, :total_price, :nm_id, :brand, :status
			)
		`, item)
		if err != nil {
			r.logger.Error("failed to insert order item", zap.Error(err), zap.String("order_uid", order.OrderUID))
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		r.logger.Error("failed to commit transaction", zap.Error(err), zap.String("order_uid", order.OrderUID))
		return err
	}

	return nil
}

func (r *orderRepository) Get(ctx context.Context, orderUID string) (entity.Order, error) {
	var orderModel dto.OrderModel
	var itemsModel []dto.ItemModel

	err := r.db.GetContext(ctx, &orderModel, `
        SELECT * FROM orders WHERE order_uid = $1
    `, orderUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Info("Order not found in database", zap.String("order_uid", orderUID))
			return entity.Order{}, ErrOrderNotFound
		}
		return entity.Order{}, fmt.Errorf("get order: %w", err)
	}

	err = r.db.SelectContext(ctx, &itemsModel, `
        SELECT * FROM items WHERE order_uid = $1
    `, orderUID)
	if err != nil {
		return entity.Order{}, fmt.Errorf("get order items: %w", err)
	}

	return dto.ToOrderEntity(orderModel, itemsModel), nil
}

func (r *orderRepository) GetAll(ctx context.Context) ([]entity.Order, error) {
	var orderModels []dto.OrderModel
	var itemModels []dto.ItemModel

	err := r.db.SelectContext(ctx, &orderModels, `SELECT * FROM orders`)
	if err != nil {
		return nil, fmt.Errorf("get orders: %w", err)
	}

	err = r.db.SelectContext(ctx, &itemModels, `SELECT * FROM items`)
	if err != nil {
		return nil, fmt.Errorf("get items: %w", err)
	}

	itemsMap := make(map[string][]dto.ItemModel)
	for _, item := range itemModels {
		itemsMap[item.OrderUID] = append(itemsMap[item.OrderUID], item)
	}

	var orders []entity.Order
	for _, orderModel := range orderModels {
		items := itemsMap[orderModel.OrderUID]
		order := dto.ToOrderEntity(orderModel, items)
		orders = append(orders, order)
	}

	return orders, nil
}
