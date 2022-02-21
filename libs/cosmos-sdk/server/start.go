package server

// DONTCOVER

import (
	"fmt"
	"os"
	"runtime/pprof"

	"github.com/okex/exchain/libs/cosmos-sdk/baseapp"
	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/okex/exchain/libs/cosmos-sdk/client/lcd"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	"github.com/okex/exchain/libs/cosmos-sdk/store/flatkv"
	"github.com/okex/exchain/libs/cosmos-sdk/store/iavl"
	storetypes "github.com/okex/exchain/libs/cosmos-sdk/store/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	tmiavl "github.com/okex/exchain/libs/iavl"
	"github.com/okex/exchain/libs/system"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	tcmd "github.com/okex/exchain/libs/tendermint/cmd/tendermint/commands"
	"github.com/okex/exchain/libs/tendermint/libs/cli"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	tmos "github.com/okex/exchain/libs/tendermint/libs/os"
	"github.com/okex/exchain/libs/tendermint/mempool"
	"github.com/okex/exchain/libs/tendermint/node"
	"github.com/okex/exchain/libs/tendermint/p2p"
	pvm "github.com/okex/exchain/libs/tendermint/privval"
	"github.com/okex/exchain/libs/tendermint/proxy"
	"github.com/okex/exchain/libs/tendermint/rpc/client/local"
	"github.com/okex/exchain/libs/tendermint/state"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Tendermint full-node start flags
const (
	flagAddress            = "address"
	flagTraceStore         = "trace-store"
	flagCPUProfile         = "cpu-profile"
	FlagMinGasPrices       = "minimum-gas-prices"
	FlagHaltHeight         = "halt-height"
	FlagHaltTime           = "halt-time"
	FlagInterBlockCache    = "inter-block-cache"
	FlagUnsafeSkipUpgrades = "unsafe-skip-upgrades"
	FlagTrace              = "trace"

	FlagPruning           = "pruning"
	FlagPruningKeepRecent = "pruning-keep-recent"
	FlagPruningKeepEvery  = "pruning-keep-every"
	FlagPruningInterval   = "pruning-interval"
	FlagLocalRpcPort      = "local-rpc-port"
	FlagPortMonitor       = "netstat"
	FlagEvmImportPath     = "evm-import-path"
	FlagEvmImportMode     = "evm-import-mode"
	FlagGoroutineNum      = "goroutine-num"

	FlagPruningMaxWsNum = "pruning-max-worldstate-num"
	FlagExportKeystore  = "export-keystore"
	FlagLogServerUrl    = "log-server"
)

// StartCmd runs the service passed in, either stand-alone or in-process with
// Tendermint.
func StartCmd(ctx *Context,
	cdc *codec.Codec,
	appCreator AppCreator,
	appStop AppStop,
	registerRoutesFn func(restServer *lcd.RestServer),
	registerAppFlagFn func(cmd *cobra.Command),
	appPreRun func(ctx *Context) error,
	subFunc func(logger log.Logger) log.Subscriber,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Run the full node",
		Long: `Run the full node application with Tendermint in or out of process. By
default, the application will run with Tendermint in process.

Pruning options can be provided via the '--pruning' flag or alternatively with '--pruning-keep-recent',
'pruning-keep-every', and 'pruning-interval' together.

For '--pruning' the options are as follows:

default: the last 100 states are kept in addition to every 500th state; pruning at 10 block intervals
nothing: all historic states will be saved, nothing will be deleted (i.e. archiving node)
everything: all saved states will be deleted, storing only the current state; pruning at 10 block intervals
custom: allow pruning options to be manually specified through 'pruning-keep-recent', 'pruning-keep-every', and 'pruning-interval'

Node halting configurations exist in the form of two flags: '--halt-height' and '--halt-time'. During
the ABCI Commit phase, the node will check if the current block height is greater than or equal to
the halt-height or if the current block time is greater than or equal to the halt-time. If so, the
node will attempt to gracefully shutdown and the block will not be committed. In addition, the node
will not be able to commit subsequent blocks.

For profiling and benchmarking purposes, CPU profiling can be enabled via the '--cpu-profile' flag
which accepts a path for the resulting pprof file.
`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// app pre run
			if err := appPreRun(ctx); err != nil {
				return err
			}
			// set external package flags
			setExternalPackageValue(cmd)
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx.Logger.Info("Starting ABCI with Tendermint")

			sub := subFunc(ctx.Logger)
			log.SetSubscriber(sub)

			setPID(ctx)
			_, err := startInProcess(ctx, cdc, appCreator, appStop, registerRoutesFn)
			if err != nil {
				tmos.Exit(err.Error())
			}
			return nil
		},
	}

	// core flags for the ABCI application
	cmd.Flags().String(flagAddress, "tcp://0.0.0.0:26658", "Listen address")
	cmd.Flags().String(flagTraceStore, "", "Enable KVStore tracing to an output file")
	cmd.Flags().Bool(FlagTrace, false, "Provide full stack traces for errors in ABCI Log")
	cmd.Flags().String(
		FlagMinGasPrices, "",
		"Minimum gas prices to accept for transactions; Any fee in a tx must meet this minimum (e.g. 0.01photino;0.0001stake)",
	)
	cmd.Flags().IntSlice(FlagUnsafeSkipUpgrades, []int{}, "Skip a set of upgrade heights to continue the old binary")
	cmd.Flags().Uint64(FlagHaltHeight, 0, "Block height at which to gracefully halt the chain and shutdown the node")
	cmd.Flags().Uint64(FlagHaltTime, 0, "Minimum block time (in Unix seconds) at which to gracefully halt the chain and shutdown the node")
	cmd.Flags().Bool(FlagInterBlockCache, true, "Enable inter-block caching")
	cmd.Flags().String(flagCPUProfile, "", "Enable CPU profiling and write to the provided file")

	cmd.Flags().String(FlagPruning, storetypes.PruningOptionDefault, "Pruning strategy (default|nothing|everything|custom)")
	cmd.Flags().Uint64(FlagPruningKeepRecent, 0, "Number of recent heights to keep on disk (ignored if pruning is not 'custom')")
	cmd.Flags().Uint64(FlagPruningKeepEvery, 0, "Offset heights to keep on disk after 'keep-every' (ignored if pruning is not 'custom')")
	cmd.Flags().Uint64(FlagPruningInterval, 0, "Height interval at which pruned heights are removed from disk (ignored if pruning is not 'custom')")
	cmd.Flags().Uint64(FlagPruningMaxWsNum, 0, "Max number of historic states to keep on disk (ignored if pruning is not 'custom')")
	cmd.Flags().String(FlagLocalRpcPort, "", "Local rpc port for mempool and block monitor on cosmos layer(ignored if mempool/block monitoring is not required)")
	cmd.Flags().String(FlagPortMonitor, "", "Local target ports for connecting number monitoring(ignored if connecting number monitoring is not required)")
	cmd.Flags().String(FlagEvmImportMode, "default", "Select import mode for evm state (default|files|db)")
	cmd.Flags().String(FlagEvmImportPath, "", "Evm contract & storage db or files used for InitGenesis")
	cmd.Flags().Uint64(FlagGoroutineNum, 0, "Limit on the number of goroutines used to import evm data(ignored if evm-import-mode is 'default')")

	cmd.Flags().Bool(tmtypes.FlagDownloadDDS, false, "Download delta")
	cmd.Flags().Bool(tmtypes.FlagUploadDDS, false, "Upload delta")
	cmd.Flags().Bool(tmtypes.FlagAppendPid, false, "Append pid to the identity of delta producer")
	cmd.Flags().String(tmtypes.FlagRedisUrl, "localhost:6379", "redis url")
	cmd.Flags().String(tmtypes.FlagRedisAuth, "", "redis auth")
	cmd.Flags().Int(tmtypes.FlagRedisExpire, 300, "delta expiration time. unit is second")
	cmd.Flags().Int(tmtypes.FlagRedisDB, 0, "delta db num")
	cmd.Flags().Int(tmtypes.FlagDDSCompressType, 0, "delta compress type. 0|1|2|3")
	cmd.Flags().Int(tmtypes.FlagDDSCompressFlag, 0, "delta compress flag. 0|1|2")
	cmd.Flags().Int(tmtypes.FlagBufferSize, 10, "delta buffer size")
	cmd.Flags().String(FlagLogServerUrl, "", "log server url")
	cmd.Flags().Int(tmtypes.FlagDeltaVersion, tmtypes.DeltaVersion, "Specify delta version")

	cmd.Flags().Int(iavl.FlagIavlCacheSize, 1000000, "Max size of iavl cache")
	cmd.Flags().StringToInt(tmiavl.FlagOutputModules, map[string]int{"evm": 1, "acc": 1}, "decide which module in iavl to be printed")
	cmd.Flags().Int64(tmiavl.FlagIavlCommitIntervalHeight, 100, "Max interval to commit node cache into leveldb")
	cmd.Flags().Int64(tmiavl.FlagIavlMinCommitItemCount, 500000, "Min nodes num to triggle node cache commit")
	cmd.Flags().Int(tmiavl.FlagIavlHeightOrphansCacheSize, 8, "Max orphan version to cache in memory")
	cmd.Flags().Int(tmiavl.FlagIavlMaxCommittedHeightNum, 30, "Max committed version to cache in memory")
	cmd.Flags().Bool(tmiavl.FlagIavlEnableAsyncCommit, false, "Enable async commit")
	cmd.Flags().Bool(abci.FlagDisableABCIQueryMutex, false, "Disable local client query mutex for better concurrency")
	cmd.Flags().Bool(abci.FlagEnableWrappedTx, false, "Wrapped tx")
	cmd.Flags().Bool(abci.FlagDisableCheckTx, false, "Disable checkTx for test")
	cmd.Flags().MarkHidden(abci.FlagDisableCheckTx)
	cmd.Flags().Bool(abci.FlagCloseMutex, false, fmt.Sprintf("Deprecated in v0.19.13 version, use --%s instead.", abci.FlagDisableABCIQueryMutex))
	cmd.Flags().MarkHidden(abci.FlagCloseMutex)
	cmd.Flags().Bool(FlagExportKeystore, false, "export keystore file when call newaccount ")
	cmd.Flags().Bool(system.FlagEnableGid, false, "Display goroutine id in log")

	cmd.Flags().Int(state.FlagApplyBlockPprofTime, -1, "time(ms) of executing ApplyBlock, if it is higher than this value, save pprof")

	cmd.Flags().Float64Var(&baseapp.GasUsedFactor, baseapp.FlagGasUsedFactor, 0.4, "factor to calculate history gas used")

	cmd.Flags().BoolVar(&sdk.EnableMultiCache, sdk.FlagMultiCache, false, "Enable multi cache")
	cmd.Flags().Int(sdk.MaxAccInMultiCache, 0, "max acc in multi cache")
	cmd.Flags().Int(sdk.MaxStorageInMultiCache, 0, "max storage in multi cache")
	cmd.Flags().Bool(flatkv.FlagEnable, false, "Enable flat kv storage for read performance")

	// Don`t use cmd.Flags().*Var functions(such as cmd.Flags.IntVar) here, because it doesn't work with environment variables.
	// Use setExternalPackageValue function instead.
	viper.BindPFlag(FlagTrace, cmd.Flags().Lookup(FlagTrace))
	viper.BindPFlag(FlagPruning, cmd.Flags().Lookup(FlagPruning))
	viper.BindPFlag(FlagPruningKeepRecent, cmd.Flags().Lookup(FlagPruningKeepRecent))
	viper.BindPFlag(FlagPruningKeepEvery, cmd.Flags().Lookup(FlagPruningKeepEvery))
	viper.BindPFlag(FlagPruningInterval, cmd.Flags().Lookup(FlagPruningInterval))
	viper.BindPFlag(FlagPruningMaxWsNum, cmd.Flags().Lookup(FlagPruningMaxWsNum))
	viper.BindPFlag(FlagLocalRpcPort, cmd.Flags().Lookup(FlagLocalRpcPort))
	viper.BindPFlag(FlagPortMonitor, cmd.Flags().Lookup(FlagPortMonitor))
	viper.BindPFlag(FlagEvmImportMode, cmd.Flags().Lookup(FlagEvmImportMode))
	viper.BindPFlag(FlagEvmImportPath, cmd.Flags().Lookup(FlagEvmImportPath))
	viper.BindPFlag(FlagGoroutineNum, cmd.Flags().Lookup(FlagGoroutineNum))

	cmd.Flags().Bool(state.FlagParalleledTx, false, "Enable Parallel Tx")
	registerRestServerFlags(cmd)
	registerAppFlagFn(cmd)
	// add support for all Tendermint-specific command line options
	tcmd.AddNodeFlags(cmd)
	cmd.AddCommand(nodeModeCmd(ctx))
	return cmd
}

