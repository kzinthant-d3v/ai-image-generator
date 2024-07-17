package db

import (
	"context"
	"kzinthant-d3v/ai-image-generator/types"

	"github.com/google/uuid"
)

func CreateAccount(account *types.Account) error {
	_, err := Bun.NewInsert().Model(account).Exec(context.Background())
	return err
}

func GetAccountByID(userID uuid.UUID) (types.Account, error) {
	account := types.Account{}
	err := Bun.NewSelect().Model(&account).Where("user_id = ?", userID).Scan(context.Background())
	return account, err
}
