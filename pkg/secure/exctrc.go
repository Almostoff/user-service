package secure

import (
	"encoding/base64"
	"encoding/json"
	eciesgo "github.com/ecies/go"
	"time"
)

func GenerateExctrc(params *GenerateExctrcParams, publicKey string) (string, error) {
	body := GenerateExctrcBody{
		Client:          "",
		Recipient:       params.ReceiveWallet,
		RecipientSystem: params.ReceiveSystem,
		CurrencyFrom:    params.CurrencyFrom,
		CurrencyTo:      params.CurrencyTo,
		Date:            time.Now().Format("2006-01-02T15:04:05"),
		Txn:             params.RequestId,
		Ip:              params.UserIp,
		UserAgent:       "",
		Email:           params.UserEmail,
		Cid:             "",
		Owner:           params.SiteBaseLink,
		OwnerEmail:      params.SiteEmail,
	}
	data, err := json.Marshal(&body)
	if err != nil {
		return "", err
	}
	pubkey, err := eciesgo.NewPublicKeyFromHex(publicKey)
	if err != nil {
		return "", err
	}
	result, err := eciesgo.Encrypt(pubkey, data)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(result), nil
}
