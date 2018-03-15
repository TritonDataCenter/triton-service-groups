package cli

import (
	"fmt"
	"io"
	stdlog "log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"

	gops "github.com/google/gops/agent"
	"github.com/joyent/triton-service-groups/buildtime"
	"github.com/joyent/triton-service-groups/config"
	isatty "github.com/mattn/go-isatty"
	"github.com/pkg/errors"
	zerolog "github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sean-/conswriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// CLI flags
var (
	cfgFile string
)

var RootCmd = &cobra.Command{
	Use:   buildtime.PROGNAME,
	Short: buildtime.PROGNAME + `is the Triton Service Groups API`,
	Long: fmt.Sprintf(`
%s - Triton Service Groups API

Everything used to configure and run the central API server of the Triton
Service Groups service. Includes configuring default values, database and client
connections settings, environment variable overrides, binding and serving HTTP
access, and optionally enabling introspective utilities like gops(1) and pprof.

`, buildtime.PROGNAME),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Re-initialize logging with user-supplied configuration parameters
		{
			// os.Stdout isn't guaranteed to be thread-safe, wrap in a sync writer.
			// Files are guaranteed to be safe, terminals are not.
			var logWriter io.Writer
			if isatty.IsTerminal(os.Stdout.Fd()) ||
				isatty.IsCygwinTerminal(os.Stdout.Fd()) {
				logWriter = conswriter.GetTerminal()
			} else {
				logWriter = os.Stdout
			}

			agentFmt := viper.GetString(config.KeyAgentLogFormat)
			logFmt, err := config.LogLevelParse(agentFmt)
			if err != nil {
				return errors.Wrap(err, "unable to parse log format")
			}

			if logFmt == config.LogFormatAuto {
				if isatty.IsTerminal(os.Stdout.Fd()) ||
					isatty.IsCygwinTerminal(os.Stdout.Fd()) {
					logFmt = config.LogFormatHuman
				} else {
					logFmt = config.LogFormatZerolog
				}
			}

			var zlog zerolog.Logger
			switch logFmt {
			case config.LogFormatZerolog:
				zlog = zerolog.New(logWriter).With().Timestamp().Logger()
			case config.LogFormatHuman:
				w := zerolog.ConsoleWriter{
					Out:     logWriter,
					NoColor: false,
				}
				zlog = zerolog.New(w).With().Timestamp().Logger()
			default:
				return fmt.Errorf("unsupported log format: %q", logFmt)
			}

			log.Logger = zlog

			stdlog.SetFlags(0)
			stdlog.SetOutput(zlog)
		}

		// Perform input validation

		logLevel := strings.ToUpper(viper.GetString(config.KeyLogLevel))
		switch logLevel {
		case "DEBUG":
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		case "INFO":
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		case "WARN":
			zerolog.SetGlobalLevel(zerolog.WarnLevel)
		case "ERROR":
			zerolog.SetGlobalLevel(zerolog.ErrorLevel)
		case "FATAL":
			zerolog.SetGlobalLevel(zerolog.FatalLevel)
		default:
			// FIXME(seanc@): move the supported log levels into a global constant
			return fmt.Errorf("unsupported error level: %q (supported levels: %s)", logLevel,
				strings.Join([]string{
					"DEBUG",
					"INFO",
					"WARN",
					"ERROR",
					"FATAL",
				}, " "))
		}

		go func() {
			if !viper.GetBool(config.KeyGoogleAgentEnable) {
				log.Debug().Msg("disabled gops(1) agent")
				return
			}

			log.Info().Msg("enabled gops(1) agent")

			var (
				bind = viper.GetString(config.KeyGoogleAgentBind)
				port = viper.GetInt(config.KeyGoogleAgentPort)
				addr = fmt.Sprintf("%s:%d", bind, port)
				opts = gops.Options{
					Addr: addr,
				}
			)

			log.Debug().
				Str("gops-bind", bind).
				Int("gops-port", port).
				Msgf("starting gops(1) agent at %q", addr)

			if err := gops.Listen(opts); err != nil {
				log.Fatal().
					Err(err).
					Msg("failed to start gops(1) agent thread")
				return
			}
		}()

		go func() {
			if !viper.GetBool(config.KeyPProfEnable) {
				log.Debug().Msg("pprof endpoint disabled by request")
				return
			}

			log.Info().Msg("enabled pprof endpoint")

			var (
				bind = viper.GetString(config.KeyPProfBind)
				port = viper.GetInt(config.KeyPProfPort)
				addr = fmt.Sprintf("%s:%d", bind, port)
			)

			log.Debug().
				Str("pprof-bind", bind).
				Int("pprof-port", port).
				Msgf("starting pprof endpoint at %q", addr)

			if err := http.ListenAndServe(addr, nil); err != nil {
				log.Fatal().
					Err(err).
					Msg("failed to start pprof listener")
				return
			}
		}()

		return nil
	},
}

func Execute() error {
	return RootCmd.Execute()
}

func init() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("TSG")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile,
		"config", buildtime.PROGNAME+".toml", "config file")

	{
		const (
			key          = config.KeyLogLevel
			longName     = "log-level"
			shortName    = "l"
			defaultValue = "INFO"
			description  = "Log level"
		)

		RootCmd.PersistentFlags().StringP(
			longName,
			shortName,
			defaultValue,
			description,
		)
		err := viper.BindPFlag(key, RootCmd.PersistentFlags().Lookup(longName))
		if err != nil {
			log.Warn().Err(err)
		}
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key         = config.KeyAgentLogFormat
			longName    = "log-format"
			shortName   = "F"
			description = `Specify the log format ("auto", "zerolog", or "human")`
		)

		defaultValue := config.LogFormatAuto.String()
		RootCmd.PersistentFlags().StringP(
			longName,
			shortName,
			defaultValue,
			description,
		)
		err := viper.BindPFlag(key, RootCmd.PersistentFlags().Lookup(longName))
		if err != nil {
			log.Warn().Err(err)
		}
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = config.KeyPProfEnable
			longName     = "enable-pprof"
			shortName    = ""
			defaultValue = true
			description  = "Enable the pprof endpoint interface"
		)

		RootCmd.PersistentFlags().BoolP(
			longName,
			shortName,
			defaultValue,
			description,
		)
		err := viper.BindPFlag(key, RootCmd.PersistentFlags().Lookup(longName))
		if err != nil {
			log.Warn().Err(err)
		}
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = config.KeyPProfPort
			longName     = "pprof-port"
			shortName    = ""
			defaultValue = 4242
			description  = "Specify the pprof port"
		)

		RootCmd.PersistentFlags().Uint16P(
			longName,
			shortName,
			defaultValue,
			description,
		)
		err := viper.BindPFlag(key, RootCmd.PersistentFlags().Lookup(longName))
		if err != nil {
			log.Warn().Err(err)
		}
		viper.SetDefault(key, defaultValue)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigName(buildtime.PROGNAME)

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		d, err := os.Getwd()
		if err != nil {
			// TODO(justinwr): Maybe output this some other way, but we should
			// have defaults...
			// stdlog.Println("unable to find the current working directory")
		} else {
			viper.AddConfigPath(d)
		}
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		stdlog.Println("Unable to read config file")
	}
}
