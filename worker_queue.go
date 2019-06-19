package lucy

//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE -self_package=github.com/getumen/lucy

import "context"

// WorkerQueue is queue operations for scheduled requests.
type WorkerQueue interface {
	SubscribeRequests(ctx context.Context) (<-chan Request, error)
	RetryRequest(request *Request) error
}
