package update_balance

import (
	"context"

	"github.com.br/devfullcycle/fc-ms-wallet/internal/gateway"
	"github.com.br/devfullcycle/fc-ms-wallet/pkg/uow"
)

type UpdateBalanceInputDTO struct {
	AccountIDFrom        string  `json:"account_id_from"`
	AccountIDTo          string  `json:"account_id_to"`
	BalanceAccountIDFrom float64 `json:"balance_account_id_from"`
	BalanceAccountIDTo   float64 `json:"balance_account_id_to"`
}

type UpdateBalanceUseCase struct {
	Uow uow.UowInterface
}

func NewUpdateBalanceUseCase(Uow uow.UowInterface) *UpdateBalanceUseCase {
	return &UpdateBalanceUseCase{
		Uow: Uow,
	}
}

func (uc *UpdateBalanceUseCase) Execute(ctx context.Context, input UpdateBalanceInputDTO) error {
	err := uc.Uow.Do(ctx, func(_ *uow.Uow) error {
		accountRepository := uc.getAccountRepository(ctx)
		accountFrom, err := accountRepository.FindByID(input.AccountIDFrom)
		if err != nil {
			return err
		}

		accountTo, err := accountRepository.FindByID(input.AccountIDTo)
		if err != nil {
			return err
		}

		accountFrom.Balance = input.BalanceAccountIDFrom
		accountTo.Balance = input.BalanceAccountIDTo

		err = accountRepository.UpdateBalance(accountFrom)
		if err != nil {
			return err
		}
		err = accountRepository.UpdateBalance(accountTo)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func (uc *UpdateBalanceUseCase) getAccountRepository(ctx context.Context) gateway.AccountGateway {
	repo, err := uc.Uow.GetRepository(ctx, "AccountDB")
	if err != nil {
		panic(err)
	}
	return repo.(gateway.AccountGateway)
}
