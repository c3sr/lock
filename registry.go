package lock

import (
	"errors"
	"strings"

	"golang.org/x/sync/syncmap"
)

var lockers syncmap.Map

func FromName(s string) (Locker, error) {
	s = strings.ToLower(s)
	val, ok := lockers.Load(s)
	if !ok {
		log.WithField("locker", s).
			Warn("cannot find locker")
		return nil, errors.New("cannot find locker")
	}
	locker, ok := val.(Locker)
	if !ok {
		log.WithField("locker", s).
			Warn("invalid locker")
		return nil, errors.New("invalid locker")
	}
	return locker, nil
}

func Register(s Locker) {
	lockers.Store(strings.ToLower(s.Name()), s)
}

func Lockers() []string {
	names := []string{}
	lockers.Range(func(key, _ interface{}) bool {
		if name, ok := key.(string); ok {
			names = append(names, name)
		}
		return true
	})
	return names
}
