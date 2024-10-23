package cli

import (
	"github.com/nahria/timetrace/core"
	"github.com/nahria/timetrace/out"

	"github.com/spf13/cobra"
)

func stopCommand(t *core.Timetrace) *cobra.Command {
	stop := &cobra.Command{
		Use:   "stop",
		Short: "Stop tracking your time",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			if err := t.Stop(); err != nil {
				out.Err("failed to stop tracking: %s", err.Error())
				return
			}

			out.Success("Stopped tracking time")
		},
	}

	return stop
}
