package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"time"
	"trading_app_cli/db"
	"trading_app_cli/engine"
	"trading_app_gui/app"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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
		// クッキーからユーザー名を取得
		username, err := c.Cookie("trading_session")
		if err != nil {
			// クッキーがない場合はログインページにリダイレクト
			c.Redirect(http.StatusSeeOther, "/login")
			return
		}
		// クッキーの値をログ出力
		log.Printf("[DEBUG] クッキーの値 (username): %s", username)
		c.HTML(http.StatusOK, "dashboard.html", gin.H{"username": username})
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

	// ログイン処理
	r.POST("/login", func(c *gin.Context) {
		var loginData struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		// JSON バインド
		if err := c.ShouldBindJSON(&loginData); err != nil {
			log.Printf("[ERROR] JSON バインドエラー: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		// バインドされたデータをログに出力
		log.Printf("[DEBUG] LoginData from request: Username=%s, Password=%s", loginData.Username, loginData.Password)

		// ユーザー認証
		isValid, err := app.AuthenticateUser(loginData.Username, loginData.Password)
		if err != nil {
			log.Printf("[ERROR] ユーザー認証中のエラー: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		if !isValid {
			log.Printf("[DEBUG] ユーザー認証失敗: Username=%s", loginData.Username)
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid credentials"})
			return
		}

		// クッキーを設定
		log.Printf("[DEBUG] クッキーに設定する値: Username=%s", loginData.Username)
		c.SetCookie("trading_session", loginData.Username, 3600, "/", "localhost", false, true)

		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	// ログインページ
	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})
	// ログアウト処理
	r.POST("/logout", func(c *gin.Context) {
		c.SetCookie("trading_session", "", -1, "/", "localhost", false, true)
		c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
	})
	// ユーザー登録処理
	r.POST("/register", func(c *gin.Context) {
		var userData struct {
			Username string `json:"username"` // JSONのキー名を修正
			Password string `json:"password"`
		}

		// JSONバインド
		if err := c.ShouldBindJSON(&userData); err != nil {
			log.Printf("[DEBUG] JSONバインドエラー: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid input"})
			return
		}

		log.Printf("[DEBUG] 受信したデータ: Username=%s, Password=%s", userData.Username, userData.Password)

		// パスワードのハッシュ化
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("[DEBUG] パスワードハッシュエラー: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Password hashing failed"})
			return
		}

		// ユーザー名の重複チェック
		var existingUser string
		err = db.DB.QueryRow("SELECT name FROM users WHERE name = $1", userData.Username).Scan(&existingUser)
		if err == nil {
			log.Printf("[DEBUG] ユーザー名が既に存在します: %s", userData.Username)
			c.JSON(http.StatusConflict, gin.H{"success": false, "message": "Username already taken"})
			return
		} else if err != sql.ErrNoRows {
			log.Printf("[DEBUG] ユーザー名チェックエラー: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Database error"})
			return
		}

		// 新規ユーザー登録
		_, err = db.DB.Exec("INSERT INTO users (name, password, balance) VALUES ($1, $2, $3)", userData.Username, string(hashedPassword), 0)
		if err != nil {
			log.Printf("[DEBUG] ユーザー登録エラー: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "User registration failed"})
			return
		}

		log.Printf("[DEBUG] ユーザー登録成功: %s", userData.Username)
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "User registered successfully"})
	})

	r.Run(":8081")
}
