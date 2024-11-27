package main

import (
	"html/template"
	"net/http"
	"time"
	"trading_app_cli/db"
	"trading_app_cli/engine"
	"trading_app_gui/app"

	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	defer db.DB.Close()

	r := gin.Default()

	// カスタム関数の登録
	r.SetFuncMap(template.FuncMap{
		"formatDate": func(t time.Time) string {
			return t.Format("2006-01-02 15:04:05")
		},
	})

	// テンプレートのロード
	r.LoadHTMLGlob("templates/*.html")

	// ダッシュボードの表示
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "dashboard.html", nil)
	})

	// APIエンドポイント
	r.POST("/api/trade", func(c *gin.Context) {
		var trade struct {
			User   string  `json:"user"`
			Amount float64 `json:"amount"`
			Price  float64 `json:"price"`
			Type   string  `json:"type"`
		}
		if err := c.ShouldBindJSON(&trade); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		newOrder, err := db.CreateOrder(trade.User, trade.Type, trade.Amount, trade.Price)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
			return
		}

		matched, message, err := engine.MatchOrder(db.DB, newOrder)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to match order"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": message,
			"matched": matched,
		})
	})

	r.GET("/api/orders", func(c *gin.Context) {
		orders, err := app.FetchOrders(db.DB)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"orders": orders})
	})

	r.GET("/api/trades", func(c *gin.Context) {
		trades, err := app.FetchTrades(db.DB)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch trades"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"trades": trades})
	})

	r.Run(":8081")
}

// package main

// import (
// 	"html/template"
// 	"net/http"
// 	"strconv"
// 	"time"
// 	"trading_app_cli/db"
// 	"trading_app_cli/engine"
// 	"trading_app_gui/app"

// 	"github.com/gin-gonic/gin"
// )

// func main() {
// 	// データベースの初期化
// 	db.InitDB()
// 	defer db.DB.Close()

// 	r := gin.Default()
// 	// カスタム関数を登録
// 	r.SetFuncMap(template.FuncMap{
// 		"formatDate": func(t time.Time) string {
// 			return t.Format("2006-01-02 15:04:05")
// 		},
// 	})
// 	r.LoadHTMLGlob("templates/*")

// 	r.GET("/", func(c *gin.Context) {
// 		c.HTML(200, "index.html", nil)
// 	})

// 	// トレードページ
// 	r.GET("/trade", func(c *gin.Context) {
// 		c.HTML(http.StatusOK, "trade.html", nil)
// 	})

// 	// 買い注文
// 	r.POST("/trade/buy", func(c *gin.Context) {
// 		user := c.PostForm("user")
// 		amount, err := strconv.ParseFloat(c.PostForm("amount"), 64)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount"})
// 			return
// 		}

// 		price, err := strconv.ParseFloat(c.PostForm("price"), 64)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price"})
// 			return
// 		}

// 		newOrder, err := db.CreateOrder(user, "buy", amount, price)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
// 			return
// 		}
// 		c.JSON(http.StatusOK, gin.H{"message": "Buy order placed successfully"})

// 		isMatched, message, err := engine.MatchOrder(db.DB, newOrder)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
// 			return
// 		}
// 		c.JSON(http.StatusInternalServerError, gin.H{"trade_message": message})
// 		if !isMatched {
// 			c.JSON(http.StatusInternalServerError, gin.H{"trade_message": "Order added to the open order book."})
// 		}
// 	})

// 	// 売り注文
// 	r.POST("/trade/sell", func(c *gin.Context) {
// 		user := c.PostForm("user")
// 		amount, err := strconv.ParseFloat(c.PostForm("amount"), 64)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount"})
// 			return
// 		}
// 		price, err := strconv.ParseFloat(c.PostForm("price"), 64)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price"})
// 			return
// 		}

// 		newOrder, err := db.CreateOrder(user, "sell", amount, price)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
// 			return
// 		}
// 		c.JSON(http.StatusOK, gin.H{"message": "Sell order placed successfully"})

// 		isMatched, message, err := engine.MatchOrder(db.DB, newOrder)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
// 			return
// 		}
// 		c.JSON(http.StatusInternalServerError, gin.H{"trade_message": message})
// 		if !isMatched {
// 			c.JSON(http.StatusInternalServerError, gin.H{"trade_message": "Order added to the open order book."})
// 		}
// 	})
// 	// オーダー一覧表示
// 	r.GET("/orders", func(c *gin.Context) {
// 		orders, err := app.FetchOrders(db.DB)
// 		if err != nil {
// 			c.HTML(http.StatusInternalServerError, "orders.html", gin.H{"error": "Failed to fetch orders"})
// 			return
// 		}
// 		c.HTML(http.StatusOK, "orders.html", gin.H{"orders": orders})
// 	})

// 	// トレード履歴表示
// 	r.GET("/trades", func(c *gin.Context) {
// 		trades, err := app.FetchTrades(db.DB)
// 		if err != nil {
// 			c.HTML(http.StatusInternalServerError, "trades.html", gin.H{"error": "Failed to fetch trades"})
// 			return
// 		}
// 		c.HTML(http.StatusOK, "trades.html", gin.H{"trades": trades})
// 	})

// 	r.Run(":8081")
// }
