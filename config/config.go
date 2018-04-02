package config

import (
	"fmt"
	"strings"

	nomad "github.com/hashicorp/nomad/api"
	"github.com/jackc/pgx"
	"github.com/joyent/triton-service-groups/buildtime"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

type DBPool = pgx.ConnPoolConfig

type Config struct {
	DBPool
	Agent
	HTTPServer
	Nomad
}

type Agent struct {
	LogFormat LogFormat
}

type HTTPServer struct {
	Bind   string
	Port   uint16
	Logger zerolog.Logger
	DC     string
}

type PGXLogger struct {
	logger zerolog.Logger
}

type Nomad struct {
	Addr      string
	Port      uint16
	TLSConfig *nomad.TLSConfig
}

// Custom logging facade that implements the pgx.Logger interface in order to
// log through Zerolog
func (l *PGXLogger) Log(level pgx.LogLevel, msg string, data map[string]interface{}) {
	var zlevel zerolog.Level
	switch level {
	case pgx.LogLevelNone:
		zlevel = zerolog.NoLevel
	case pgx.LogLevelError:
		zlevel = zerolog.ErrorLevel
	case pgx.LogLevelWarn:
		zlevel = zerolog.WarnLevel
	case pgx.LogLevelInfo:
		// NOTE(justinwr): We want to force into debug output through zerolog.
		zlevel = zerolog.DebugLevel
	case pgx.LogLevelDebug:
		zlevel = zerolog.DebugLevel
	default:
		zlevel = zerolog.DebugLevel
	}

	pgxlog := l.logger.With().Fields(data).Logger()
	pgxlog.WithLevel(zlevel).Msg(msg)
}

func NewDefault() (cfg *Config, err error) {
	var pgxLogLevel int = pgx.LogLevelInfo
	switch logLevel := strings.ToUpper(viper.GetString(KeyLogLevel)); logLevel {
	case "FATAL":
		pgxLogLevel = pgx.LogLevelNone
	case "ERROR":
		pgxLogLevel = pgx.LogLevelError
	case "WARN":
		pgxLogLevel = pgx.LogLevelWarn
	case "INFO":
		// NOTE(justinwr): If the app was set for INFO than we'll want to force
		// pgx to output to Debug.
		pgxLogLevel = pgx.LogLevelInfo
	case "DEBUG":
		pgxLogLevel = pgx.LogLevelDebug
	default:
		panic(fmt.Sprintf("unsupported log level: %q", logLevel))
	}

	agentConfig := Agent{}
	{
		agentConfig.LogFormat, err = LogLevelParse(viper.GetString(KeyAgentLogFormat))
		if err != nil {
			return nil, errors.Wrap(err, "unable to parse the log format")
		}
	}

	httpServerConfig := HTTPServer{}
	{
		httpServerConfig.Logger = log.Logger.With().Str("module", "http").Logger()

		httpServerConfig.Bind = "127.0.0.1"
		if bind := viper.GetString(KeyHTTPServerBind); bind != "" {
			httpServerConfig.Bind = bind
		}

		httpServerConfig.Port = uint16(3000)
		if port := viper.GetInt(KeyHTTPServerPort); port != 0 {
			httpServerConfig.Port = uint16(port)
		}

		httpServerConfig.DC = "us-east-1"
		if dc := viper.GetString(KeyHTTPServerDC); dc != "" {
			httpServerConfig.DC = dc
		}
	}

	pgxLogger := &PGXLogger{}
	{
		pgxLogger.logger = log.Logger.With().Str("module", "pgx").Logger()
	}

	// default to commonly configured CockroachDB port
	viper.SetDefault(KeyCRDBPort, uint16(26257))

	nomadConfig := Nomad{}
	{
		nomadConfig.Addr = "127.0.0.1"
		if url := viper.GetString(KeyNomadURL); url != "" {
			nomadConfig.Addr = url
		}

		nomadConfig.Port = uint16(4242)
		if port := cast.ToUint16(viper.GetInt(KeyNomadPort)); port != 0 {
			nomadConfig.Port = port
		}
	}

	return &Config{
		DBPool: pgx.ConnPoolConfig{
			MaxConnections: 5,
			AfterConnect:   nil,
			AcquireTimeout: 0,

			ConnConfig: pgx.ConnConfig{
				Database: viper.GetString(KeyCRDBDatabase),
				User:     viper.GetString(KeyCRDBUser),
				Password: viper.GetString(KeyCRDBPassword),
				Host:     viper.GetString(KeyCRDBHost),
				Port:     cast.ToUint16(viper.GetInt(KeyCRDBPort)),
				// TLSConfig: &tls.Config{},
				Logger:   pgxLogger,
				LogLevel: pgxLogLevel,
				RuntimeParams: map[string]string{
					"application_name": buildtime.PROGNAME,
				},
			},
		},
		Agent:      agentConfig,
		HTTPServer: httpServerConfig,
		Nomad:      nomadConfig,
	}, nil
}

// IsDebug returns true when the server is configured for debug level
func IsDebug() bool {
	switch logLevel := strings.ToUpper(viper.GetString(KeyLogLevel)); logLevel {
	case "DEBUG":
		return true
	default:
		return false
	}
}
