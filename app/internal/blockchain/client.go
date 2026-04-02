package blockchain

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"strings"

	"app/contracts"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Client struct {
	client          *ethclient.Client
	contractAddress common.Address
	chainID         *big.Int
	private         string
}

func NewClient(ctx context.Context, rpcURL string, contractAddress string, privateKey string) (*Client, error) {
	if !common.IsHexAddress(contractAddress) {
		return nil, fmt.Errorf("invalid CONTRACT_ADDRESS format: %s", contractAddress)
	}
	address := common.HexToAddress(contractAddress)
	if address == (common.Address{}) {
		return nil, fmt.Errorf("invalid CONTRACT_ADDRESS: zero address is not allowed")
	}
	ethClient, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		return nil, fmt.Errorf("dial rpc: %w", err)
	}
	chainID, err := ethClient.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("chain id: %w", err)
	}
	return &Client{
		client:          ethClient,
		contractAddress: address,
		chainID:         chainID,
		private:         strings.TrimPrefix(privateKey, "0x"),
	}, nil
}

func (client *Client) ExecContract(ctx context.Context, value uint64) error {
	contractABI, err := abi.JSON(strings.NewReader(contracts.SimpleStorageABI))
	if err != nil {
		return fmt.Errorf("parse abi: %w", err)
	}
	slog.Info("querying chain id")
	boundContract := bind.NewBoundContract(
		client.contractAddress,
		contractABI,
		client.client,
		client.client,
		client.client,
	)
	priv, err := crypto.HexToECDSA(client.private)
	if err != nil {
		return fmt.Errorf("parse private key: %w", err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(priv, client.chainID)
	if err != nil {
		return fmt.Errorf("create transactor: %w", err)
	}
	auth.Context = ctx
	tx, err := boundContract.Transact(auth, "set", new(big.Int).SetUint64(value))
	if err != nil {
		return fmt.Errorf("transact: %w", err)
	}
	fmt.Println("waiting until transaction is mined", "tx", tx.Hash().Hex())
	receipt, err := bind.WaitMined(ctx, client.client, tx)
	if err != nil {
		return fmt.Errorf("wait mined: %w", err)
	}
	fmt.Printf("transaction mined: %v\n", receipt)
	return nil
}

func (client *Client) CallContract(ctx context.Context) (uint64, error) {
	contractABI, err := abi.JSON(strings.NewReader(contracts.SimpleStorageABI))
	if err != nil {
		return 0, fmt.Errorf("parse abi: %w", err)
	}
	caller := bind.CallOpts{
		Pending: false,
		Context: ctx,
	}
	boundContract := bind.NewBoundContract(
		client.contractAddress,
		contractABI,
		client.client,
		client.client,
		client.client,
	)
	var output []interface{}
	if err := boundContract.Call(&caller, &output, "get"); err != nil {
		return 0, fmt.Errorf("call contract: %w", err)
	}
	if len(output) == 0 {
		return 0, fmt.Errorf("empty output from get")
	}
	result, ok := output[0].(*big.Int)
	if !ok {
		return 0, fmt.Errorf("unexpected output type from get")
	}
	fmt.Println("Successfully called contract!", output)
	return result.Uint64(), nil
}

func (client *Client) Close() {
	client.client.Close()
}
