package base

//import (
//	"encoding/json"
//	"github.com/gorilla/mux"
//	"github.com/grpc-ecosystem/grpc-gateway/runtime"
//	"github.com/okex/exchain/libs/cosmos-sdk/codec"
//	"github.com/okex/exchain/libs/cosmos-sdk/types/module"
//	abci "github.com/okex/exchain/libs/tendermint/abci/types"
//	"github.com/spf13/cobra"
//)
//
//var (
//	_ module.AppModuleAdapter    = IBCAppModuleProxu{}
//	_ module.AppModuleBasic      = IBCAppModuleProxuBasic{}
//	_ module.AppModuleSimulation = IBCAppModuleProxu{}
//)
//
//type IBCAppModuleProxu struct {
//	impl module.AppModuleAdapter
//}
//
//func (p IBCAppModuleProxu) Name() string {
//	return "ibc"
//}
//
//func (p IBCAppModuleProxu) RegisterCodec(codec *codec.Codec) {
//	p.impl.RegisterCodec(codec)
//}
//
//func (p IBCAppModuleProxu) DefaultGenesis() json.RawMessage {
//	return p.impl.DefaultGenesis()
//}
//
//func (p IBCAppModuleProxu) ValidateGenesis(message json.RawMessage) error {
//}
//
//func (p IBCAppModuleProxu) RegisterRESTRoutes(context context.CLIContext, router *mux.Router) {
//	panic("implement me")
//}
//
//func (p IBCAppModuleProxu) GetTxCmd(codec *codec.Codec) *cobra.Command {
//	panic("implement me")
//}
//
//func (p IBCAppModuleProxu) GetQueryCmd(codec *codec.Codec) *cobra.Command {
//	panic("implement me")
//}
//
//func (p IBCAppModuleProxu) InitGenesis(context sdk.Context, message json.RawMessage) []abci.ValidatorUpdate {
//	panic("implement me")
//}
//
//func (p IBCAppModuleProxu) ExportGenesis(context sdk.Context) json.RawMessage {
//	panic("implement me")
//}
//
//func (p IBCAppModuleProxu) RegisterInterfaces(registry types.InterfaceRegistry) {
//	panic("implement me")
//}
//
//func (p IBCAppModuleProxu) RegisterGRPCGatewayRoutes(context context.CLIContext, mux *runtime.ServeMux) {
//	panic("implement me")
//}
//
//func (p IBCAppModuleProxu) RegisterInvariants(registry sdk.InvariantRegistry) {
//	panic("implement me")
//}
//
//func (p IBCAppModuleProxu) RegisterServices(configurator module.Configurator) {
//	panic("implement me")
//}
//
//func (p IBCAppModuleProxu) Upgrade(req *abci.UpgradeReq) (*abci.ModuleUpgradeResp, error) {
//	panic("implement me")
//}
//
//type IBCAppModuleProxuBasic struct {
//}
