# Welcome to GoPool

GoPool is a very simple and easy to use pool of Goroutines.

## Usage

```golang
// Create a pool with 10 Goroutines:
pool := New(ctx, 10)

// Do any amount of work:
for range rand.Int() %10 {
    p.Go(func() {
    	time.Sleep(100*time.Millisecond)
    })
}

// Wait for the work to finish:
err := p.Wait()

// Do it again:
for range rand.Int() %10 {
    p.Go(func() {
    	time.Sleep(100*time.Millisecond)
    })
}

// Close the pool to stop all Goroutines:
err := p.Close()
```
