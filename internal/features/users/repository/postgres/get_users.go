package users_postgres_repository

import (
	"context"
	"fmt"

	"github.com/kiricenkokbravl5-beep/Golang-todoapp-/tree/infra/env-setup/internal/core/domain"
	core_postgres_pool "github.com/kiricenkokbravl5-beep/Golang-todoapp-/tree/infra/env-setup/internal/core/repository/postgres/pool"
)

func userDomainsFromModels(userModels []UserModel) []*domain.User {
	userDomains := make([]*domain.User, len(userModels))
	for i, userModel := range userModels {
		userDomains[i] = &domain.User{
			ID:          userModel.ID,
			Version:     userModel.Version,
			FullName:    userModel.FullName,
			PhoneNumber: userModel.PhoneNumber,
		}
	}
	return userDomains
}

func (r *UsersRepository) GetUsers(
	ctx context.Context,
	limit *int,
	offset *int,
) ([]*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `SELECT id, version, full_name, phone_number
		FROM todoapp.users
		ORDER BY id ASC
		LIMIT $1 
		OFFSET $2;
		`

	rows, err := r.pool.Query(
		ctx,
		query,
		limit,
		offset,
	)
	if err != nil {
		return nil, fmt.Errorf("select users: %w", err)
	}
	defer rows.Close()
	var userModels []UserModel
	for rows.Next() {
		var userModel UserModel

		err := rows.Scan(
			&userModel.ID,
			&userModel.Version,
			&userModel.FullName,
			&userModel.PhoneNumber,
		)
		if err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		userModels = append(userModels, userModel)
	}
	if err := core_postgres_pool.ErrNoRows; err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	userDomains := userDomainsFromModels(userModels)
	return userDomains, nil

}
