package sign

import (
	"github.com/kaifei-bianjie/mock/util/helper/tx"
	"github.com/kaifei-bianjie/mock/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kaifei-bianjie/mock/conf"
	"github.com/kaifei-bianjie/mock/util/helper"
	"bytes"
	"fmt"
	"log"
	"github.com/kaifei-bianjie/mock/util/constants"
)

const (
	// Bech32PrefixAccAddr defines the Bech32 prefix of an account's address
	Bech32PrefixAccAddr = "faa"
	// Bech32PrefixAccPub defines the Bech32 prefix of an account's public key
	Bech32PrefixAccPub = "fap"
	// Bech32PrefixValAddr defines the Bech32 prefix of a validator's operator address
	Bech32PrefixValAddr = "fva"
	// Bech32PrefixValPub defines the Bech32 prefix of a validator's operator public key
	Bech32PrefixValPub = "fvp"
	// Bech32PrefixConsAddr defines the Bech32 prefix of a consensus node address
	Bech32PrefixConsAddr = "fca"
	// Bech32PrefixConsPub defines the Bech32 prefix of a consensus node public key
	Bech32PrefixConsPub = "fcp"
)

var (
	Cdc *codec.Codec
)

//custom tx codec
func init() {
	InitBech32()

	var cdc = codec.New()
	bank.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	Cdc = cdc
}

func InitBech32() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(Bech32PrefixAccAddr, Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(Bech32PrefixValAddr, Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(Bech32PrefixConsAddr, Bech32PrefixConsPub)
	config.Seal()
}

// sign tx
func signTx(unsignedTx auth.StdTx, senderInfo types.AccountInfo) ([]byte, error) {
	// build request
	accountNumber, err := helper.ConvertStrToInt64(senderInfo.AccountNumber)
	if err != nil {
		return nil, err
	}
	sequence, err := helper.ConvertStrToInt64(senderInfo.Sequence)
	if err != nil {
		return nil, err
	}
	signTxReq := types.SignTxReq{
		Tx:            unsignedTx,
		Name:          senderInfo.LocalAccountName,
		Password:      senderInfo.Password,
		ChainID:       conf.MockChainId,
		AccountNumber: accountNumber,
		Sequence:      sequence,
		AppendSig:     true,
	}

	// send sign tx request
	reqBytes, err := Cdc.MarshalJSON(signTxReq)
	if err != nil {
		return nil, err
	}
	reqBuffer := bytes.NewBuffer(reqBytes)
	statusCode, resBytes, err := helper.HttpClientPostJsonData(constants.UriTxSign, reqBuffer)

	// handle response
	if err != nil {
		return nil, err
	}

	if statusCode != constants.StatusCodeOk {
		return nil, fmt.Errorf("unexcepted status code: %v", statusCode)
	}

	return resBytes, nil
}

// broadcast signed tx
func BroadcastSignedTx(senderInfo types.AccountInfo, receiver string) ([]byte, error) {
	var (
		unsignedTx, signedTx auth.StdTx
	)

	// build unsigned tx
	unsignedTxBytes, err := tx.SendTransferTx(senderInfo, receiver, true)
	if err != nil {
		log.Printf("build unsigned tx failed: %v\n", err)
		return nil, err
	}
	err = Cdc.UnmarshalJSON(unsignedTxBytes, &unsignedTx)
	if err != nil {
		log.Printf("build unsigned tx failed: %v\n", err)
		return nil, err
	}

	// sign tx
	signedTxBytes, err := signTx(unsignedTx, senderInfo)
	if err != nil {
		log.Printf("sign tx failed: %v\n", err)
		return nil, err
	}
	err = Cdc.UnmarshalJSON(signedTxBytes, &signedTx)
	if err != nil {
		log.Printf("sign tx failed: %v\n", err)
		return nil, err
	}

	// broadcast signed tx
	broadcastTxReq := types.BoradcaseTxReq{
		Tx: signedTx,
	}
	reqBytes, err := Cdc.MarshalJSON(broadcastTxReq)
	if err != nil {
		log.Printf("broadcast tx failed: %v\n", err)
		return nil, err
	}
	reqBuffer := bytes.NewBuffer(reqBytes)

	statusCode, resBytes, err := helper.HttpClientPostJsonData(constants.UriTxBroadcastTx, reqBuffer)

	if err != nil {
		log.Printf("broadcast tx failed: %v\n", err)
		return nil, err
	}

	if statusCode != constants.StatusCodeOk {
		log.Printf("broadcast tx failed, unexcepted status code: %v\n", err)
		return nil, err
	}

	return resBytes, nil
}