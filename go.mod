module tropos

go 1.13

require (
	github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78 // indirect
	github.com/Microsoft/go-winio v0.4.13 // indirect
	github.com/Nvveen/Gotty v0.0.0-20120604004816-cd527374f1e5 // indirect
	github.com/containerd/continuity v0.0.0-20190426062206-aaeac12a7ffc // indirect
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v1.13.1
	github.com/docker/go-connections v0.4.0
	github.com/docker/go-units v0.4.0 // indirect
	github.com/fatih/color v1.7.0
	github.com/googleapis/gnostic v0.3.0 // indirect
	github.com/gorilla/mux v1.7.3 // indirect
	github.com/imdario/mergo v0.3.7 // indirect
	github.com/mattn/go-colorable v0.1.4 // indirect
	github.com/mattn/go-isatty v0.0.11 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/opencontainers/runc v0.1.1 // indirect
	github.com/spf13/cobra v0.0.6-0.20191202130430-b04b5bfc50cb
	github.com/spf13/viper v1.4.0
	golang.org/x/crypto v0.0.0-20190418161225-b43e412143f9
	golang.org/x/net v0.0.0-20190628185345-da137c7871d7 // indirect
	google.golang.org/grpc v1.25.1 // indirect
	gotest.tools v2.2.0+incompatible // indirect
	k8s.io/api v0.0.0-20190620084959-7cf5895f2711
	k8s.io/apimachinery v0.0.0-20190612205821-1799e75a0719
	k8s.io/client-go v0.0.0-20190620085101-78d2af792bab
)

replace (
	github.com/Sirupsen/logrus v1.0.5 => github.com/sirupsen/logrus v1.0.5
	github.com/Sirupsen/logrus v1.3.0 => github.com/Sirupsen/logrus v1.0.6
	github.com/Sirupsen/logrus v1.4.0 => github.com/sirupsen/logrus v1.0.6
	github.com/docker/docker v1.13.1 => github.com/docker/engine v0.0.0-20190717160951-456712c5b8d9
)
