package cmd

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/sentinel-official/cli-client/context"
	"github.com/sentinel-official/cli-client/middlewares"
	keysrest "github.com/sentinel-official/cli-client/rest/keys"
	servicerest "github.com/sentinel-official/cli-client/rest/service"
	"github.com/sentinel-official/cli-client/types"
	configtypes "github.com/sentinel-official/cli-client/types/config"
	randutils "github.com/sentinel-official/cli-client/utils/rand"
)

func StartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start REST API server",
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				home       = viper.GetString(flags.FlagHome)
				configPath = filepath.Join(home, types.ConfigFilename)
				buildPath  = filepath.Join(home, types.BuildFolderName)
			)

			v := viper.New()
			v.SetConfigFile(configPath)

			config, err := configtypes.ReadInConfig(v)
			if err != nil {
				return err
			}
			if err := config.Validate(); err != nil {
				return err
			}

			config.Token = randutils.RandomStringHex(types.TokenLength)
			if err := config.SaveToPath(configPath); err != nil {
				return err
			}

			var (
				ctx = context.NewContext().
					WithConfig(config)

				muxRouter    = mux.NewRouter()
				prefixRouter = muxRouter.PathPrefix(types.APIPathPrefix).Subrouter()
			)

			muxRouter.Use(middlewares.Log)
			muxRouter.PathPrefix("/").
				Handler(http.FileServer(http.Dir(buildPath)))

			prefixRouter.Use(middlewares.AddHeaders)
			keysrest.RegisterRoutes(prefixRouter, ctx)
			servicerest.RegisterRoutes(prefixRouter, ctx)

			router := cors.New(
				cors.Options{
					AllowedOrigins: strings.Split(config.CORS.AllowedOrigins, ","),
					AllowedMethods: []string{http.MethodPost},
					AllowedHeaders: []string{"Content-Type"},
				},
			).Handler(muxRouter)

			return http.ListenAndServe(config.ListenOn, router)
		},
	}

	return cmd
}
