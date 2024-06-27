package utils

import "time"

const Day = time.Hour * 24

func Sleep(stop <-chan struct{}, dur time.Duration) {
	t := time.NewTimer(dur)
	select {
	case <-stop:
	case <-t.C:
	}
}
