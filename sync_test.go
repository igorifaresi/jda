//TODO

package jda

import (
  "testing"
  "runtime"
  "time"
)

func TestSemaphore(t *testing.T) {
	semaphore := NewSemaphore()
	actives := 0
	routines := 0
	
	t.Log("Stress test: Lock/Unlock, 10 goroutines, 250ms between each call, 1s goroutine duration")

	for i := 0; i < 10; i = i + 1 {
		go func() {
			semaphore.Lock()
			actives = actives+1
			time.Sleep(time.Second)
			t.Log(actives)
			actives = actives-1
			routines = routines + 1
			semaphore.Unlock()
		}()
		time.Sleep(time.Millisecond*250)
	}

	for routines < 10 {
		runtime.Gosched()
	}

	t.Logf("OK -> %d", routines)

	t.Log("Stress test: Lock/Unlock, 100 goroutines, 50ms between each call, 180ms goroutine duration")

	routines = 0

	for i := 0; i < 100; i = i + 1 {
		go func(){
			semaphore.Lock()
			time.Sleep(time.Millisecond*180)
			routines = routines + 1
			semaphore.Unlock()
		}()
		time.Sleep(time.Millisecond*50)
	}

	for routines < 100 {
		runtime.Gosched()
	}

	t.Logf("OK -> %d", routines)
}

func TestTunnel(t *testing.T) {
	tunnel := NewTunnel(2)
	semaphore := NewSemaphore()
	routines := 0

	for i := 0; i < 100; i = i + 1 {
		go func(){
			tunnel.Pass()
			runtime.Gosched()
			t.Log(tunnel.ActualWidth)
			time.Sleep(time.Millisecond*150)
			semaphore.Lock()
			routines = routines + 1
			semaphore.Unlock()
			tunnel.GoOut()
		}()
		time.Sleep(time.Millisecond*50)
	}

	for routines < 100 {
		runtime.Gosched()
	}
}