package operation

import (
	"context"
	"fmt"
	"time"
)

func loopUntil(ctx context.Context, interval time.Duration, maxRetries int, f func() (bool, error)) error {
	count := 0
	for {
		if stop, err := f(); err != nil {
			if count++; count > maxRetries {
				return err
			}
		} else if stop {
			break
		}
		select {
		case <-time.After(interval):
		case <-ctx.Done():
			return fmt.Errorf("execute canceled")
		}
	}
	return nil
}

func loopItemsUntil(ctx context.Context, interval time.Duration, maxRetries int, items []*executeItem, f func(*executeItem) (bool, error)) error {
	completes := make(map[int]bool)
	retries := make(map[int]int)
	for i := range items {
		completes[i] = false
	}
	for {
		for i, item := range items {
			if completes[i] {
				continue
			}
			if stop, err := f(item); err != nil {
				fmt.Printf("execute failed: %v, retrying\n", err)
				if retries[i] += 1; retries[i] > maxRetries {
					return err
				}
			} else if stop {
				completes[i] = true
			}
		}
		if allCompleted(completes) {
			return nil
		}
		select {
		case <-time.After(interval):
		case <-ctx.Done():
			err := fmt.Errorf("execute canceled")
			for i, done := range completes {
				if !done {
					items[i].err = err
				}
			}
			return err
		}
	}
}

func allCompleted(completes map[int]bool) bool {
	for _, c := range completes {
		if !c {
			return false
		}
	}
	return true
}
