package system

import (
	"sync"

	"github.com/pkg/errors"
	"github.com/rai-project/lock"
	"github.com/rai-project/systemlock"
)

type locker struct {
	locks map[string]*systemlock.FileMutex
	mutex *sync.RWMutex
}

func New() lock.Locker {
	locker := &locker{}
	locker.Init()
	return locker
}

func (locker *Locker) Init() error {
	locker.locks = make(map[string]*systemlock.FileMutex)
	locker.mutex = new(sync.RWMutex)
	return nil
}

func (locker *locker) Name() string {
	return "system"
}

func (locker *locker) Lock(id string) error {
	flock, err := systemlock.New(id)
	if err != nil {
		return err
	}

	flock.Lock()

	locker.mutex.Lock()
	defer locker.mutex.Unlock()

	locker.locks[id] = flock

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
	lock.Register(&locker{})
}
