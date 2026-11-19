package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
)

// --- Money makers: flipping + vendor snipes --------------------------------
//
// Both are pure math on data we already cache — no upstream calls, no AI.
//
// 1. FLIPS: items listed ≥20% under their 7-day usual price that actually
//    sell (possible_24h ≥ 1). Buy now, resell at the usual price, minus the
//    auction house cut. Ranked by expected daily profit (profit × sold/24h).
//    Caveat: the usual price comes from blackwater's baseline — the same
//    baseline our ⚠️ flags watch, so deep-discount rows stay suspicious.
//
// 2. VENDOR SNIPES: items whose cheapest listing is BELOW what a vendor pays.
//    Buy → sell to a vendor → guaranteed profit, no auctioning needed.
//
// Data source sizes: deals = top 100 (upstream limit), listings = 2,000.

type ahFlip struct {
	ItemID      int64   `json:"item_id"`
	Name        string  `json:"name"`
	Icon        string  `json:"icon"`
	BuyCopper   float64 `json:"buy_copper"`
	SellCopper  float64 `json:"sell_copper"`
	ProfitEach  float64 `json:"profit_each"`
	Sold24h     int64   `json:"sold_24h"`
	Units       int64   `json:"units"`
	Discount    float64 `json:"discount"`
	DailyProfit float64 `json:"daily_profit"` // profit_each × sold_24h
}

type ahSnipe struct {
	ItemID      int64   `json:"item_id"`
	Name        string  `json:"name"`
	Icon        string  `json:"icon"`
	BuyCopper   float64 `json:"buy_copper"`
	VendorCopper float64 `json:"vendor_copper"`
	ProfitEach  float64 `json:"profit_each"`
	Units       int64   `json:"units"`
	TotalProfit float64 `json:"total_profit"` // profit_each × units available
}

// ahFlips serves ranked flipping opportunities from the cached deals.
func ahFlips(c *gin.Context, db *sql.DB) {
	var body string
	if err := db.QueryRow(`SELECT body FROM ah_cache WHERE key = 'deals'`).Scan(&body); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "deals not polled yet"})
		return
	}

	var d struct {
		AuctionCut float64 `json:"auction_cut"`
		Items      []struct {
			ItemID   int64   `json:"item_id"`
			Name     string  `json:"name"`
			Icon     string  `json:"icon"`
			Min      float64 `json:"minimum_ask_copper"`
			Usual    float64 `json:"usual_copper"`
			Discount float64 `json:"discount"`
			Sold     int64   `json:"possible_24h"`
			Units    int64   `json:"units"`
		} `json:"items"`
	}
	if err := json.Unmarshal([]byte(body), &d); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "bad deals cache"})
		return
	}
	// auction cut: upstream sends a fraction (0.05); default to the classic 5%.
	cut := d.AuctionCut
	if cut <= 0 || cut >= 1 {
		cut = 0.05
	}

	flips := []ahFlip{}
	for _, it := range d.Items {
		if it.Sold < 1 || it.Min <= 0 || it.Usual <= 0 {
			continue // nothing sells it, or no usable prices
		}
		profit := it.Usual*(1-cut) - it.Min
		if profit <= 0 {
			continue
		}
		flips = append(flips, ahFlip{
			ItemID: it.ItemID, Name: it.Name, Icon: it.Icon,
			BuyCopper: it.Min, SellCopper: it.Usual, ProfitEach: profit,
			Sold24h: it.Sold, Units: it.Units, Discount: it.Discount,
			DailyProfit: profit * float64(it.Sold),
		})
	}
	sort.Slice(flips, func(i, j int) bool {
		return flips[i].DailyProfit > flips[j].DailyProfit
	})
	if len(flips) > 30 {
		flips = flips[:30]
	}
	c.JSON(http.StatusOK, gin.H{"items": flips, "auction_cut": cut})
}

// ahVendorSnipes serves items cheaper than their vendor price.
func ahVendorSnipes(c *gin.Context, db *sql.DB) {
	listings := listingIndex(db)
	snipes := []ahSnipe{}
	for _, li := range listings {
		if li.VendorSell <= 0 || li.MinAsk <= 0 || li.MinAsk >= li.VendorSell {
			continue
		}
		profit := li.VendorSell - li.MinAsk
		snipes = append(snipes, ahSnipe{
			ItemID: li.ID, Name: li.Name, Icon: li.Icon,
			BuyCopper: li.MinAsk, VendorCopper: li.VendorSell,
			ProfitEach: profit, Units: li.Units,
			TotalProfit: profit * float64(li.Units),
		})
	}
	sort.Slice(snipes, func(i, j int) bool {
		return snipes[i].TotalProfit > snipes[j].TotalProfit
	})
	if len(snipes) > 30 {
		snipes = snipes[:30]
	}
	c.JSON(http.StatusOK, gin.H{"items": snipes})
}