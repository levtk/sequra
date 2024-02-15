package interfaces

import (
	"github.com/google/uuid"
	"time"
)

type Disburser interface {
	ProcessOrder() error
	Import() error
	Print(method models.PrintMethod) error
}

type DisburserRepo interface {
	GetOrdersByMerchantUUID(id uuid.UUID) ([]models.Order, error)
	GetOrdersByDate(date time.Time) ([]models.Order, error)
	GetMerchantByReference(merchReference string) (models.Merchant, error)
}
