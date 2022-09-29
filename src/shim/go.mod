module github.com/confidential-containers/enclave-cc/src/shim

go 1.16

require (
	github.com/BurntSushi/toml v1.2.0
	github.com/containerd/cgroups v1.0.4
	github.com/containerd/containerd v1.6.1
	github.com/containerd/continuity v0.3.0
	github.com/containerd/go-runc v1.0.0
	github.com/containerd/ttrpc v1.1.0
	github.com/containerd/typeurl v1.0.2
	github.com/gogo/protobuf v1.3.2
	github.com/google/go-cmp v0.5.7 // indirect
	github.com/opencontainers/image-spec v1.0.2 // indirect
	github.com/opencontainers/runtime-spec v1.0.3-0.20210326190908-1c3f411f0417
	github.com/prashantv/gostub v1.1.0
	github.com/sirupsen/logrus v1.9.0
	github.com/stretchr/testify v1.8.0
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd // indirect
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8
	google.golang.org/grpc v1.49.0
	gotest.tools/v3 v3.1.0 // indirect
)

replace (
	github.com/containerd/containerd => github.com/confidential-containers/containerd v1.6.0-beta.0.0.20220303142103-c8f5e4509dcc
	google.golang.org/genproto => google.golang.org/genproto v0.0.0-20180817151627-c66870c02cf8
)
