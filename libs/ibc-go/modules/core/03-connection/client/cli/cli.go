package cli

import (
	"github.com/okex/exchain/libs/cosmos-sdk/client"
	"github.com/okex/exchain/libs/ibc-go/modules/core/03-connection/types"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the query commands for IBC connections
func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        types.SubModuleName,
		Short:                      "IBC connection query subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
	}

	queryCmd.AddCommand(
		//GetCmdQueryConnections(),
		//GetCmdQueryConnection(),
		//GetCmdQueryClientConnections(),
	)

	return queryCmd
}

// NewTxCmd returns a CLI command handler for all x/ibc connection transaction commands.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.SubModuleName,
		Short:                      "IBC connection transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		//NewConnectionOpenInitCmd(),
		//NewConnectionOpenTryCmd(),
		//NewConnectionOpenAckCmd(),
		//NewConnectionOpenConfirmCmd(),
	)

	return txCmd
}