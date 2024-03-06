package utils

import "time"

func Sleep(stop <-chan struct{}, dur time.Duration) {
	t := time.NewTimer(dur)
	select {
	case <-stop:
	case <-t.C:
	}
}
