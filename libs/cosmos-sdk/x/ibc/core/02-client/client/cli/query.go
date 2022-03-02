package cli
//
//import (
//	"bufio"
//	"errors"
//	"fmt"
//	"github.com/cosmos/cosmos-sdk/version"
//	"github.com/okex/exchain/libs/cosmos-sdk/client"
//	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
//	"github.com/okex/exchain/libs/cosmos-sdk/client/flags"
//	"github.com/okex/exchain/libs/cosmos-sdk/codec"
//	"github.com/okex/exchain/libs/cosmos-sdk/x/ibc/core/02-client/client/utils"
//	"github.com/okex/exchain/libs/cosmos-sdk/x/ibc/core/02-client/types"
//	host "github.com/okex/exchain/libs/cosmos-sdk/x/ibc/core/24-host"
//	"github.com/spf13/cobra"
//)
//
//const (
//	flagLatestHeight = "latest-height"
//)
//
//// GetCmdQueryClientStates defines the command to query all the light clients
//// that this chain mantains.
//func GetCmdQueryClientStates(cdc *codec.Codec) *cobra.Command {
//	cmd := &cobra.Command{
//		Use:     "states",
//		Short:   "Query all available light clients",
//		Long:    "Query all available light clients",
//		Example: fmt.Sprintf("%s query %s %s states", version.AppName, host.ModuleName, types.SubModuleName),
//		Args:    cobra.NoArgs,
//		RunE: func(cmd *cobra.Command, _ []string) error {
//			clientCtx :=client.GetContext(cmd,cdc)
//			queryClient := types.NewQueryClient(clientCtx)
//
//			pageReq, err := client.ReadPageRequest(cmd.Flags())
//			if err != nil {
//				return err
//			}
//
//			req := &types.QueryClientStatesRequest{
//				Pagination: pageReq,
//			}
//
//			res, err := queryClient.ClientStates(cmd.Context(), req)
//			if err != nil {
//				return err
//			}
//
//			return clientCtx.Response(res)
//		},
//	}
//	flags.AddQueryFlagsToCmd(cmd)
//	flags.AddPaginationFlagsToCmd(cmd, "client states")
//
//	return cmd
//}
//
//// GetCmdQueryClientState defines the command to query the state of a client with
//// a given id as defined in https://github.com/cosmos/ibc/tree/master/spec/core/ics-002-client-semantics#query
//func GetCmdQueryClientState(cdc *codec.Codec) *cobra.Command {
//	cmd := &cobra.Command{
//		Use:     "state [client-id]",
//		Short:   "Query a client state",
//		Long:    "Query stored client state",
//		Example: fmt.Sprintf("%s query %s %s state [client-id]", version.AppName, host.ModuleName, types.SubModuleName),
//		Args:    cobra.ExactArgs(1),
//		RunE: func(cmd *cobra.Command, args []string) error {
//			inBuf := bufio.NewReader(cmd.InOrStdin())
//			clientCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
//			clientID := args[0]
//			prove, _ := cmd.Flags().GetBool(flags.FlagProve)
//
//			clientStateRes, err := utils.QueryClientState(clientCtx, clientID, prove)
//			if err != nil {
//				return err
//			}
//
//			return clientCtx.Response(clientStateRes)
//		},
//	}
//
//	cmd.Flags().Bool(flags.FlagProve, true, "show proofs for the query results")
//	flags.AddQueryFlagsToCmd(cmd)
//
//	return cmd
//}
//
//// GetCmdQueryClientStatus defines the command to query the status of a client with a given id
////func GetCmdQueryClientStatus(cdc *codec.Codec) *cobra.Command {
////	cmd := &cobra.Command{
////		Use:     "status [client-id]",
////		Short:   "Query client status",
////		Long:    "Query client activity status. Any client without an 'Active' status is considered inactive",
////		Example: fmt.Sprintf("%s query %s %s status [client-id]", version.AppName, host.ModuleName, types.SubModuleName),
////		Args:    cobra.ExactArgs(1),
////		RunE: func(cmd *cobra.Command, args []string) error {
////
////			clientCtx := client.GetContext(cmd, cdc)
////
////			clientID := args[0]
////			queryClient := types.NewQueryClient(clientCtx)
////
////			req := &types.QueryClientStateRequest{
////				ClientId: clientID,
////			}
////
////			clientStatusRes, err := queryClient.ClientStatus(cmd.Context(), req)
////			if err != nil {
////				return err
////			}
////
////			return clientCtx.PrintProto(clientStatusRes)
////		},
////	}
////
////	return cmd
////}
//
//// GetCmdQueryConsensusStates defines the command to query all the consensus states from a given
//// client state.
//func GetCmdQueryConsensusStates(cdc *codec.Codec) *cobra.Command {
//	cmd := &cobra.Command{
//		Use:     "consensus-states [client-id]",
//		Short:   "Query all the consensus states of a client.",
//		Long:    "Query all the consensus states from a given client state.",
//		Example: fmt.Sprintf("%s query %s %s consensus-states [client-id]", version.AppName, host.ModuleName, types.SubModuleName),
//		Args:    cobra.ExactArgs(1),
//		RunE: func(cmd *cobra.Command, args []string) error {
//			clientCtx:=client.GetContext(cmd,cdc)
//			clientID := args[0]
//
//			queryClient := types.NewQueryClient(clientCtx)
//
//			pageReq, err := client.ReadPageRequest(cmd.Flags())
//			if err != nil {
//				return err
//			}
//
//			req := &types.QueryConsensusStatesRequest{
//				ClientId:   clientID,
//				Pagination: pageReq,
//			}
//
//			res, err := queryClient.ConsensusStates(cmd.Context(), req)
//			if err != nil {
//				return err
//			}
//
//			return clientCtx.Response(res)
//		},
//	}
//	flags.AddQueryFlagsToCmd(cmd)
//	flags.AddPaginationFlagsToCmd(cmd, "consensus states")
//
//	return cmd
//}
//
//// GetCmdQueryConsensusState defines the command to query the consensus state of
//// the chain as defined in https://github.com/cosmos/ibc/tree/master/spec/core/ics-002-client-semantics#query
//func GetCmdQueryConsensusState(cdc *codec.Codec) *cobra.Command {
//	cmd := &cobra.Command{
//		Use:   "consensus-state [client-id] [height]",
//		Short: "Query the consensus state of a client at a given height",
//		Long: `Query the consensus state for a particular light client at a given height.
//If the '--latest' flag is included, the query returns the latest consensus state, overriding the height argument.`,
//		Example: fmt.Sprintf("%s query %s %s  consensus-state [client-id] [height]", version.AppName, host.ModuleName, types.SubModuleName),
//		Args:    cobra.RangeArgs(1, 2),
//		RunE: func(cmd *cobra.Command, args []string) error {
//			clientCtx:=client.GetContext(cmd,cdc)
//			clientID := args[0]
//			queryLatestHeight, _ := cmd.Flags().GetBool(flagLatestHeight)
//			var height types.Height
//			var err error
//
//			if !queryLatestHeight {
//				if len(args) != 2 {
//					return errors.New("must include a second 'height' argument when '--latest-height' flag is not provided")
//				}
//
//				height, err = types.ParseHeight(args[1])
//				if err != nil {
//					return err
//				}
//			}
//
//			prove, _ := cmd.Flags().GetBool(flags.FlagProve)
//
//			csRes, err := utils.QueryConsensusState(clientCtx, clientID, height, prove, queryLatestHeight)
//			if err != nil {
//				return err
//			}
//
//			return clientCtx.Response(csRes)
//		},
//	}
//
//	cmd.Flags().Bool(flags.FlagProve, true, "show proofs for the query results")
//	cmd.Flags().Bool(flagLatestHeight, false, "return latest stored consensus state")
//	flags.AddQueryFlagsToCmd(cmd)
//
//	return cmd
//}
//
//// GetCmdQueryHeader defines the command to query the latest header on the chain
//func GetCmdQueryHeader(cdc *codec.Codec) *cobra.Command {
//	cmd := &cobra.Command{
//		Use:     "header",
//		Short:   "Query the latest header of the running chain",
//		Long:    "Query the latest Tendermint header of the running chain",
//		Example: fmt.Sprintf("%s query %s %s  header", version.AppName, host.ModuleName, types.SubModuleName),
//		Args:    cobra.NoArgs,
//		RunE: func(cmd *cobra.Command, _ []string) error {
//			clientCtx:=client.GetContext(cmd,cdc)
//			header, _, err := utils.QueryTendermintHeader(clientCtx)
//			if err != nil {
//				return err
//			}
//
//			return clientCtx.Response(&header)
//		},
//	}
//
//	flags.AddQueryFlagsToCmd(cmd)
//
//	return cmd
//}
//
//// GetCmdSelfConsensusState defines the command to query the self consensus state of a chain
//func GetCmdSelfConsensusState(cdc *codec.Codec) *cobra.Command {
//	cmd := &cobra.Command{
//		Use:     "self-consensus-state",
//		Short:   "Query the self consensus state for this chain",
//		Long:    "Query the self consensus state for this chain. This result may be used for verifying IBC clients representing this chain which are hosted on counterparty chains.",
//		Example: fmt.Sprintf("%s query %s %s self-consensus-state", version.AppName, host.ModuleName, types.SubModuleName),
//		Args:    cobra.NoArgs,
//		RunE: func(cmd *cobra.Command, _ []string) error {
//			clientCtx:=client.GetContext(cmd,cdc)
//			state, _, err := utils.QuerySelfConsensusState(clientCtx)
//			if err != nil {
//				return err
//			}
//
//			return clientCtx.Response(state)
//		},
//	}
//
//	flags.AddQueryFlagsToCmd(cmd)
//
//	return cmd
//}
//
//
//// GetCmdNodeConsensusState defines the command to query the latest consensus state of a node
//// The result is feed to client creation
//func GetCmdNodeConsensusState(cdc *codec.Codec) *cobra.Command {
//	cmd := &cobra.Command{
//		Use:     "node-state",
//		Short:   "Query a node consensus state",
//		Long:    "Query a node consensus state. This result is feed to the client creation transaction.",
//		Example: fmt.Sprintf("%s query %s %s node-state", version.AppName, host.ModuleName, types.SubModuleName),
//		Args:    cobra.NoArgs,
//		RunE: func(cmd *cobra.Command, _ []string) error {
//			clientCtx:=client.GetContext(cmd,cdc)
//			state, height, err := utils.QueryNodeConsensusState(clientCtx)
//			if err != nil {
//				return err
//			}
//
//			clientCtx = clientCtx.WithHeight(height)
//			return clientCtx.Response(state)
//		},
//	}
//
//	flags.AddQueryFlagsToCmd(cmd)
//
//	return cmd
//}
//
//// GetCmdParams returns the command handler for ibc client parameter querying.
//func GetCmdParams(cdc *codec.Codec) *cobra.Command {
//	cmd := &cobra.Command{
//		Use:     "params",
//		Short:   "Query the current ibc client parameters",
//		Long:    "Query the current ibc client parameters",
//		Args:    cobra.NoArgs,
//		Example: fmt.Sprintf("%s query %s %s params", version.AppName, host.ModuleName, types.SubModuleName),
//		RunE: func(cmd *cobra.Command, _ []string) error {
//			clientCtx:=client.GetContext(cmd,cdc)
//			queryClient := types.NewQueryClient(clientCtx)
//
//			res, _ := queryClient.ClientParams(cmd.Context(), &types.QueryClientParamsRequest{})
//			return clientCtx.Response(res.Params)
//		},
//	}
//
//	flags.AddQueryFlagsToCmd(cmd)
//
//	return cmd
//}