package database

import (
	"encoding/json"
	"fmt"

	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/elastic/go-elasticsearch/v7"
)

const (
	numTopTransactions           = 20
	numTransactionFromAMiniblock = 100
)

type elasticSearchConnector struct {
	client *elasticsearch.Client
}

// NewElasticSearchConnector create a new elastic search database reader object
func NewElasticSearchConnector(url, username, password string) (*elasticSearchConnector, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{url},
		Username:  username,
		Password:  password,
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("cannot create database reader %w", err)
	}

	return &elasticSearchConnector{
		client: client,
	}, nil
}

// GetTransactionsByAddress gets transactions TO or FROM the specified address
func (esc *elasticSearchConnector) GetTransactionsByAddress(address string) ([]data.DatabaseTransaction, error) {
	query := txsByAddrQuery(address)
	decodedBody, err := esc.doSearchRequest(query, "transactions", numTopTransactions)
	if err != nil {
		return nil, err
	}

	return convertObjectToTransactions(decodedBody)
}

func (esc *elasticSearchConnector) GetLatestBlockHeight() (uint64, error) {
	query := latestBlockQuery()
	decodedBody, err := esc.doSearchRequest(query, "blocks", 1)
	if err != nil {
		return 0, err
	}

	block, _, err := convertObjectToBlock(decodedBody)
	if err != nil {
		return 0, err
	}

	return block.Nonce, nil
}

// GetBlockByNonce -
func (esc *elasticSearchConnector) GetBlockByNonce(nonce uint64) (data.ApiBlock, error) {
	query := blockByNonceAndShardIDQuery(nonce, core.MetachainShardId)
	decodedBody, err := esc.doSearchRequest(query, "blocks", 1)
	if err != nil {
		return data.ApiBlock{}, err
	}

	metaBlock, metaBlockHash, err := convertObjectToBlock(decodedBody)
	if err != nil {
		return data.ApiBlock{}, err
	}

	txs, err := esc.getTxsByMiniblockHashes(metaBlock.MiniBlocksHashes)
	if err != nil {
		return data.ApiBlock{}, err
	}

	transactions, err := esc.getTxsByNotarizedBlockHashes(metaBlock.NotarizedBlocksHashes)
	if err != nil {
		return data.ApiBlock{}, err
	}

	txs = append(txs, transactions...)

	return data.ApiBlock{
		Nonce:        metaBlock.Nonce,
		Hash:         metaBlockHash,
		Transactions: txs,
	}, nil
}

func (esc *elasticSearchConnector) getTxsByNotarizedBlockHashes(hashes []string) ([]data.DatabaseTransaction, error) {
	txs := make([]data.DatabaseTransaction, 0)
	for _, hash := range hashes {
		query := blockByHashQuery(hash)
		decodedBody, err := esc.doSearchRequest(query, "blocks", 1)
		if err != nil {
			return nil, err
		}

		shardBlock, _, err := convertObjectToBlock(decodedBody)
		if err != nil {
			return nil, err
		}

		transactions, err := esc.getTxsByMiniblockHashes(shardBlock.MiniBlocksHashes)
		if err != nil {
			return nil, err
		}

		txs = append(txs, transactions...)
	}
	return txs, nil
}

func (esc *elasticSearchConnector) getTxsByMiniblockHashes(hashes []string) ([]data.DatabaseTransaction, error) {
	txs := make([]data.DatabaseTransaction, 0)
	for _, hash := range hashes {
		query := txsByMiniblockHashQuery(hash)
		decodedBody, err := esc.doSearchRequest(query, "transactions", numTransactionFromAMiniblock)
		if err != nil {
			return nil, err
		}

		transactions, err := convertObjectToTransactions(decodedBody)
		if err != nil {
			return nil, err
		}

		txs = append(txs, transactions...)
	}
	return txs, nil
}

func (esc *elasticSearchConnector) doSearchRequest(query object, index string, size int) (object, error) {
	buff, err := encodeQuery(query)
	if err != nil {
		return nil, err
	}

	res, err := esc.client.Search(
		esc.client.Search.WithIndex(index),
		esc.client.Search.WithSize(size),
		esc.client.Search.WithBody(&buff),
	)
	if err != nil {
		return nil, fmt.Errorf("cannot get data from database: %w", err)
	}

	defer func() {
		_ = res.Body.Close()
	}()
	if res.IsError() {
		return nil, fmt.Errorf("cannot get data from database: %v", res)
	}

	var decodedBody map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&decodedBody); err != nil {
		return nil, err
	}

	return decodedBody, nil
}

func (esc *elasticSearchConnector) IsInterfaceNil() bool {
	return esc == nil
}
