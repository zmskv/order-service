package dto

import (
	"time"

	"github.com/zmskv/order-service/internal/domain/order/entity"
)

func ToOrderEntity(d Order) (entity.Order, error) {
	dateCreated, err := time.Parse(time.RFC3339, d.DateCreated)
	if err != nil {
		return entity.Order{}, err
	}
	items := make([]entity.Item, len(d.Items))
	for i, it := range d.Items {
		items[i] = entity.Item{
			ChrtID:      it.ChrtID,
			TrackNumber: it.TrackNumber,
			Price:       it.Price,
			RID:         it.RID,
			Name:        it.Name,
			Sale:        it.Sale,
			Size:        it.Size,
			TotalPrice:  it.TotalPrice,
			NmID:        it.NmID,
			Brand:       it.Brand,
			Status:      it.Status,
		}
	}

	return entity.Order{
		OrderUID:          d.OrderUID,
		TrackNumber:       d.TrackNumber,
		Entry:             d.Entry,
		Locale:            d.Locale,
		InternalSignature: d.InternalSignature,
		CustomerID:        d.CustomerID,
		DeliveryService:   d.DeliveryService,
		ShardKey:          d.ShardKey,
		SmID:              d.SmID,
		DateCreated:       dateCreated,
		OofShard:          d.OofShard,
		Delivery: entity.Delivery{
			Name:    d.Delivery.Name,
			Phone:   d.Delivery.Phone,
			Zip:     d.Delivery.Zip,
			City:    d.Delivery.City,
			Address: d.Delivery.Address,
			Region:  d.Delivery.Region,
			Email:   d.Delivery.Email,
		},
		Payment: entity.Payment{
			Transaction:  d.Payment.Transaction,
			RequestID:    d.Payment.RequestID,
			Currency:     d.Payment.Currency,
			Provider:     d.Payment.Provider,
			Amount:       d.Payment.Amount,
			PaymentDT:    d.Payment.PaymentDT,
			Bank:         d.Payment.Bank,
			DeliveryCost: d.Payment.DeliveryCost,
			GoodsTotal:   d.Payment.GoodsTotal,
			CustomFee:    d.Payment.CustomFee,
		},
		Items: items,
	}, nil
}
