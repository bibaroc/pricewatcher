package amzn

type (
	Config struct {
		AmazonMarkerplaces    map[string]Categories `json:"amazon_markerplaces,omitempty" yaml:"amazon_markerplaces"`
		AmazonSellerBlacklist map[string]*string    `json:"amazon_seller_blacklist,omitempty" yaml:"amazon_seller_blacklist"`
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
