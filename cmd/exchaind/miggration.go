package main

import (
	"fmt"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	ethstate "github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/okex/exchain/app"
	types2 "github.com/okex/exchain/app/types"
	"github.com/okex/exchain/libs/cosmos-sdk/server"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	authexported "github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	"github.com/okex/exchain/libs/mpt"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/spf13/cobra"
	"log"
	"path/filepath"
)

func migrateCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "miggrate-state",
		Short: "miggrate iavl state to mpt state",
	}

	cmd.AddCommand(
		migrateAccountCmd(ctx),
		migrateContractCmd(ctx),
		cleanRawDBCmd(ctx),
	)

	return cmd
}

func migrateAccountCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "miggrate-account",
		Short: "1. miggrate iavl account to mpt account",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("--------- miggrate account start ---------")
			migrateAccount(ctx)
			log.Println("--------- miggrate account end ---------")
		},
	}
	cmd.Flags().String(FlagDisplayContractAddr, "", "target contract address to display")
	cmd.Flags().Int64(FlagDisplayVersion, 0, "target state version to display")
	return cmd
}

func migrateContractCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "miggrate-contract",
		Short: "2. miggrate iavl contract state to mpt contract state",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("--------- display state start ---------")
			migrateContract(ctx)
			log.Println("--------- display state end ---------")
		},
	}
	cmd.Flags().String(FlagDisplayContractAddr, "", "target contract address to display")
	cmd.Flags().Int64(FlagDisplayVersion, 0, "target state version to display")
	return cmd
}

func cleanRawDBCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clean-rawdb",
		Short: "3. clean up miggrated iavl state",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("--------- display state start ---------")
			cleanRawDB(ctx)
			log.Println("--------- display state end ---------")
		},
	}
	cmd.Flags().String(FlagDisplayContractAddr, "", "target contract address to display")
	cmd.Flags().Int64(FlagDisplayVersion, 0, "target state version to display")
	return cmd
}

//----------------------------------------------------------------
func migrateAccount(ctx *server.Context) {
	migApp := newMigrationApp(ctx)

	ver, err := migApp.GetCommitVersion()
	panicError(err)

	// init deliver state
	migApp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: ver + 1}})
	cmCtx := migApp.GetDeliverStateCtx()

	accMptDb := mpt.InstanceOfAccStore()
	accTrie, err := accMptDb.OpenTrie(ethcmn.Hash{})
	panicError(err)

	evmMptDb := mpt.InstanceOfEvmStore()
	evmTrie, err := evmMptDb.OpenTrie(ethcmn.Hash{})
	panicError(err)

	cnt := 0
	contractCnt := 0
	emptyRootHashByte := types.EmptyRootHash.Bytes()

	migApp.AccountKeeper.MiggrateAccounts(cmCtx, func(account authexported.Account, key, value []byte) (stop bool) {
		cnt += 1
		err := accTrie.TryUpdate(key, value)
		panicError(err)

		if cnt % 100 == 0 {
			pushData2Database(accMptDb, accTrie, cmCtx.BlockHeight() - 1)
		}

		// contract account
		switch account.(type) {
		case *types2.EthAccount:
			contractCnt += 1

			ethAcc := account.(*types2.EthAccount)
			err = evmTrie.TryUpdate(ethAcc.EthAddress().Bytes(), emptyRootHashByte)
			panicError(err)

			if len(ethAcc.CodeHash) > 0 {
				cHash := ethcmn.BytesToHash(ethAcc.CodeHash)

				// migrate code
				codeWriter := evmMptDb.TrieDB().DiskDB().NewBatch()
				code := migApp.EvmKeeper.GetCodeByHash(cmCtx, cHash)
				rawdb.WriteCode(codeWriter, cHash, code)
				err = codeWriter.Write()
				panicError(err)
			}

			if contractCnt % 100 == 0 {
				pushData2Database(evmMptDb, evmTrie, cmCtx.BlockHeight() - 1)
			}
		default:
			//do nothing
		}

		return false
	})
	pushData2Database(accMptDb, accTrie, cmCtx.BlockHeight() - 1)
	pushData2Database(evmMptDb, evmTrie, cmCtx.BlockHeight() - 1)

	fmt.Println(fmt.Sprintf("Successfule migrate %d account (include %d contract account)", cnt, contractCnt))
}

