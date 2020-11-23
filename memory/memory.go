package memory

// import (
// 	"sync"

// 	"github.com/pkg/errors"
// 	"github.com/c3sr/lock"
// )

// type Locker struct {
// 	locks map[string]*sync.Mutex
// 	mutex *sync.RWMutex
// }

// func New() lock.Locker {
// 	locker := new(Locker)
// 	locker.Init()
// 	return locker
// }
// func (locker *Locker) Init() error {
// 	*locker = Locker{
// 		locks: make(map[string]*sync.Mutex),
// 		mutex: new(sync.RWMutex),
// 	}
// }

// func (locker *locker) Name() string {
// 	return "memory"
// }

// func (locker *locker) Lock(id string) error {
// 	mutex := &sync.Mutex{}
// 	mutex.Lock()

// 	locker.mutex.Lock()
// 	defer locker.mutex.Unlock()

// 	locker.locks[id] = mutex

// 	return nil
// }

// func (locker *locker) Unlock(id string) error {
// 	locker.mutex.Lock()
// 	defer locker.mutex.Unlock()

// 	// Complain if no lock has been found. This can only happen if LockUpload
// 	// has not been invoked before or UnlockUpload multiple times.
// 	lock, ok := locker.locks[id]
// 	if !ok {
// 		return errors.Errorf("no lock has been found %v", id)
// 	}

// 	defer delete(locker.locks, id)

// 	return lock.Unlock()
// }

// func init() {
// 	lock.Register(&locker{})
// }

// func init() {
// 	lock.Register(&locker{})
// }
