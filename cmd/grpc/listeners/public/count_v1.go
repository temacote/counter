package listeners

import (
	"context"
	"log"

	counter "sber_cloud/tw/proto"
)

// CountV1
func (l CounterPublicListener) CountV1(
	ctx context.Context,
	req *counter.EmptyMessage,
) (response *counter.EmptyMessage, err error) {
	log.Println("CountV1")
	return &counter.EmptyMessage{}, nil
}
