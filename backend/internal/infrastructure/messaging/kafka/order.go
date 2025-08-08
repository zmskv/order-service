package kafka

import (
	"encoding/json"

	"github.com/zmskv/order-service/internal/infrastructure/messaging/kafka/dto"
)

func HandleOrder(data []byte) (dto.Order, error) {
	var order dto.Order
	err := json.Unmarshal(data, &order)
	if err != nil {
		return dto.Order{}, err
	}
	return order, nil
}
