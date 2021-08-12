module github.com/sentinel-official/cli-client

go 1.16

require (
	github.com/alessio/shellescape v1.4.1
	github.com/cosmos/cosmos-sdk v0.42.5
	github.com/cosmos/go-bip39 v1.0.0
	github.com/go-kit/kit v0.10.0
	github.com/gorilla/mux v1.8.0
	github.com/natefinch/atomic v1.0.1
	github.com/olekukonko/tablewriter v0.0.5
	github.com/pkg/errors v0.9.1
	github.com/rs/cors v1.7.0
	github.com/sentinel-official/hub v0.7.0
	github.com/spf13/cobra v1.1.3
	github.com/tendermint/tendermint v0.34.10
	golang.org/x/crypto v0.0.0-20210616213533-5ff15b29337e
)

replace (
	github.com/99designs/keyring => github.com/99designs/keyring v1.1.7-0.20210324095724-d9b6b92e219f
	github.com/cosmos/cosmos-sdk => github.com/sentinel-official/cosmos-sdk v0.42.6-sentinel
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
	google.golang.org/grpc => google.golang.org/grpc v1.33.2
)
