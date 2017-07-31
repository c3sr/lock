package lock

type locker struct {
	name     string
	provider Locker
}

func New(providerName string) (Locker, error) {
	provider, err := FromName(providerName)
	if err != nil {
		return nil, err
	}
	if err := provider.Init(); err != nil {
		return nil, err
	}
	return &locker{name: providerName, provider: provider}, nil
}

func (locker *locker) Init() error {
	return nil
}

func (locker *locker) Lock(id string) error {
	return locker.provider.Lock(id)
}

func (locker *locker) Unlock(id string) error {
	return locker.provider.Unlock(id)
}

func (locker *locker) Name() string {
	return locker.name
}
