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
	github.com/gorilla/mux v1.7.3 // indirect
	github.com/imdario/mergo v0.3.7 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/opencontainers/runc v0.1.1 // indirect
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v0.0.6-0.20191202130430-b04b5bfc50cb
	github.com/spf13/viper v1.4.0
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	gotest.tools v2.2.0+incompatible // indirect
	k8s.io/api v0.20.0-alpha.2
	k8s.io/apimachinery v0.20.0-alpha.2
	k8s.io/client-go v0.20.0-alpha.2
)

replace (
	github.com/Sirupsen/logrus v1.0.5 => github.com/sirupsen/logrus v1.0.5
	github.com/Sirupsen/logrus v1.3.0 => github.com/Sirupsen/logrus v1.0.6
	github.com/Sirupsen/logrus v1.4.0 => github.com/sirupsen/logrus v1.0.6
	github.com/docker/docker v1.13.1 => github.com/docker/engine v0.0.0-20190717160951-456712c5b8d9
)
