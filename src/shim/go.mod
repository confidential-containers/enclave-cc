module github.com/confidential-containers/enclave-cc/src/shim

go 1.20

require (
	github.com/BurntSushi/toml v1.3.2
	github.com/containerd/cgroups v1.1.0
	github.com/containerd/containerd v1.6.1
	github.com/containerd/continuity v0.4.1
	github.com/containerd/go-runc v1.0.0
	github.com/containerd/ttrpc v1.1.2
	github.com/containerd/typeurl v1.0.2
	github.com/gogo/protobuf v1.3.2
	github.com/opencontainers/runtime-spec v1.1.0-rc.3
	github.com/prashantv/gostub v1.1.0
	github.com/sirupsen/logrus v1.9.3
	github.com/stretchr/testify v1.8.4
	golang.org/x/sys v0.10.0
	google.golang.org/grpc v1.56.1
)

require (
	github.com/Microsoft/go-winio v0.5.2 // indirect
	github.com/Microsoft/hcsshim v0.9.3 // indirect
	github.com/cilium/ebpf v0.9.1 // indirect
	github.com/containerd/console v1.0.3 // indirect
	github.com/containerd/fifo v1.0.0 // indirect
	github.com/coreos/go-systemd/v22 v22.3.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/docker/go-events v0.0.0-20190806004212-e31b211e4f1c // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/godbus/dbus/v5 v5.0.6 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/moby/sys/mountinfo v0.5.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.3-0.20211202183452-c5a74bcca799 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.opencensus.io v0.23.0 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	google.golang.org/genproto v0.0.0-20230410155749-daa745c078e1 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gotest.tools/v3 v3.1.0 // indirect
)

replace (
	github.com/containerd/containerd => github.com/confidential-containers/containerd v1.6.7-0.20220619164525-4b12e77e79dc
	google.golang.org/genproto => google.golang.org/genproto v0.0.0-20180817151627-c66870c02cf8
)
