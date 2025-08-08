package dto

import (
	"github.com/zmskv/order-service/internal/domain/order/entity"
)

func ToOrderModel(e entity.Order) OrderModel {
	return OrderModel{
		OrderUID:            e.OrderUID,
		TrackNumber:         e.TrackNumber,
		Entry:               e.Entry,
		Locale:              e.Locale,
		InternalSignature:   e.InternalSignature,
		CustomerID:          e.CustomerID,
		DeliveryService:     e.DeliveryService,
		ShardKey:            e.ShardKey,
		SmID:                e.SmID,
		DateCreated:         e.DateCreated,
		OofShard:            e.OofShard,
		DeliveryName:        e.Delivery.Name,
		DeliveryPhone:       e.Delivery.Phone,
		DeliveryZip:         e.Delivery.Zip,
		DeliveryCity:        e.Delivery.City,
		DeliveryAddress:     e.Delivery.Address,
		DeliveryRegion:      e.Delivery.Region,
		DeliveryEmail:       e.Delivery.Email,
		PaymentTransaction:  e.Payment.Transaction,
		PaymentRequestID:    e.Payment.RequestID,
		PaymentCurrency:     e.Payment.Currency,
		PaymentProvider:     e.Payment.Provider,
		PaymentAmount:       e.Payment.Amount,
		PaymentDT:           e.Payment.PaymentDT,
		PaymentBank:         e.Payment.Bank,
		PaymentDeliveryCost: e.Payment.DeliveryCost,
		PaymentGoodsTotal:   e.Payment.GoodsTotal,
		PaymentCustomFee:    e.Payment.CustomFee,
	}
}

func ToItemModels(orderUID string, items []entity.Item) []ItemModel {
	var result []ItemModel
	for _, i := range items {
		result = append(result, ItemModel{
			OrderUID:    orderUID,
			ChrtID:      i.ChrtID,
			TrackNumber: i.TrackNumber,
			Price:       i.Price,
			RID:         i.RID,
			Name:        i.Name,
			Sale:        i.Sale,
			Size:        i.Size,
			TotalPrice:  i.TotalPrice,
			NmID:        i.NmID,
			Brand:       i.Brand,
			Status:      i.Status,
		})
	}
	return result
}

func ToOrderEntity(order OrderModel, items []ItemModel) entity.Order {
	var itemEntities []entity.Item
	for _, item := range items {
		itemEntities = append(itemEntities, entity.Item{
			ChrtID:      item.ChrtID,
			TrackNumber: item.TrackNumber,
			Price:       item.Price,
			RID:         item.RID,
			Name:        item.Name,
			Sale:        item.Sale,
			Size:        item.Size,
			TotalPrice:  item.TotalPrice,
			NmID:        item.NmID,
			Brand:       item.Brand,
			Status:      item.Status,
		})
	}

	return entity.Order{
		OrderUID:          order.OrderUID,
		TrackNumber:       order.TrackNumber,
		Entry:             order.Entry,
		Locale:            order.Locale,
		InternalSignature: order.InternalSignature,
		CustomerID:        order.CustomerID,
		DeliveryService:   order.DeliveryService,
		ShardKey:          order.ShardKey,
		SmID:              order.SmID,
		DateCreated:       order.DateCreated,
		OofShard:          order.OofShard,
		Delivery: entity.Delivery{
			Name:    order.DeliveryName,
			Phone:   order.DeliveryPhone,
			Zip:     order.DeliveryZip,
			City:    order.DeliveryCity,
			Address: order.DeliveryAddress,
			Region:  order.DeliveryRegion,
			Email:   order.DeliveryEmail,
		},
		Payment: entity.Payment{
			Transaction:  order.PaymentTransaction,
			RequestID:    order.PaymentRequestID,
			Currency:     order.PaymentCurrency,
			Provider:     order.PaymentProvider,
			Amount:       order.PaymentAmount,
			PaymentDT:    order.PaymentDT,
			Bank:         order.PaymentBank,
			DeliveryCost: order.PaymentDeliveryCost,
			GoodsTotal:   order.PaymentGoodsTotal,
			CustomFee:    order.PaymentCustomFee,
		},
		Items: itemEntities,
	}
}
