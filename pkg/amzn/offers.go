package amzn

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/bibaroc/pricewatcher/pkg"
)

//nolint
var priceReplacer = strings.NewReplacer("\n", " ", "  ", " ", ".", "", "\u00a0", "")

type Offer struct {
	Condition    string
	ShippingFrom string
	Notes        string
	Price        int64
}

func getOffers(r io.Reader) ([]Offer, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("failed to parse document: %w", err)
	}

	var (
		offers = make([]Offer, 0)
		topErr error
	)

	doc.Find("#aod-offer").Each(func(_ int, s *goquery.Selection) {
		condition := pkg.SingleLineString(s.Find("#aod-offer-heading h5").Text())
		wholePriceS := priceReplacer.Replace(s.Find("span.a-price-whole").Text())
		wholePrice, err := strconv.ParseInt(wholePriceS[:len(wholePriceS)-1], 10, 64)
		if err != nil {
			topErr = fmt.Errorf("failed to parse %s as a number: %w", wholePriceS, err)
			return
		}

		shippingFrom := pkg.SingleLineString(s.Find("#aod-offer-shipsFrom span.a-color-base").Text())
		notes := pkg.SingleLineString(s.Find("#condition-text-block-title").Text())

		offers = append(offers, Offer{
			Condition:    condition,
			Price:        wholePrice,
			ShippingFrom: shippingFrom,
			Notes:        notes,
		})
	})

	if topErr != nil {
		return nil, topErr
	}

	return offers, nil
}
