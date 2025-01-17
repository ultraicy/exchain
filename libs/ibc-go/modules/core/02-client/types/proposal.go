package types

import (
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	govtypes "github.com/okex/exchain/libs/cosmos-sdk/x/gov/types"
)

const (
	// ProposalTypeClientUpdate defines the type for a ClientUpdateProposal
	ProposalTypeClientUpdate = "ClientUpdate"
)

var (
	_ govtypes.Content = &ClientUpdateProposal{}
	//_ codectypes.UnpackInterfacesMessage = &UpgradeProposal{}
)

func init() {
	govtypes.RegisterProposalType(ProposalTypeClientUpdate)
}

// NewClientUpdateProposal creates a new client update proposal.
// NewClientUpdateProposal creates a new client update proposal.
func NewClientUpdateProposal(title, description, subjectClientID, substituteClientID string) govtypes.Content {
	return &ClientUpdateProposal{
		Title:              title,
		Description:        description,
		SubjectClientId:    subjectClientID,
		SubstituteClientId: substituteClientID,
	}
}

// GetTitle returns the title of a client update proposal.
func (cup *ClientUpdateProposal) GetTitle() string { return cup.Title }

// GetDescription returns the description of a client update proposal.
func (cup *ClientUpdateProposal) GetDescription() string { return cup.Description }

// ProposalRoute returns the routing key of a client update proposal.
func (cup *ClientUpdateProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of a client update proposal.
func (cup *ClientUpdateProposal) ProposalType() string { return ProposalTypeClientUpdate }

// ValidateBasic runs basic stateless validity checks
func (cup *ClientUpdateProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(cup)
	if err != nil {
		return err
	}

	if cup.SubjectClientId == cup.SubstituteClientId {
		return sdkerrors.Wrap(ErrInvalidSubstitute, "subject and substitute client identifiers are equal")
	}
	if _, _, err := ParseClientIdentifier(cup.SubjectClientId); err != nil {
		return err
	}
	if _, _, err := ParseClientIdentifier(cup.SubstituteClientId); err != nil {
		return err
	}

	return nil
}

// UnpackInterfaces implements the UnpackInterfacesMessage interface.
// func (cup ClientUpdateProposal) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
// 	var header exported.Header
// 	return unpacker.UnpackAny(cup.Header, &header)
// }
