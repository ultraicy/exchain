package types

import (
	codec2 "github.com/okex/exchain/libs/cosmos-sdk/codec"
	codectypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
)

func (acc BaseAccount) String() string {
	out, _ := acc.MarshalYAML()
	return out.(string)
}

// MarshalYAML returns the YAML representation of an account.
func (acc BaseAccount) MarshalYAML() (interface{}, error) {
	bz, err := codec2.MarshalYAML(codec2.NewProtoCodec(codectypes.NewInterfaceRegistry()), &acc)
	if err != nil {
		return nil, err
	}
	return string(bz), err
}



func(this *Params)String()string{
	return ""
}