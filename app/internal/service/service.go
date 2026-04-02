package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"app/internal/model"
)

type BlockchainClient interface {
	ExecContract(ctx context.Context, value uint64) error
	CallContract(ctx context.Context) (uint64, error)
}

type StateRepository interface {
	SaveState(ctx context.Context, value uint64) error
	GetState(ctx context.Context) (model.ContractState, error)
}

type Service struct {
	blockchainClient BlockchainClient
	stateRepository  StateRepository
}

func NewService(blockchainClient BlockchainClient, stateRepository StateRepository) *Service {
	return &Service{
		blockchainClient: blockchainClient,
		stateRepository:  stateRepository,
	}
}

func (service *Service) SetValue(ctx context.Context, value uint64) error {
	return service.blockchainClient.ExecContract(ctx, value)
}

func (service *Service) GetValue(ctx context.Context) (uint64, error) {
	return service.blockchainClient.CallContract(ctx)
}

func (service *Service) SyncValue(ctx context.Context) (uint64, error) {
	blockchainValue, err := service.blockchainClient.CallContract(ctx)
	if err != nil {
		return 0, fmt.Errorf("get blockchain value: %w", err)
	}
	if err := service.stateRepository.SaveState(ctx, blockchainValue); err != nil {
		return 0, fmt.Errorf("save database value: %w", err)
	}
	return blockchainValue, nil
}

func (service *Service) CheckValue(ctx context.Context) (bool, uint64, uint64, error) {
	blockchainValue, err := service.blockchainClient.CallContract(ctx)
	if err != nil {
		return false, 0, 0, fmt.Errorf("get blockchain value: %w", err)
	}
	storedState, err := service.stateRepository.GetState(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, blockchainValue, 0, nil
		}
		return false, 0, 0, fmt.Errorf("get database value: %w", err)
	}
	return blockchainValue == storedState.Value, blockchainValue, storedState.Value, nil
}
