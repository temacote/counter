package listeners

import (
	"context"

	counter "sber_cloud/tw/proto"
)

// CreateUserV1 метод создания нового пользователя
func (l CounterPublicListener) CountV1(
	ctx context.Context,
	req *counter.EmptyMessage,
) (response *counter.EmptyMessage, err error) {
	return
}
