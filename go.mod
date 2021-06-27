module github.com/sentinel-official/cli-client

go 1.16

require (
	github.com/alessio/shellescape v1.4.1
	github.com/cosmos/cosmos-sdk v0.42.5
	github.com/go-kit/kit v0.10.0
	github.com/olekukonko/tablewriter v0.0.5
	github.com/pkg/errors v0.9.1
	github.com/sentinel-official/hub v0.7.0
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.8.0
	golang.org/x/crypto v0.0.0-20210616213533-5ff15b29337e
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e
)

replace (
	github.com/99designs/keyring => github.com/99designs/keyring v1.1.7-0.20210324095724-d9b6b92e219f
	github.com/cosmos/cosmos-sdk => github.com/sentinel-official/cosmos-sdk v0.42.6-sentinel
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
	google.golang.org/grpc => google.golang.org/grpc v1.33.2
)
