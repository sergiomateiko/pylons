package msgs

import (
	"encoding/json"

	"github.com/Pylons-tech/pylons/x/pylons/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgCreateRecipe defines a CreateRecipe message
type MsgCreateRecipe struct {
	// optional RecipeID if someone
	RecipeID      string `json:",omitempty"`
	Name          string
	CookbookID    string // the cookbook guid
	CoinInputs    types.CoinInputList
	ItemInputs    types.ItemInputList
	Entries       types.EntriesList
	Outputs       types.WeightedOutputsList
	BlockInterval int64
	Sender        sdk.AccAddress
	Description   string
}

// NewMsgCreateRecipe a constructor for CreateRecipe msg
func NewMsgCreateRecipe(recipeName, cookbookID, recipeID, description string,
	coinInputs types.CoinInputList,
	itemInputs types.ItemInputList,
	entries types.EntriesList,
	outputs types.WeightedOutputsList,
	blockInterval int64,
	sender sdk.AccAddress) MsgCreateRecipe {
	return MsgCreateRecipe{
		Name:          recipeName,
		CookbookID:    cookbookID,
		RecipeID:      recipeID,
		Description:   description,
		CoinInputs:    coinInputs,
		ItemInputs:    itemInputs,
		Entries:       entries,
		Outputs:       outputs,
		BlockInterval: int64(blockInterval),
		Sender:        sender,
	}
}

// Route should return the name of the module
func (msg MsgCreateRecipe) Route() string { return "pylons" }

// Type should return the action
func (msg MsgCreateRecipe) Type() string { return "create_recipe" }

// ValidateBasic validates the Msg
func (msg MsgCreateRecipe) ValidateBasic() sdk.Error {

	// validation for the item input index overflow on entries
	for _, entry := range msg.Entries {
		switch entry.(type) {
		case types.CoinOutput:
		case types.ItemOutput:
			itemOutput, _ := entry.(types.ItemOutput)
			if itemOutput.ItemInputRef > 0 {
				if itemOutput.ItemInputRef > len(msg.ItemInputs) {
					return sdk.ErrInternal("ItemInputRef overflow length of ItemInputs")
				}
				if itemOutput.ItemInputRef < 0 {
					return sdk.ErrInternal("ItemInputRef is less than 0 which is invalid")
				}
			}
			// TODO should do basic validation for program of ItemOutput weight
			// TODO should do basic validation coin output program
			// TODO should do basic validation double param program for ToModify
			// TODO should do basic validation string param program for ToModify
			// TODO should do basic validation long param program for ToModify
			// TODO should do basic validation double param program for ItemOutput (generation)
			// TODO should do basic validation string param program for ItemOutput (generation)
			// TODO should do basic validation long param program for ItemOutput (generation)
		default:
			return sdk.ErrInternal("invalid entry type available")
		}
	}

	// validation for same ItemInputRef on output
	for _, output := range msg.Outputs {
		usedItemInputRefs := make(map[int]bool)
		usedEntries := make(map[int]bool)
		for _, result := range output.Result {
			if result >= len(msg.Entries) || result < 0 {
				return sdk.ErrInternal("output is refering to index which is out of entries range")
			}
			if usedEntries[result] {
				return sdk.ErrInternal("double use of entries within single output result")
			}
			usedEntries[result] = true
			entry := msg.Entries[result]
			switch entry.(type) {
			case types.ItemOutput:
				itemOutput, _ := entry.(types.ItemOutput)
				if itemOutput.ItemInputRef > 0 {
					if usedItemInputRefs[itemOutput.ItemInputRef] {
						return sdk.ErrInternal("double use of item input within single output result")
					}
					usedItemInputRefs[itemOutput.ItemInputRef] = true
				}
			}
		}
	}

	if msg.Sender.Empty() {
		return sdk.ErrInvalidAddress(msg.Sender.String())
	}

	if len(msg.Description) < 20 {
		return sdk.ErrInternal("the description should have more than 20 characters")
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgCreateRecipe) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners gets the signer who should have signed the message
func (msg MsgCreateRecipe) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
