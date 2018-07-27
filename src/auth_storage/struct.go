package auth_storage

type Key struct {
	Access      string   `json:"access"`
	Secret      string   `json:"secret"`
	ProductList []string `json:"product_list"`
}

type Keys struct {
	Version string `json:"version"`
	Keys    []Key  `json:"keys"`
}
