package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/bibaroc/pricewatcher/pkg/amzn"
	"github.com/bibaroc/pricewatcher/pkg/tlgrm"
	"golang.org/x/sync/errgroup"
)

func main() {
	cfg, err := readConfig()
	if err != nil {
		log.Println(err)
		return
	}
	sendToAllGroups, err := tlgrm.NewGroupSender(
		mustString("TELEGRAM__TOKEN"),
		mustI64A("TELEGRAM__CHAT_IDS")...)
	if err != nil {
		log.Println(err)
		return
	}

	g, ctx := errgroup.WithContext(context.Background())

	for marketPlaceURL, groups := range cfg.Amazon.Markerplaces {
		log.Printf("starting to watch %s", marketPlaceURL)

		offerHandler := amzn.NewOfferHandler(sendToAllGroups)

		g.Go(amzn.FetchResults(ctx, marketPlaceURL, groups, cfg.Amazon.SellerBlacklist, offerHandler))
	}

	if err := g.Wait(); err != nil {
		log.Println(err)
	}
}

func mustString(s string) string {
	v, ok := os.LookupEnv(s)
	if !ok {
		panic(fmt.Sprintf("environment value for %q not set", s))
	}

	return v
}

func mustI64A(s string) []int64 {
	acc := []int64{}

	for _, stringValue := range strings.Split(mustString(s), ",") {
		intValue, err := strconv.ParseInt(stringValue, 10, 64)
		if err != nil {
			panic(fmt.Errorf("bad environment value for %q: cannot parse %s as integer: %w", s, stringValue, err))
		}

		acc = append(acc, intValue)
	}

	return acc
}
