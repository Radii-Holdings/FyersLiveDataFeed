package main

import (
	"fmt"

	"github.com/rishi-anand/fyers-go-client"
)

func main() {
	apiKey := "85VLN1I8IV-100"
	accessToken := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJhcGkuZnllcnMuaW4iLCJpYXQiOjE2NjUxMTQ4OTAsImV4cCI6MTY2NTE4OTA1MCwibmJmIjoxNjY1MTE0ODkwLCJhdWQiOlsieDowIiwieDoxIiwieDoyIiwiZDoxIiwiZDoyIiwieDoxIiwieDowIl0sInN1YiI6ImFjY2Vzc190b2tlbiIsImF0X2hhc2giOiJnQUFBQUFCalA2TUtvMXppTjNPaVEzQi1tX2phMmo0WDNSQlJiaWFYZGxTUmlva1VkLUQ2d0wtamdyY29aTkVKTHlOVUJUVV90dWpqMmVJVlFmSzJOYTd5SlZFTUVDNnBlT2JMV1BzZ0pDRy1hWDB2TEQxclRIYz0iLCJkaXNwbGF5X25hbWUiOiJBTE9ZIEFESVRZQSBTRU4iLCJmeV9pZCI6IlhBMDEyODgiLCJhcHBUeXBlIjoxMDAsInBvYV9mbGFnIjoiTiJ9.EAnXW1yyY1YUiB-y8SCNw9Ww33ECE-bLnDUxwiIM7JM"
	symbols := []string{"NSE:SBIN-EQ", "NSE:ONGC-EQ"}

	cli := fyers.New(apiKey, accessToken)
	if quotes, err := cli.GetQuote(symbols); err != nil {
		fmt.Errorf("failed to get quote from fyers. %v", err)
	} else {
		fmt.Println(quotes)
	}
}
