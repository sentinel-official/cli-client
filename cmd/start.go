package cmd

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/spf13/cobra"

	"github.com/sentinel-official/cli-client/context"
	restmiddlewares "github.com/sentinel-official/cli-client/rest/middlewares"
	restmodules "github.com/sentinel-official/cli-client/rest/modules"
	clitypes "github.com/sentinel-official/cli-client/types"
)

func StartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the management server",
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			home, err := cmd.Flags().GetString(clitypes.FlagHome)
			if err != nil {
				return err
			}

			if err := os.MkdirAll(home, os.ModePerm); err != nil {
				return err
			}

			tty, err := cmd.Flags().GetBool(clitypes.FlagTTY)
			if err != nil {
				return err
			}

			if !tty {
				os.Stderr, os.Stdin, os.Stdout = nil, nil, nil
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			home, err := cmd.Flags().GetString(clitypes.FlagHome)
			if err != nil {
				return err
			}

			listen, err := cmd.Flags().GetString(clitypes.FlagListen)
			if err != nil {
				return err
			}

			withKeyring, err := cmd.Flags().GetBool(clitypes.FlagWithKeyring)
			if err != nil {
				return err
			}

			withService, err := cmd.Flags().GetBool(clitypes.FlagWithService)
			if err != nil {
				return err
			}

			var (
				ctx = context.NewServiceContext().
					WithHome(home)
				muxRouter    = mux.NewRouter()
				prefixRouter = muxRouter.
						PathPrefix(clitypes.APIPathPrefix).
						Subrouter()
			)

			muxRouter.Use(restmiddlewares.Log)
			prefixRouter.Use(restmiddlewares.AddHeaders)

			if withKeyring {
				restmodules.RegisterKeyring(prefixRouter, &ctx)
			}
			if withService {
				restmodules.RegisterService(prefixRouter, &ctx)
			}

			if err := os.WriteFile(
				filepath.Join(home, "url.txt"),
				[]byte("http"+"://"+listen+clitypes.APIPathPrefix),
				os.ModePerm,
			); err != nil {
				return err
			}

			router := cors.New(
				cors.Options{
					AllowedOrigins: []string{"*"},
					AllowedMethods: []string{http.MethodPost},
					AllowedHeaders: []string{"Content-Type"},
				},
			).Handler(muxRouter)

			log.Printf("Listening on %s", listen)
			return http.ListenAndServe(listen, router)
		},
	}

	cmd.Flags().Bool(clitypes.FlagWithKeyring, false, "include the endpoints of keyring module")
	cmd.Flags().Bool(clitypes.FlagWithService, false, "include the endpoints of service module")
	cmd.Flags().Bool(clitypes.FlagTTY, false, "enable the standard error, input and output")
	cmd.Flags().String(clitypes.FlagListen, clitypes.Listen, "listen address of the server")
	cmd.Flags().String(clitypes.FlagHome, clitypes.Home, "home directory of the server")

	return cmd
}
