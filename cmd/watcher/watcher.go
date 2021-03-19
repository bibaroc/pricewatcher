package main

import (
	"context"
	"log"

	"github.com/bibaroc/pricewatcher/pkg/amzn"
	"golang.org/x/sync/errgroup"
)

func main() {
	cfg, err := readConfig()
	if err != nil {
		log.Println(err)
		return
	}

	g, ctx := errgroup.WithContext(context.Background())

	for marketPlaceURL, groups := range cfg.Amazon.Markerplaces {
		log.Printf("starting to watch %s", marketPlaceURL)
		g.Go(amzn.FetchResults(ctx, marketPlaceURL, groups, cfg.Amazon.SellerBlacklist, nil))
	}

	if err := g.Wait(); err != nil {
		log.Println(err)
	}
}
