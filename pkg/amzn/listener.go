package amzn

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/bibaroc/pricewatcher/pkg/msg"
)

func NewOfferHandler(messageStore msg.MessageRepository, sender func(string) error) func(string, string, string, Offer) {
	return func(domain, group, asin string, offer Offer) {
		messageText := fmt.Sprintf(
			"Found a %s, on %s.\n\n"+
				"Sold as %s by %s for %d.\n"+
				"Notes: %s\n\n"+
				"Check in out at https://%s/dp/%s/",
			group, domain,
			offer.Condition, offer.ShippingFrom, offer.Price,
			offer.Notes,
			domain, asin)
		messageID := doSha(messageText)

		message, err := messageStore.Get(context.Background(), messageID)
		if err != nil {
			log.Println(err)
			return
		} else if message == nil {
			if err := sender(messageText); err != nil {
				log.Println(err)
				return
			}
			if err := messageStore.Save(context.Background(), messageID, messageText); err != nil {
				log.Println(err)
			}
		} else {
			fmt.Printf("skipping message %s for a %s dating to %s\n", messageID, group, message.CreatedAt.Format(time.RFC3339))
		}
	}
}

func doSha(s string) string {
	hash := sha256.New()
	_, _ = hash.Write([]byte(s))
	return hex.EncodeToString(hash.Sum(nil))
}
