package lucy

import (
	"context"
	"sync"

	"golang.org/x/sync/semaphore"
	"golang.org/x/xerrors"
)

// Worker handles scheduled requests.
type Worker struct {
	workerQueue                WorkerQueue
	requestRestrictionStrategy RequestRestrictionStrategy
	requestSemaphore           RequestSemaphore
	httpClient                 HTTPClient
	logger                     Logger
	maxRequestNum              int64
	requestMiddlewares         []func(request *Request)
	responseMidlewares         []func(response *Response)
	spider                     func(response *Response) ([]*Request, error)
}

func newWorker(
	workerQueue WorkerQueue,
	requestRestrictionStrategy RequestRestrictionStrategy,
	requestSemaphore RequestSemaphore,
	httpClient HTTPClient,
	logger Logger,
	maxRequestNum int64,
	requestMiddlewares []func(request *Request),
	responseMidlewares []func(response *Response),
	spider func(response *Response) ([]*Request, error),
) *Worker {
	return &Worker{
		workerQueue:                workerQueue,
		requestRestrictionStrategy: requestRestrictionStrategy,
		requestSemaphore:           requestSemaphore,
		httpClient:                 httpClient,
		logger:                     logger,
		maxRequestNum:              maxRequestNum,
		requestMiddlewares:         requestMiddlewares,
		responseMidlewares:         responseMidlewares,
		spider:                     spider,
	}
}

// Start kicks worker off.
func (*Worker) Start(ctx context.Context) {

}

func (w *Worker) subscribe(ctx context.Context) (<-chan *Request, error) {
	output := make(chan *Request)
	requestChan, err := w.workerQueue.SubscribeRequests(ctx)
	if err != nil {
		return nil, xerrors.Errorf("fail to subscribe requests: %w",
			err)
	}
	go func() {
		defer close(output)

		for request := range requestChan {
			output <- request
		}
	}()
	return output, nil
}

func (w *Worker) doRequest(requestChan <-chan *Request) (<-chan *Response, error) {
	responseChan := make(chan *Response)

	// TODO: add goroutine supervisor
	go func() {
		defer close(responseChan)

		// max request per worker
		workerRequestSemaphore := semaphore.NewWeighted(w.maxRequestNum)

		requestWaitGroup := sync.WaitGroup{}

		for request := range requestChan {

			ctx := context.Background()
			err := workerRequestSemaphore.Acquire(ctx, 1)
			if err != nil {
				w.logger.Errorf("fail to acquire workerRequestSemaphore")
				continue
			}

			handleRequest := func(request *Request) {
				defer workerRequestSemaphore.Release(1)
				defer requestWaitGroup.Done()

				// request restriction
				if w.requestRestrictionStrategy.CheckRestriction() {
					resource, err := w.requestRestrictionStrategy.Resource(request)
					if err != nil {
						w.logger.Warnf("fail to get resource name for semaphore.: %v",
							err)
						return
					}
					err = w.requestSemaphore.Acquire(ctx, resource)
					if err != nil {
						w.logger.Infof("retry %s because worker failed to acquire resource.",
							request.URL)
						w.requestRestrictionStrategy.ChangePriorityWhenRestricted(request)
						err = w.workerQueue.RetryRequest(request)
						if err != nil {
							w.logger.Errorf("fail to retry %s. this request is lost.")
						}
						return
					}
					defer w.requestSemaphore.Release(resource)
				}

				// apply requestMiddlewares
				for _, middlewareFunc := range w.requestMiddlewares {
					middlewareFunc(request)
					if request == nil {
						// discard request if nil
						return
					}
				}

				// send request
				httpRequest, err := request.HTTPRequest()
				if err != nil {
					w.logger.Warnf("fail to construct http.Request. %v: %v",
						request, err)
					return
				}
				httpResponse, err := w.httpClient.Do(httpRequest)
				if err != nil {
					w.logger.Warnf("fail to get http.Response of http.Request(%v): %v",
						request, err)
					return
				}
				response, err := NewResponseFromHTTPResponse(httpResponse)
				if err != nil {
					w.logger.Warnf("fail to construct Response of http.Response(%v): %v",
						httpResponse, err)
					return
				}

				// apply responseMiddlewares
				for _, middlewareFunc := range w.responseMidlewares {
					middlewareFunc(response)
					if response == nil {
						// discard response if nil
						return
					}
				}

				responseChan <- response
			}

			requestWaitGroup.Add(1)
			go handleRequest(request)
		}

		requestWaitGroup.Wait()
	}()

	return responseChan, nil
}

func (w *Worker) applySpider(responseChan <-chan *Response) (<-chan *Request, error) {
	requestChan := make(chan *Request)

	go func() {
		defer close(requestChan)
		for response := range responseChan {
			nextRequests, err := w.spider(response)
			if err != nil {
				w.logger.Infof("spider error: %v", err)
				continue
			}
			for _, request := range nextRequests {
				requestChan <- request
			}
		}
	}()

	return requestChan, nil
}

func (w *Worker) publishRequest(requestChan <-chan *Request) error {
	for request := range requestChan {
		err := w.workerQueue.PublishRequest(request)
		if err != nil {
			w.logger.Errorf("fail to publish request: %s", request.URL)
		}
	}
	return nil
}

// WorkerBuilder is the builder of Worker.
type WorkerBuilder struct {
}
