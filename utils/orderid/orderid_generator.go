package orderid

import (
	"github.com/segmentio/ksuid"
)

func Generator() string {
	orderID := ksuid.New()
	return orderID.String()
}