package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com.br/devfullcycle/fc-ms-wallet/internal/database"
	"github.com.br/devfullcycle/fc-ms-wallet/internal/usecase/get_balance"
	"github.com.br/devfullcycle/fc-ms-wallet/internal/web"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", "root", "root", "mysql-balance", "3306", "balance"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	accountDb := database.NewAccountDB(db)
	getBalanceUseCase := get_balance.NewGetBalanceUseCase(accountDb)

	balanceHander := web.NewWebBalanceHandler(getBalanceUseCase)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/balance/{accountID}", balanceHander.GetBalance)

	log.Println("Server running on port 3003")
	http.ListenAndServe(":3003", r)
}
