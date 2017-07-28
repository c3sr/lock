package lock

type Locker interface {
	Name() string
	Lock(id string) error
	Unlock(id string) error
}
