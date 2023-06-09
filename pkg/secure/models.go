package secure

type GenerateExctrcParams struct {
	ReceiveWallet string
	ReceiveSystem string
	CurrencyFrom  string
	CurrencyTo    string
	RequestId     int64
	UserIp        string
	UserEmail     string
	SiteBaseLink  string
	SiteEmail     string
}

type GenerateExctrcBody struct {
	Client          string `json:"client"`
	Recipient       string `json:"recipient"`
	RecipientSystem string `json:"recipient_system"`
	CurrencyFrom    string `json:"c_from"`
	CurrencyTo      string `json:"c_to"`
	Date            string `json:"date"`
	Txn             int64  `json:"txn"`
	Ip              string `json:"ip"`
	UserAgent       string `json:"user_agent"`
	Email           string `json:"email"`
	Cid             string `json:"cid"`
	Owner           string `json:"owner"`
	OwnerEmail      string `json:"owner_email"`
}
