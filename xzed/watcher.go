package xzed

import "os"
import "time"

type Watcher struct {
	C    chan (string)
	Done chan (struct{})
}

// Watch will send the name ro the file by the Watcher C channel
// when the named file is updated.
// If it is deleted events will not be sent until it is recreated.
// close done to stop watching.
func Watch(name string) *Watcher {
	watcher := &Watcher{}
	watcher.C = make(chan (string))
	watcher.Done = make(chan (struct{}))
	dur := time.Second * 7
	ticker := time.NewTicker(dur)

	go func() {
		prev, err := os.Stat(name)
		if err != nil {
			prev = nil
		}
		for {
			select {
			case <-ticker.C:
				now, err := os.Stat(name)
				if err != nil {
					now = nil
				} else if prev != nil &&
					(now.ModTime().After(prev.ModTime()) ||
						now.Size() != (prev.Size())) {
					watcher.C <- name
				} else if prev == nil && now != nil {
					watcher.C <- name
				}
				prev = now
			case <-watcher.Done:
				close(watcher.C)
				return
			}
		}
	}()
	return watcher
}
