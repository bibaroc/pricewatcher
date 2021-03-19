package amzn

type (
	Config struct {
		Markerplaces    map[string]Categories `json:"markerplaces,omitempty" yaml:"markerplaces"`
		SellerBlacklist map[string]*string    `json:"seller_blacklist,omitempty" yaml:"seller_blacklist"`
	}
	Categories map[string]Group
	Group      struct {
		MaxPrice *int64 `json:"max_price,omitempty" yaml:"max_price"`
		Items    []Item `json:"items,omitempty" yaml:"items"`
	}
	Item struct {
		MaxPrice *int64 `json:"max_price,omitempty" yaml:"max_price"`
		ASIN     string `json:"asin,omitempty" yaml:"asin"`
	}
)
