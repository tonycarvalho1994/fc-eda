package get_balance

import "github.com.br/devfullcycle/fc-ms-wallet/internal/gateway"

type GetBalanceUseCase struct {
	DB gateway.AccountGateway
}

type GetBalanceOutputDTO struct {
	Balance float64 `json:"balance"`
}

func NewGetBalanceUseCase(db gateway.AccountGateway) GetBalanceUseCase {
	return GetBalanceUseCase{DB: db}
}

func (g GetBalanceUseCase) Execute(accountID string) (*GetBalanceOutputDTO, error) {
	account, err := g.DB.FindByID(accountID)
	if err != nil {
		return nil, err
	}

	return &GetBalanceOutputDTO{Balance: account.Balance}, nil
}
