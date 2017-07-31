package lock

type Locker interface {
	Name() string
	Init() error
	Lock(id string) error
	Unlock(id string) error
}
