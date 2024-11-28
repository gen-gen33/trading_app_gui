package app

import (
	"database/sql"
	"trading_app_cli/db"

	"golang.org/x/crypto/bcrypt"
)

// AuthenticateUser checks if the username and password match a valid user in the database.
func AuthenticateUser(username, password string) (bool, error) {
	var hashedPassword string

	// データベースからハッシュ化されたパスワードを取得
	err := db.DB.QueryRow("SELECT password FROM users WHERE name = $1", username).Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			// ユーザーが存在しない場合
			return false, nil
		}
		// その他のエラー
		return false, err
	}

	// ハッシュ化されたパスワードと入力されたパスワードを比較
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		// パスワードが一致しない場合
		return false, nil
	}

	return true, nil
}
