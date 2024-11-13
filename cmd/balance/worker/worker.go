package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com.br/devfullcycle/fc-ms-wallet/internal/database"
	"github.com.br/devfullcycle/fc-ms-wallet/internal/usecase/create_transaction"
	"github.com.br/devfullcycle/fc-ms-wallet/internal/usecase/update_balance"
	"github.com.br/devfullcycle/fc-ms-wallet/pkg/kafka"
	"github.com.br/devfullcycle/fc-ms-wallet/pkg/uow"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"

	_ "github.com/go-sql-driver/mysql"
)

type Worker struct {
	Consumer             *kafka.Consumer
	UpdateBalanceUseCase *update_balance.UpdateBalanceUseCase
	Context              context.Context
}

func NewWorker(consumer *kafka.Consumer, updateBalanceUseCase *update_balance.UpdateBalanceUseCase, ctx context.Context) *Worker {
	return &Worker{
		Consumer:             consumer,
		UpdateBalanceUseCase: updateBalanceUseCase,
		Context:              ctx,
	}
}

func (w *Worker) Run() {
	msgChan := make(chan *ckafka.Message)
	go w.Consumer.Consume(msgChan)

	for msg := range msgChan {
		var input update_balance.UpdateBalanceInputDTO
		var message struct {
			Payload create_transaction.BalanceUpdatedOutputDTO `json:"Payload"`
		}
		err := json.Unmarshal(msg.Value, &message)
		if err != nil {
			log.Println(err)
			continue
		}

		input.AccountIDFrom = message.Payload.AccountIDFrom
		input.AccountIDTo = message.Payload.AccountIDTo
		input.BalanceAccountIDFrom = message.Payload.BalanceAccountIDFrom
		input.BalanceAccountIDTo = message.Payload.BalanceAccountIDTo

		err = w.UpdateBalanceUseCase.Execute(w.Context, input)
		if err != nil {
			log.Println(err)
			continue
		}

		log.Println("Balance updated")
	}
}

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", "root", "root", "mysql-balance", "3306", "balance"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	configMap := ckafka.ConfigMap{
		"bootstrap.servers": "kafka:29092",
		"group.id":          "wallet",
	}

	ctx := context.Background()
	uow := uow.NewUow(ctx, db)

	uow.Register("AccountDB", func(tx *sql.Tx) interface{} {
		return database.NewAccountDB(db)
	})

	updateBalanceUseCase := update_balance.NewUpdateBalanceUseCase(uow)
	consumer := kafka.NewConsumer(&configMap, []string{"balances"})
	worker := NewWorker(consumer, updateBalanceUseCase, ctx)

	worker.Run()
}