func startInProcess(ctx *Context, cdc *codec.Codec, appCreator AppCreator, appStop AppStop,
	registerRoutesFn func(restServer *lcd.RestServer)) (*node.Node, error) {

	cfg := ctx.Config
	home := cfg.RootDir
	//startInProcess hooker
	callHooker(FlagHookstartInProcess, ctx)

	traceWriterFile := viper.GetString(flagTraceStore)
	db, err := openDB(home)
	if err != nil {
		return nil, err
	}

	traceWriter, err := openTraceWriter(traceWriterFile)
	if err != nil {
		return nil, err
	}

	app := appCreator(ctx.Logger, db, traceWriter)

	nodeKey, err := p2p.LoadOrGenNodeKey(cfg.NodeKeyFile())
	if err != nil {
		return nil, err
	}

	// create & start tendermint node
	tmNode, err := node.NewNode(
		cfg,
		pvm.LoadOrGenFilePV(cfg.PrivValidatorKeyFile(), cfg.PrivValidatorStateFile()),
		nodeKey,
		proxy.NewLocalClientCreator(app),
		node.DefaultGenesisDocProviderFunc(cfg),
		node.DefaultDBProvider,
		node.DefaultMetricsProvider(cfg.Instrumentation),
		ctx.Logger.With("module", "node"),
	)
	if err != nil {
		return nil, err
	}

	app.SetOption(abci.RequestSetOption{
		Key:   "CheckChainID",
		Value: tmNode.ConsensusState().GetState().ChainID,
	})

	ctx.Logger.Info("startInProcess",
		"ConsensusStateChainID", tmNode.ConsensusState().GetState().ChainID,
		"GenesisDocChainID", tmNode.GenesisDoc().ChainID,
	)
	if err := tmNode.Start(); err != nil {
		return nil, err
	}

	var cpuProfileCleanup func()

	if cpuProfile := viper.GetString(flagCPUProfile); cpuProfile != "" {
		f, err := os.Create(cpuProfile)
		if err != nil {
			return nil, err
		}

		ctx.Logger.Info("starting CPU profiler", "profile", cpuProfile)
		if err := pprof.StartCPUProfile(f); err != nil {
			return nil, err
		}

		cpuProfileCleanup = func() {
			ctx.Logger.Info("stopping CPU profiler", "profile", cpuProfile)
			pprof.StopCPUProfile()
			f.Close()
		}
	}

	TrapSignal(func() {
		if tmNode.IsRunning() {
			_ = tmNode.Stop()
		}
		appStop(app)

		if cpuProfileCleanup != nil {
			cpuProfileCleanup()
		}

		ctx.Logger.Info("exiting...")
	})

	if registerRoutesFn != nil {
		go lcd.StartRestServer(cdc, registerRoutesFn, tmNode, viper.GetString(FlagListenAddr))
	}

	baseapp.SetGlobalMempool(tmNode.Mempool(), cfg.Mempool.SortTxByGp, cfg.Mempool.EnablePendingPool)

	if cfg.Mempool.EnablePendingPool {
		cliCtx := context.NewCLIContext().WithCodec(cdc)
		cliCtx.Client = local.New(tmNode)
		cliCtx.TrustNode = true
		accRetriever := types.NewAccountRetriever(cliCtx)
		tmNode.Mempool().SetAccountRetriever(accRetriever)
	}

	if parser, ok := app.(mempool.TxInfoParser); ok {
		tmNode.Mempool().SetTxInfoParser(parser)
	}

	// run forever (the node will not be returned)
	select {}
}

