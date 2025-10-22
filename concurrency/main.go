package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"
)

func fanIn[T any](ctx context.Context, wg *sync.WaitGroup, channels ...<-chan T) <-chan T {
	fannedInStream := make(chan T)

	transfer := func(c <-chan T) {
		defer wg.Done()
		for i := range c {
			select {
			case <-ctx.Done():
				return
			case fannedInStream <- i:
			}
		}
	}

	for _, c := range channels {
		wg.Add(1)
		go transfer(c)
	}

	return fannedInStream
}

func repeatFunc[T any](ctx context.Context, fn func() T) <-chan T {
	stream := make(chan T)
	go func() {
		defer close(stream)
		for {
			select {
			case <-ctx.Done():
				return
			case stream <- fn():
			}
		}
	}()

	return stream
}

func take[T any](ctx context.Context, stream <-chan T, n int) <-chan T {
	taken := make(chan T)

	go func() {
		defer close(taken)
		for i := 0; i < n; i++ {
			select {
			case <-ctx.Done():
				return
			case taken <- <-stream:
			}
		}
	}()

	return taken
}

func primeFinder(ctx context.Context, randIntStream <-chan int) <-chan int {
	isPrime := func(randomInt int) bool {
		for i := randomInt - 1; i > 1; i-- {
			if randomInt%i == 0 {
				return false
			}
		}
		return true
	}

	primes := make(chan int)
	go func() {
		defer close(primes)
		for {
			select {
			case <-ctx.Done():
				return
			case randomInt := <-randIntStream:
				if isPrime(randomInt) {
					primes <- randomInt
				}
			}
		}
	}()

	return primes
}

func main() {
	start := time.Now()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// randIntStream continuously produces random integers from randomNumberFetcher
	randomNumberFetcher := func() int { return rand.Intn(500000000) }
	randIntStream := repeatFunc(ctx, randomNumberFetcher)

	// fan out: launch multiple primeFinder goroutines to process randIntStream in parallel
	CPUCount := runtime.NumCPU()
	primeFinderChannels := make([]<-chan int, CPUCount)
	for i := 0; i < CPUCount; i++ {
		primeFinderChannels[i] = primeFinder(ctx, randIntStream)
	}

	// fan in: merge all primeFinder outputs into a single channel, fannedInStream
	fannedInStream := fanIn(ctx, &wg, primeFinderChannels...)

	// print all prime numbers
	for rando := range take(ctx, fannedInStream, 10) {
		fmt.Println(rando)
	}
	fmt.Println(time.Since(start))

	// start the graceful shutdown
	select {
	case <-sigChan:
		fmt.Println("Received the shutdown signal, Shutting down gracefully...")
		cancel()
	}

	// wait until all goroutines have finished
	wg.Wait()

	fmt.Println("Application shutdown complete.")
}
