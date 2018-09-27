package main

import (
	"github.com/gin-gonic/gin"
	"github.com/json-iterator/go"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const cmcTickerApi = "https://api.coinmarketcap.com/v1/ticker/?limit=9999"

type Ticker struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Symbol           string `json:"symbol"`
	Rank             string `json:"rank"`
	PriceUSD         string `json:"price_usd"`
	PriceBTC         string `json:"price_btc"`
	Volume24hUSD     string `json:"24h_volume_usd"`
	MarketCapUSD     string `json:"market_cap_usd"`
	AvailableSupply  string `json:"available_supply"`
	TotalSupply      string `json:"total_supply"`
	MaxSupply        string `json:"max_supply"`
	PercentChange1h  string `json:"percent_change_1h"`
	PercentChange24h string `json:"percent_change_24h"`
	PercentChange7d  string `json:"percent_change_7d"`
	LastUpdated      string `json:"last_updated"`
}

var tickerMap = make(map[string]Ticker)
var tickerList = make([]Ticker, 0)

func main() {
	log.Info("started!")
	go updateCmcTicker()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.GET("/ticker", func(c *gin.Context) {
		c.JSON(http.StatusOK, tickerList)
	})

	r.GET("/ticker/:id", func(c *gin.Context) {
		id := c.Param("id")

		if ticker, ok := tickerMap[id]; ok {
			// keep same format with cmc
			oneTickerList := [1]Ticker{ticker}
			c.JSON(http.StatusOK, oneTickerList)
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "id not found"})
		}
	})

	r.Run()
}

func updateCmcTicker() {
	for true {
		log.Info("update from cmc.")
		resp, err := http.Get(cmcTickerApi)
		if err != nil {
			panic(err)
		}

		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&tickerList)
		if err != nil {
			panic(err)
		}

		for _, ticker := range tickerList {
			tickerMap[ticker.ID] = ticker
		}
		log.Info("updated!")
		time.Sleep(30 * time.Second)
	}
}
