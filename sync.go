package jda

import (
	"sync"
    "sync/atomic"
    "runtime"
)

const (
	LockLocked = iota
	LockUnlocked
)

type Semaphore struct {
	State    int32
	Freeness int32
	FreenessMutex sync.Mutex
}

func NewSemaphore() Semaphore {
	return Semaphore{
		State: LockUnlocked,
		Freeness: 0,
	}
}

func (m *Semaphore) Lock() {
	for !atomic.CompareAndSwapInt32(&m.State, LockUnlocked, LockLocked) {
		runtime.Gosched()
	}
}

func (m *Semaphore) LockAndWaitFreeness() {
	m.State = LockLocked
	for m.Freeness > 0 {
		runtime.Gosched()
	}
}

func (m *Semaphore) WaitFreeness() {
	for m.Freeness > 0 {
		runtime.Gosched()
	}
}

func (m *Semaphore) Unlock() {
	for !atomic.CompareAndSwapInt32(&m.State, LockLocked, LockUnlocked) {
		runtime.Gosched()
	}
}

func (m *Semaphore) CompareAndUnlockForOne(state bool) bool {
	if state {
		m.Lock()
		return true
	}
	return false
}

func (m *Semaphore) ContinueIfUnlocked() {
	for m.State != LockUnlocked {
		runtime.Gosched()
	}
}

func (m *Semaphore) IncrementFreeness() {
	m.FreenessMutex.Lock()
	m.Freeness = m.Freeness + 1
	m.FreenessMutex.Unlock()
}

func (m *Semaphore) DecrementFreeness() {
	m.FreenessMutex.Lock()
	m.Freeness = m.Freeness - 1
	m.FreenessMutex.Unlock()
}

type Tunnel struct {
	ActualWidth           int32
	ActualWidthSemaphore Semaphore
	MaxWidth              int32
	MaxWidthMutex         sync.Mutex
}

func NewTunnel(maxWidth int32) Tunnel {
	return Tunnel{
		MaxWidth: maxWidth,
		ActualWidthSemaphore: NewSemaphore(),
	}
}

func (t *Tunnel) Pass() {
	t.ActualWidthSemaphore.Lock()
	for t.ActualWidth >= t.MaxWidth {
		runtime.Gosched()
		t.ActualWidthSemaphore.Unlock()
		t.ActualWidthSemaphore.Lock()
	}
	t.ActualWidth = t.ActualWidth + 1
	t.ActualWidthSemaphore.Unlock()
}

func (t *Tunnel) GoOut() {
	t.ActualWidthSemaphore.Lock()
	t.ActualWidth = t.ActualWidth - 1
	t.ActualWidthSemaphore.Unlock()
}

func (t *Tunnel) SetMaxWidth(width int32) {
	t.MaxWidthMutex.Lock()
	t.MaxWidth = width
	t.MaxWidthMutex.Unlock()
}