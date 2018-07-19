package account

import "github.com/mixbee/mixbee-crypto/keypair"

type Controller struct {
	ID     string `json:"id"`
	Public string `json:"publicKey,omitemtpy"`
	keypair.PublicKey
}

type Identity struct {
	ID      string       `json:"ontid"`
	Label   string       `json:"label,omitempty"`
	Lock    bool         `json:"lock"`
	Control []Controller `json:"controls,omitempty"`
	Extra   interface{}  `json:"extra,omitempty"`
}
