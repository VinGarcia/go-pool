package pool

import (
	"context"
	"testing"
	"time"

	tt "github.com/vingarcia/go-pool/internal/testtools"
)

func TestPool(t *testing.T) {
	ctx := context.Background()

	t.Run("should process multiple jobs with a single Goroutine", func(t *testing.T) {
		p := New(ctx, 1)
		var r1, r2, r3 string
		p.Go(func() {
			r1 = "result1"
		})
		p.Go(func() {
			r2 = "result2"
		})
		p.Go(func() {
			r3 = "result3"
		})

		err := p.Wait()
		tt.AssertNoErr(t, err)

		tt.AssertEqual(t, r1, "result1")
		tt.AssertEqual(t, r2, "result2")
		tt.AssertEqual(t, r3, "result3")
	})

	t.Run("should process multiple jobs with a multiple Goroutines", func(t *testing.T) {
		p := New(ctx, 3)
		var r1, r2, r3 string
		p.Go(func() {
			r1 = "result1"
		})
		p.Go(func() {
			r2 = "result2"
		})
		p.Go(func() {
			r3 = "result3"
		})

		err := p.Wait()
		tt.AssertNoErr(t, err)

		tt.AssertEqual(t, r1, "result1")
		tt.AssertEqual(t, r2, "result2")
		tt.AssertEqual(t, r3, "result3")
	})

	t.Run("should be reusable after a call to Wait()", func(t *testing.T) {
		p := New(ctx, 1)
		var r1 string
		p.Go(func() {
			r1 = "result1"
		})

		err := p.Wait()
		tt.AssertNoErr(t, err)

		tt.AssertEqual(t, r1, "result1")

		var r2 string
		p.Go(func() {
			r2 = "result2"
		})
		err = p.Wait()
		tt.AssertNoErr(t, err)

		tt.AssertEqual(t, r2, "result2")
	})

	t.Run("close should wait for all goroutines to finish", func(t *testing.T) {
		waiterCh := make(chan struct{})

		p := New(ctx, 1)

		p.Go(func() {
			<-waiterCh
		})

		closed := false
		go func() {
			p.Close()
			closed = true
		}()

		// Wait for the p.Close() do be called:
		time.Sleep(100 * time.Millisecond)
		tt.AssertEqual(t, closed, false)

		close(waiterCh)

		// Wait for it to finish closing:
		time.Sleep(100 * time.Millisecond)
		tt.AssertEqual(t, closed, true)
	})
}
