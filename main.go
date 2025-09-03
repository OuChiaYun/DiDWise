package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	// DB 連線字串
	dsn := "postgres://postgres:pass@localhost:8080/didwisesql?sslmode=disable"
	if env := os.Getenv("DATABASE_URL"); env != "" {
		dsn = env
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("connect error:", err)
	}
	defer db.Close()

	http.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		rows, err := db.QueryContext(ctx, "SELECT id, name FROM items ORDER BY id")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var items []Item
		for rows.Next() {
			var it Item
			if err := rows.Scan(&it.ID, &it.Name); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			items = append(items, it)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(items)
	})

	fmt.Println("server on :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