// Use setExternalPackageValue to set external package config value.
func setExternalPackageValue(cmd *cobra.Command) {
	iavl.IavlCacheSize = viper.GetInt(iavl.FlagIavlCacheSize)
	tmiavl.OutputModules, _ = cmd.Flags().GetStringToInt(tmiavl.FlagOutputModules)
	tmiavl.CommitIntervalHeight = viper.GetInt64(tmiavl.FlagIavlCommitIntervalHeight)
	tmiavl.MinCommitItemCount = viper.GetInt64(tmiavl.FlagIavlMinCommitItemCount)
	tmiavl.HeightOrphansCacheSize = viper.GetInt(tmiavl.FlagIavlHeightOrphansCacheSize)
	tmiavl.MaxCommittedHeightNum = viper.GetInt(tmiavl.FlagIavlMaxCommittedHeightNum)
	tmiavl.EnableAsyncCommit = viper.GetBool(tmiavl.FlagIavlEnableAsyncCommit)
	system.EnableGid = viper.GetBool(system.FlagEnableGid)

	state.ApplyBlockPprofTime = viper.GetInt(state.FlagApplyBlockPprofTime)
	state.HomeDir = viper.GetString(cli.HomeFlag)

	abci.SetDisableABCIQueryMutex(viper.GetBool(abci.FlagDisableABCIQueryMutex))
	abci.SetDisableCheckTx(viper.GetBool(abci.FlagDisableCheckTx))

	tmtypes.DownloadDelta = viper.GetBool(tmtypes.FlagDownloadDDS)
	tmtypes.UploadDelta = viper.GetBool(tmtypes.FlagUploadDDS)
	tmtypes.FastQuery = viper.GetBool(tmtypes.FlagFastQuery)
	tmtypes.DeltaVersion = viper.GetInt(tmtypes.FlagDeltaVersion)
}
