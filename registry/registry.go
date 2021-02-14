package registry

import (
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/c3sr/libkv/store"
	"github.com/c3sr/lock"
	"github.com/c3sr/registry"
)

type Locker struct {
	// Client used to connect to the registry server
	Client registry.Store

	// ConnectionName is an optional field which may contain a human-readable
	// description for the connection. It is only used for composing error
	// messages and can be used to match them to a specific registry instance.
	ConnectionName string

	// locks is used for storing registry.Lock structs before they are unlocked.
	// If you want to release a lock, you need the same registry.Lock instance
	// and therefore we need to save them temporarily.
	locks map[string]store.Locker
	mutex *sync.RWMutex
}

// New constructs a new locker using the provided client.
func New(client registry.Store) *Locker {
	locker := new(Locker)
	locker.Client = client
	locker.Init()
	return locker
}
func (locker *Locker) Init() error {
	locker.locks = make(map[string]store.Locker)
	locker.mutex = new(sync.RWMutex)
	return nil
}

func (locker *Locker) Name() string {
	return "registry"
}

// LockUpload tries to obtain the exclusive lock.
func (locker *Locker) Lock(id string) error {
	lock, err := locker.Client.NewLock(id+"/"+".lock", &store.LockOptions{
		Value: []byte(".lock"),
		TTL:   20 * time.Second,
	})
	if err != nil {
		return err
	}

	ch, err := lock.Lock(nil)
	if ch == nil {
		if err == nil || err == store.ErrCannotLock {
			return errors.Wrapf(err, "failed to acquire lock %v", id)
		}
		return err
	}

	locker.mutex.Lock()
	defer locker.mutex.Unlock()
	// Only add the lock to our list if the acquire was successful and no error appeared.
	locker.locks[id] = lock

	go func() {
		// This channel will be closed once we lost the lock. This can either happen
		// wanted (using the Unlock method) or by accident, e.g. if the connection
		// to the registry server is lost.
		<-ch

		locker.mutex.RLock()
		defer locker.mutex.RUnlock()
		// Only proceed if the lock has been lost by accident. If we cannot find it
		// in the map, it has already been gracefully removed (see UnlockUpload).
		if _, ok := locker.locks[id]; !ok {
			return
		}

		msg := "registrylocker: lock for upload '" + id + "' has been lost."
		if locker.ConnectionName != "" {
			msg += " Please ensure that the connection to '" + locker.ConnectionName + "' is stable."
		} else {
			msg += " Please ensure that the connection to registry is stable (use ConnectionName to provide a printable name)."
		}

		log.Error(msg)
	}()

	return nil
}

// UnlockUpload releases a lock. If no such lock exists, no error will be returned.
func (locker *Locker) Unlock(id string) error {
	locker.mutex.Lock()
	defer locker.mutex.Unlock()

	// Complain if no lock has been found. This can only happen if LockUpload
	// has not been invoked before or UnlockUpload multiple times.
	lock, ok := locker.locks[id]
	if !ok {
		return errors.Errorf("the lock %s is not being held", id)
	}

	defer delete(locker.locks, id)

	return lock.Unlock()
}

func init() {
	lock.Register(&Locker{})
}
