module github.com/c3sr/lock

go 1.15

replace (
	github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.5
	google.golang.org/grpc => google.golang.org/grpc v1.29.1
)

require (
	github.com/c3sr/config v1.0.1
	github.com/c3sr/libkv v1.0.0
	github.com/c3sr/logger v1.0.1
	github.com/c3sr/registry v1.0.0
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.8.1
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
)
