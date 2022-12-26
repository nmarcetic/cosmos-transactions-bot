package messages

import (
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	cosmosBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/gogo/protobuf/proto"
	types2 "main/pkg/config/types"
	"main/pkg/data_fetcher"
	"main/pkg/types"
	"main/pkg/utils"
)

type MsgSend struct {
	From   types2.Link
	To     types2.Link
	Amount []*types.Amount
}

func ParseMsgSend(data []byte, chain *types2.Chain) (types.Message, error) {
	var parsedMessage cosmosBankTypes.MsgSend
	if err := proto.Unmarshal(data, &parsedMessage); err != nil {
		return nil, err
	}

	return &MsgSend{
		From: chain.GetWalletLink(parsedMessage.FromAddress),
		To:   chain.GetWalletLink(parsedMessage.ToAddress),
		Amount: utils.Map(parsedMessage.Amount, func(coin cosmosTypes.Coin) *types.Amount {
			return &types.Amount{
				Value: float64(coin.Amount.Int64()),
				Denom: coin.Denom,
			}
		}),
	}, nil
}

func (m MsgSend) Type() string {
	return "MsgSend"
}

func (m *MsgSend) GetAdditionalData(fetcher data_fetcher.DataFetcher) {
	price, found := fetcher.GetPrice()
	if !found {
		return
	}

	for _, amount := range m.Amount {
		if amount.Denom != fetcher.Chain.BaseDenom {
			continue
		}

		amount.Value /= float64(fetcher.Chain.DenomCoefficient)
		amount.Denom = fetcher.Chain.DisplayDenom
		amount.PriceUSD = amount.Value * price
	}
}

func (m *MsgSend) GetValues() map[string]string {
	return map[string]string{
		"type": "MsgSend",
		"from": m.From.Value,
		"to":   m.To.Value,
	}
}
