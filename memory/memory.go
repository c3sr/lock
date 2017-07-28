package memory

import (
	"sync"

	"github.com/pkg/errors"
	"github.com/rai-project/lock"
)

type locker struct {
	locks map[string]*sync.Mutex
	mutex *sync.RWMutex
}

func New() lock.Locker {
	return &locker{
		locks: make(map[string]*sync.Mutex),
		mutex: new(sync.RWMutex),
	}
}

func (locker *locker) Name() string {
	return "memory"
}

func (locker *locker) Lock(id string) error {
	mutex := &sync.Mutex{}
	mutex.Lock()

	locker.mutex.Lock()
	defer locker.mutex.Unlock()

	locker.locks[id] = mutex

	return nil
}

func (locker *locker) Unlock(id string) error {
	locker.mutex.Lock()
	defer locker.mutex.Unlock()

	// Complain if no lock has been found. This can only happen if LockUpload
	// has not been invoked before or UnlockUpload multiple times.
	lock, ok := locker.locks[id]
	if !ok {
		return errors.Errorf("no lock has been found %v", id)
	}

	defer delete(locker.locks, id)

	return lock.Unlock()
}

func init() {
	lock.AddLocker(&locker{})
}

func init() {
	lock.AddLocker(&locker{})
}
