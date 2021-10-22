package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"sber_cloud/tw/cmd/definition/grpc/gateway"
	"sber_cloud/tw/cmd/definition/service"
	"sber_cloud/tw/container"
	"sber_cloud/tw/definition/config"
	"sber_cloud/tw/definition/kv"

	listenersGRPC "sber_cloud/tw/cmd/definition/grpc/listeners"
)

var (
	// Config path
	configPath string

	// DI Container.
	diContainer container.Container

	// MonitoringSubProcess
	subProcess string

	// config
	conf config.Config

	stopNotification = make(chan struct{})

	// Root command.
	rootCmd = &cobra.Command{
		Use:           "counter [command]",
		Long:          "counter project",
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			if diContainer, err = container.Instance(
				[]string{
					container.App,
					container.Request,
					container.SubRequest,
					listenersGRPC.DefGRPCPublicListenerScope,
					gateway.DefGRPCGatewayScope,
					service.DefServiceScope,
				},
				map[string]interface{}{
					"config":      configPath,
					"cli_cmd":     cmd,
					"cli_args":    args,
					"sub_process": subProcess,
				}); err != nil {
				return err
			}

			if err = diContainer.Fill(config.DefConfig, &conf); err != nil {
				return err
			}

			var w kv.ConsulWatcher
			if err = diContainer.Fill(kv.DefConsulWatcher, &w); err != nil {
				return err
			}

			//TODO вынести часть конфига в консул и включить merge

			/*			if err = w.Get(func(val []byte) error {
							return conf.MergeConfig(bytes.NewBuffer(val))
						}); err != nil {
							return err
						}*/

			// graceful stop
			go func() {
				var c = make(chan os.Signal, 1)
				signal.Notify(c,
					syscall.SIGHUP,
					syscall.SIGINT,
					syscall.SIGTERM,
				)

				<-c

				stopNotification <- struct{}{}
			}()

			return err
		},
	}
)

func Execute() (err error) {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "config.yaml", "config file")
	rootCmd.PersistentFlags().StringVarP(&subProcess, "subProcess", "s", "default", "monitoring subprocess tag")
	if err = rootCmd.Execute(); err != nil {
		return err
	}

	if diContainer != nil {
		return diContainer.Delete()
	}

	return nil
}
