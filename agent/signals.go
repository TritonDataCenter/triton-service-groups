package agent

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
)

// handleSignals runs the signal handler thread
func (a *Agent) handleSignals() {
	a.signalCh = make(chan os.Signal, 1)
	signal.Notify(a.signalCh,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer a.Stop()
	for {
		select {
		case <-a.shutdownCtx.Done():
			log.Debug().Msg("agent: removed handler for process signals")

			return
		case sig := <-a.signalCh:
			switch sig {
			case syscall.SIGINT, syscall.SIGTERM:
				log.Debug().
					Str("signal", sig.String()).
					Msgf("agent: process received %s signal", sig)

				return
			default:
				panic(fmt.Sprintf("unsupported signal: %v", sig))
			}
		}
	}
}

func (a *Agent) stopSignalCh() {
	signal.Stop(a.signalCh)
}