func migrateContract(ctx *server.Context) {
	migApp := newMigrationApp(ctx)

	ver, err := migApp.GetCommitVersion()
	panicError(err)

	// init deliver state
	migApp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: ver + 1}})
	cmCtx := migApp.GetDeliverStateCtx()

	evmMptDb := mpt.InstanceOfEvmStore()
	rootHash := migApp.EvmKeeper.GetMptRootHash(uint64(cmCtx.BlockHeight() - 1))
	evmTrie, err := evmMptDb.OpenTrie(rootHash)
	panicError(err)

	cnt := 0
	itr := trie.NewIterator(evmTrie.NodeIterator(nil))
	for itr.Next() {
		cnt += 1

		addr := ethcmn.BytesToAddress(evmTrie.GetKey(itr.Key))
		addrHash := ethcrypto.Keccak256Hash(addr[:])
		contractTrie := getTrie(evmMptDb, addrHash)

		keyCnt := 0
		_ = migApp.EvmKeeper.ForEachStorage(cmCtx, addr, func(key, value ethcmn.Hash) bool {
			// Encoding []byte cannot fail, ok to ignore the error.
			v, _ := rlp.EncodeToBytes(ethcmn.TrimLeftZeroes(value[:]))
			err := contractTrie.TryUpdate(key[:], v)
			panicError(err)

			keyCnt += 1
			return false
		})
		rootHash, err := contractTrie.Commit(nil)
		panicError(err)

		fmt.Println(fmt.Sprintf("migrate contract %s with %d key-value", addr.String(), keyCnt))

		err = evmTrie.TryUpdate(addr[:], rootHash.Bytes())
		panicError(err)

		if cnt % 100 == 0 {
			pushData2Database(evmMptDb, evmTrie, cmCtx.BlockHeight() - 1)
		}
	}
	pushData2Database(evmMptDb, evmTrie, cmCtx.BlockHeight() - 1)

	fmt.Println(fmt.Sprintf("Successfule migrate %d contract stroage", cnt))
}

func cleanRawDB(ctx *server.Context) {
	fmt.Println("Not implement!!!")
}

//----------------------------------------------------------------

func pushData2Database(db ethstate.Database, tr ethstate.Trie, height int64) {
	var storageRoot ethcmn.Hash
	root, err := tr.Commit(func(_ [][]byte, _ []byte, leaf []byte, parent ethcmn.Hash) error {
		storageRoot.SetBytes(leaf)
		if storageRoot != types.EmptyRootHash {
			db.TrieDB().Reference(storageRoot, parent)
		}
		return nil
	})
	panicError(err)

	err = db.TrieDB().Commit(root, false, nil)
	panicError(err)

	setAccMptRootHash(db, uint64(height), root)

	fmt.Println("pushData2Database version: ", height, " root hash is: ", root.String())
}

func newMigrationApp(ctx *server.Context) *app.OKExChainApp {
	rootDir := ctx.Config.RootDir
	dataDir := filepath.Join(rootDir, "data")
	db, err := openDB(applicationDB, dataDir)
	if err != nil {
		panic("fail to open application db: " + err.Error())
	}

	return app.NewOKExChainApp(
		ctx.Logger,
		db,
		nil,
		true,
		map[int64]bool{},
		0,
	)
}

func getTrie(db ethstate.Database, addrHash ethcmn.Hash) ethstate.Trie {
	tr, _ := db.OpenStorageTrie(addrHash, ethcmn.Hash{})
	return tr
}

// SetMptRootHash sets the mapping from block height to root mpt hash
func setAccMptRootHash(db ethstate.Database, height uint64, hash ethcmn.Hash) {
	KeyPrefixRootMptHash := []byte{0x01}
	KeyPrefixLatestStoredHeight := []byte{0x02}

	hhash := sdk.Uint64ToBigEndian(height)
	db.TrieDB().DiskDB().Put(KeyPrefixLatestStoredHeight, hhash)
	db.TrieDB().DiskDB().Put(append(KeyPrefixRootMptHash, hhash...), hash.Bytes())
}