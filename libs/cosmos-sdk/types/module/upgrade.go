package module

import (
	"errors"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

type UpgradeModule interface {
	ModuleName() string
	RegisterTask() HeightTask
	UpgradeHeight() int64
}

type HeightTasks []HeightTask

func (h HeightTasks) Len() int {
	return len(h)
}

func (h HeightTasks) Less(i, j int) bool {
	return h[i].GetOrderer() < h[j].GetOrderer()
}

func (h HeightTasks) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

type HeightTask interface {
	GetOrderer() int16
	Execute(c sdk.Context) error
	ValidateBasic() error
}
type heightTask struct {
	orderer      int16
	taskExecutor func(ctx sdk.Context) error
}

var (
	_ HeightTask = (*heightTask)(nil)
)

func NewHeightTask(orderer int16, taskExecutor func(ctx sdk.Context) error) HeightTask {
	return &heightTask{orderer: orderer, taskExecutor: taskExecutor}
}

func (t *heightTask) GetOrderer() int16 {
	return t.orderer
}

func (t *heightTask) ValidateBasic() error {
	if t.taskExecutor == nil {
		return errors.New("executor cant be nil")
	}

	return nil
}

func (t *heightTask) Execute(ctx sdk.Context) error {
	return t.taskExecutor(ctx)
}
