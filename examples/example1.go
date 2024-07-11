package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/vingarcia/go-pool"
)

func main() {
	ctx := context.Background()

	// Create a pool with 10 Goroutines:
	pool := pool.New(ctx, 10)

	// Do any amount of work:
	for range rand.Int() % 10 {
		pool.Go(func() {
			time.Sleep(100 * time.Millisecond)
		})
	}

	// Wait for the work to finish:
	err := pool.Wait()
	if err != nil {
		panic(err)
	}

	// Do it again:
	for i := range rand.Int() % 10 {
		pool.Go(func() {
			time.Sleep(100 * time.Millisecond)
			fmt.Printf("finished job %d\n", i)
		})
	}

	// Close the pool to stop all Goroutines:
	err = pool.Close()
	if err != nil {
		panic(err)
	}
}
