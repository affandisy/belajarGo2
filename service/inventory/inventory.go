package inventory

type (
	Inventory struct {
		Code        string `json:"code"`
		Name        string `json:"name"`
		Stock       int    `json:"stock"`
		Description string `json:"description"`
		Status      string `json:"status"`
	}
)
