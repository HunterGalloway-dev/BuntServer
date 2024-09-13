package scratchpads

import (
	"fmt"
	"sync"
	"time"
)

const rateLimit = time.Minute / 10

type Client interface {
	Call(*Payload, *sync.WaitGroup)
}

type MockClient struct {
}

func (mock *MockClient) Call(payload *Payload, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println(payload.msg)
}

type Payload struct {
	msg string
}

func RateLimitCall(client Client, payloads []*Payload) {
	throttle := time.Tick(rateLimit)
	var wg sync.WaitGroup
	wg.Add(len(payloads))

	for _, payload := range payloads {
		<-throttle

		go client.Call(payload, &wg)
	}

	wg.Wait()
}

func RateLimitExmp() {
	var client Client
	client = new(MockClient)

	var payloads []*Payload

	for i := 0; i < 10; i++ {
		mockPayload := Payload{
			msg: fmt.Sprintf("Hello from payload #%v", i),
		}

		payloads = append(payloads, &mockPayload)
	}

	RateLimitCall(client, payloads)
}
