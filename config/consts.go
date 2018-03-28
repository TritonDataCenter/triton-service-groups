package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

const (
	KeyLogLevel = "log.level"

	KeyCRDBDatabase = "crdb.database"
	KeyCRDBHost     = "crdb.host"
	KeyCRDBPort     = "crdb.port"
	KeyCRDBUser     = "crdb.user"
	KeyCRDBPassword = "crdb.password"
	KeyCRDBMode     = "crdb.mode"

	KeyAgentLogFormat = "agent.log-format"

	KeyGoogleAgentEnable = "gops.enable"
	KeyGoogleAgentBind   = "gops.bind"
	KeyGoogleAgentPort   = "gops.port"

	KeyPProfEnable = "pprof.enable"
	KeyPProfBind   = "pprof.bind"
	KeyPProfPort   = "pprof.port"

	KeyHTTPServerBind = "http.bind"
	KeyHTTPServerPort = "http.port"
	KeyHTTPServerDC   = "http.dc"

	KeyNomadURL  = "nomad.url"
	KeyNomadPort = "nomad.port"
)

const (
	// Use a log format that resembles time.RFC3339Nano but includes all trailing
	// zeros so that we get fixed-width logging.
	LogTimeFormat = "2006-01-02T15:04:05.000000000Z07:00"
)

type LogFormat uint

const (
	LogFormatAuto LogFormat = iota
	LogFormatZerolog
	LogFormatHuman
)

func (f LogFormat) String() string {
	switch f {
	case LogFormatAuto:
		return "auto"
	case LogFormatZerolog:
		return "zerolog"
	case LogFormatHuman:
		return "human"
	default:
		panic(fmt.Sprintf("unknown log format: %d", f))
	}
}

func LogLevelParse(s string) (LogFormat, error) {
	switch logFormat := strings.ToLower(viper.GetString(KeyAgentLogFormat)); logFormat {
	case "auto":
		return LogFormatAuto, nil
	case "json", "zerolog":
		return LogFormatZerolog, nil
	case "human":
		return LogFormatHuman, nil
	default:
		return LogFormatAuto, fmt.Errorf("unsupported log format: %q", logFormat)
	}
}
