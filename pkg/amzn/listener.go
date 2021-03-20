package amzn

import (
	"fmt"
	"log"
)

func NewOfferHandler(sender func(string) error) func(string, string, string, Offer) {
	return func(domain, group, asin string, offer Offer) {
		if err := sender(fmt.Sprintf(
			"Found a %s, on %s.\n\n"+
				"Sold as %s by %s for %d.\n"+
				"Notes: %s\n\n"+
				"Check in out at https://%s/dp/%s/",
			group, domain,
			offer.Condition, offer.ShippingFrom, offer.Price,
			offer.Notes,
			domain, asin)); err != nil {
			log.Println(err)
		}
	}
}
