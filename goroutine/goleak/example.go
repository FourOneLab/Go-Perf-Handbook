package goleak

import "time"

func leak() error {
	go func() {
		time.Sleep(time.Minute)
	}()
	return nil
}
