package core

type CreditCard struct {
	Name      string `json:"name"`
	Number    string `json:"number"`
	Expiry    string `json:"expiry"`
	CVV       string `json:"cvv"`
	ZipCode   string `json:"zip_code"`
	Type      string `json:"type,omitempty"`
	VaultID   string `json:"vault_id,omitempty"`
	IsDefault bool   `json:"is_default"`
}
