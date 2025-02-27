package evmc

import (
	"context"
	"errors"

	"github.com/bbaktaeho/evmc/evmctypes"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/shopspring/decimal"
)

// TODO: get uncle block
// TODO: batch call
// TODO: describe custom functions

type ethNamespace struct {
	c caller
}

func (e *ethNamespace) ChainID() (uint64, error) {
	return e.chainID(context.Background())
}

func (e *ethNamespace) ChainIDWithContext(ctx context.Context) (uint64, error) {
	return e.chainID(ctx)
}

func (e *ethNamespace) chainID(ctx context.Context) (uint64, error) {
	result := new(string)
	if err := e.c.call(ctx, result, ethChainID); err != nil {
		return 0, err
	}
	id, err := hexutil.DecodeUint64(*result)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (e *ethNamespace) GetStorageAt(address, position string, blockAndTag BlockAndTag) (string, error) {
	return e.getStorageAt(context.Background(), address, position, blockAndTag)
}

func (e *ethNamespace) GetStorageAtWithContext(
	ctx context.Context,
	address string,
	position string,
	blockAndTag BlockAndTag,
) (string, error) {
	return e.getStorageAt(ctx, address, position, blockAndTag)
}

func (e *ethNamespace) getStorageAt(
	ctx context.Context,
	address string,
	position string,
	numOrTag BlockAndTag,
) (string, error) {
	result := new(string)
	parsedBT := parseBlockAndTag(numOrTag)
	if err := e.c.call(
		ctx,
		result,
		ethGetStorageAt,
		address,
		position,
		parsedBT,
	); err != nil {
		return "", err
	}
	return *result, nil
}

func (e *ethNamespace) GetBlockNumber() (uint64, error) {
	return e.getBlockNumber(context.Background())
}

func (e *ethNamespace) GetBlockNumberWithContext(ctx context.Context) (uint64, error) {
	return e.getBlockNumber(ctx)
}

func (e *ethNamespace) getBlockNumber(ctx context.Context) (uint64, error) {
	result := new(string)
	if err := e.c.call(ctx, result, ethBlockNumber); err != nil {
		return 0, err
	}
	return hexutil.MustDecodeUint64(*result), nil
}

func (e *ethNamespace) GetCode(address string, blockAndTag BlockAndTag) (string, error) {
	return e.getCode(context.Background(), address, blockAndTag)
}

func (e *ethNamespace) GetCodeWithContext(ctx context.Context, address string, blockAndTag BlockAndTag) (string, error) {
	return e.getCode(ctx, address, blockAndTag)
}

func (e *ethNamespace) getCode(
	ctx context.Context,
	address string,
	blockAndTag BlockAndTag,
) (string, error) {
	result := new(string)
	parsedBT := parseBlockAndTag(blockAndTag)
	if err := e.c.call(ctx, result, ethGetCode, address, parsedBT); err != nil {
		return "", err
	}
	return *result, nil
}

func (e *ethNamespace) GetBlockByTag(tag BlockAndTag) (*evmctypes.Block[[]string], error) {
	return e.getBlockByTag(context.Background(), tag)
}

func (e *ethNamespace) GetBlockByTagWithContext(ctx context.Context, tag BlockAndTag) (*evmctypes.Block[[]string], error) {
	return e.getBlockByTag(ctx, tag)
}

func (e *ethNamespace) getBlockByTag(ctx context.Context, tag BlockAndTag) (*evmctypes.Block[[]string], error) {
	block := new(evmctypes.Block[[]string])
	if err := e.getBlockByNumber(ctx, block, tag, false); err != nil {
		return nil, err
	}
	return block, nil
}

func (e *ethNamespace) GetBlockByTagIncTx(tag BlockAndTag) (*evmctypes.Block[[]*evmctypes.Transaction], error) {
	return e.getBlockByTagIncTx(context.Background(), tag)
}

func (e *ethNamespace) GetBlockByTagIncTxWithContext(ctx context.Context, tag BlockAndTag) (*evmctypes.Block[[]*evmctypes.Transaction], error) {
	return e.getBlockByTagIncTx(ctx, tag)
}

func (e *ethNamespace) getBlockByTagIncTx(ctx context.Context, tag BlockAndTag) (*evmctypes.Block[[]*evmctypes.Transaction], error) {
	block := new(evmctypes.Block[[]*evmctypes.Transaction])
	if err := e.getBlockByNumber(ctx, block, tag, true); err != nil {
		return nil, err
	}
	return block, nil
}

func (e *ethNamespace) GetBlockByNumber(number uint64) (*evmctypes.Block[[]string], error) {
	return e.getBlock(context.Background(), number)
}

func (e *ethNamespace) GetBlockByNumberWithContext(ctx context.Context, number uint64) (*evmctypes.Block[[]string], error) {
	return e.getBlock(ctx, number)
}

func (e *ethNamespace) getBlock(ctx context.Context, number uint64) (*evmctypes.Block[[]string], error) {
	result := new(evmctypes.Block[[]string])
	if err := e.getBlockByNumber(ctx, result, FormatNumber(number), false); err != nil {
		return nil, err
	}
	return result, nil
}

func (e *ethNamespace) GetBlockByNumberIncTx(number uint64) (*evmctypes.Block[[]*evmctypes.Transaction], error) {
	return e.getBlockByNumberIncTx(context.Background(), number)
}

func (e *ethNamespace) GetBlockByNumberIncTxWithContext(
	ctx context.Context,
	number uint64,
) (*evmctypes.Block[[]*evmctypes.Transaction], error) {
	return e.getBlockByNumberIncTx(ctx, number)
}

func (e *ethNamespace) getBlockByNumberIncTx(ctx context.Context, number uint64) (*evmctypes.Block[[]*evmctypes.Transaction], error) {
	block := new(evmctypes.Block[[]*evmctypes.Transaction])
	if err := e.getBlockByNumber(ctx, block, FormatNumber(number), true); err != nil {
		return nil, err
	}
	return block, nil
}

func (e *ethNamespace) getBlockByNumber(
	ctx context.Context,
	result interface{},
	number BlockAndTag,
	incTx bool,
) error {
	if number == Pending {
		return ErrPendingBlockNotSupported
	}
	parsedBT := parseBlockAndTag(number)
	params := []interface{}{parsedBT, incTx}
	if err := e.c.call(ctx, result, ethGetBlockByNumber, params...); err != nil {
		return err
	}
	return nil
}

func (e *ethNamespace) GetBlockByHash(hash string) (*evmctypes.Block[[]string], error) {
	block := new(evmctypes.Block[[]string])
	if err := e.getBlockByHash(context.Background(), block, hash, false); err != nil {
		return nil, err
	}
	return block, nil
}

func (e *ethNamespace) GetBlockByHashWithContext(ctx context.Context, hash string) (*evmctypes.Block[[]string], error) {
	block := new(evmctypes.Block[[]string])
	if err := e.getBlockByHash(ctx, block, hash, false); err != nil {
		return nil, err
	}
	return block, nil
}

func (e *ethNamespace) GetBlockByHashIncTx(hash string) (*evmctypes.Block[[]*evmctypes.Transaction], error) {
	return e.getBlockByHashIncTx(context.Background(), hash)
}

func (e *ethNamespace) GetBlockByHashIncTxWithContext(
	ctx context.Context,
	hash string,
) (*evmctypes.Block[[]*evmctypes.Transaction], error) {
	return e.getBlockByHashIncTx(ctx, hash)
}

func (e *ethNamespace) getBlockByHashIncTx(ctx context.Context, hash string) (*evmctypes.Block[[]*evmctypes.Transaction], error) {
	block := new(evmctypes.Block[[]*evmctypes.Transaction])
	if err := e.getBlockByHash(ctx, block, hash, true); err != nil {
		return nil, err
	}
	return block, nil
}

func (e *ethNamespace) getBlockByHash(ctx context.Context, result interface{}, hash string, incTx bool) error {
	params := []interface{}{hash, incTx}
	if err := e.c.call(ctx, result, ethGetBlockByHash, params...); err != nil {
		return err
	}
	return nil
}

func (e *ethNamespace) GetTransaction(hash string) (*evmctypes.Transaction, error) {
	return e.getTransaction(context.Background(), hash)
}

func (e *ethNamespace) GetTransactionWithContext(ctx context.Context, hash string) (*evmctypes.Transaction, error) {
	return e.getTransaction(ctx, hash)
}

func (e *ethNamespace) getTransaction(ctx context.Context, hash string) (*evmctypes.Transaction, error) {
	tx := new(evmctypes.Transaction)
	if err := e.c.call(ctx, tx, ethGetTransaction, hash); err != nil {
		return nil, err
	}
	return tx, nil
}

func (e *ethNamespace) GetTransactionReceipt(hash string) (*evmctypes.Receipt, error) {
	return e.getTransactionReceipt(context.Background(), hash)
}

func (e *ethNamespace) GetTransactionReceiptWithContext(ctx context.Context, hash string) (*evmctypes.Receipt, error) {
	return e.getTransactionReceipt(ctx, hash)
}

func (e *ethNamespace) getTransactionReceipt(ctx context.Context, hash string) (*evmctypes.Receipt, error) {
	receipt := new(evmctypes.Receipt)
	if err := e.c.call(ctx, receipt, ethGetReceipt, hash); err != nil {
		return nil, err
	}
	return receipt, nil
}

func (e *ethNamespace) GetBalance(address string, blockAndTag BlockAndTag) (decimal.Decimal, error) {
	return e.getBalance(context.Background(), address, blockAndTag)
}

func (e *ethNamespace) GetBalanceWithContext(ctx context.Context, address string, blockAndTag BlockAndTag) (decimal.Decimal, error) {
	return e.getBalance(ctx, address, blockAndTag)
}

func (e *ethNamespace) getBalance(ctx context.Context, address string, blockAndTag BlockAndTag) (decimal.Decimal, error) {
	result := new(string)
	parsedBT := parseBlockAndTag(blockAndTag)
	if err := e.c.call(ctx, result, ethGetBalance, address, parsedBT); err != nil {
		return decimal.Zero, err
	}
	if *result == "" {
		*result = "0x0"
	}
	return decimal.NewFromBigInt(hexutil.MustDecodeBig(*result), 0), nil
}

func (e *ethNamespace) GetLogs(filter *evmctypes.LogFilter) ([]*evmctypes.Log, error) {
	return e.getLogs(context.Background(), filter)
}

func (e *ethNamespace) GetLogsWithContext(ctx context.Context, filter *evmctypes.LogFilter) ([]*evmctypes.Log, error) {
	return e.getLogs(ctx, filter)
}

func (e *ethNamespace) GetLogsByBlockNumber(number uint64) ([]*evmctypes.Log, error) {
	return e.getLogsByBlockNumber(context.Background(), number)
}

func (e *ethNamespace) GetLogsByBlockNumberWithContext(ctx context.Context, number uint64) ([]*evmctypes.Log, error) {
	return e.getLogsByBlockNumber(ctx, number)
}

func (e *ethNamespace) getLogsByBlockNumber(ctx context.Context, number uint64) ([]*evmctypes.Log, error) {
	filter := &evmctypes.LogFilter{
		FromBlock: &number,
		ToBlock:   &number,
	}
	return e.getLogs(ctx, filter)
}

func (e *ethNamespace) GetLogsByBlockHash(hash string) ([]*evmctypes.Log, error) {
	return e.getLogsByBlockHash(context.Background(), hash)
}

func (e *ethNamespace) GetLogsByBlockHashWithContext(ctx context.Context, hash string) ([]*evmctypes.Log, error) {
	return e.getLogsByBlockHash(ctx, hash)
}

func (e *ethNamespace) getLogsByBlockHash(ctx context.Context, hash string) ([]*evmctypes.Log, error) {
	filter := &evmctypes.LogFilter{
		BlockHash: &hash,
	}
	return e.getLogs(ctx, filter)
}

func (e *ethNamespace) getLogs(ctx context.Context, filter *evmctypes.LogFilter) ([]*evmctypes.Log, error) {
	logs := new([]*evmctypes.Log)
	params := make(map[string]interface{})
	if filter.BlockHash != nil {
		params["blockHash"] = *filter.BlockHash
	} else {
		params["fromBlock"] = hexutil.EncodeUint64(*filter.FromBlock)
		params["toBlock"] = hexutil.EncodeUint64(*filter.ToBlock)
	}
	if filter.Address != nil {
		params["address"] = *filter.Address
	}
	if filter.Topics != nil {
		params["topics"] = filter.Topics
	}
	if err := e.c.call(ctx, logs, ethGetLogs, filter); err != nil {
		return nil, err
	}
	return *logs, nil
}

func (e *ethNamespace) GetTransactionCount(address string, blockAndTag BlockAndTag) (uint64, error) {
	return e.getTransactionCount(context.Background(), address, blockAndTag)
}

func (e *ethNamespace) GetTransactionCountWithContext(
	ctx context.Context,
	address string,
	blockAndTag BlockAndTag,
) (uint64, error) {
	return e.getTransactionCount(ctx, address, blockAndTag)
}

func (e *ethNamespace) getTransactionCount(ctx context.Context, address string, blockAndTag BlockAndTag) (uint64, error) {
	result := new(string)
	parsedBT := parseBlockAndTag(blockAndTag)
	if err := e.c.call(ctx, result, ethGetTransactionCount, address, parsedBT); err != nil {
		return 0, err
	}
	return hexutil.MustDecodeUint64(*result), nil
}

func (e *ethNamespace) GetBlockReceipts(number uint64) ([]*evmctypes.Receipt, error) {
	return e.getBlockReceipts(context.Background(), number)
}

func (e *ethNamespace) GetBlockReceiptsWithContext(ctx context.Context, number uint64) ([]*evmctypes.Receipt, error) {
	return e.getBlockReceipts(ctx, number)
}

func (e *ethNamespace) getBlockReceipts(ctx context.Context, number uint64) ([]*evmctypes.Receipt, error) {
	var (
		result        = new([]*evmctypes.Receipt)
		method        = ethGetBlockReceipts
		clientName, _ = e.c.NodeClient()
	)
	if ClientName(clientName) == Bor {
		method = ethGetTransactionReceiptsByBlock
	}
	if err := e.c.call(ctx, result, method, hexutil.EncodeUint64(number)); err != nil {
		return nil, err
	}
	return *result, nil
}

func (e *ethNamespace) GasPrice() (decimal.Decimal, error) {
	return e.gasPrice(context.Background())
}

func (e *ethNamespace) GasPriceWithContext(ctx context.Context) (decimal.Decimal, error) {
	return e.gasPrice(ctx)
}

func (e *ethNamespace) gasPrice(ctx context.Context) (decimal.Decimal, error) {
	result := new(string)
	if err := e.c.call(ctx, result, ethGasPrice); err != nil {
		return decimal.Zero, err
	}
	return decimal.NewFromBigInt(hexutil.MustDecodeBig(*result), 0), nil
}

func (e *ethNamespace) MaxPriorityFeePerGas() (decimal.Decimal, error) {
	return e.maxPriorityFeePerGas(context.Background())
}

func (e *ethNamespace) MaxPriorityFeePerGasWithContext(ctx context.Context) (decimal.Decimal, error) {
	return e.maxPriorityFeePerGas(ctx)
}

func (e *ethNamespace) maxPriorityFeePerGas(ctx context.Context) (decimal.Decimal, error) {
	result := new(string)
	if err := e.c.call(ctx, result, ethMaxPriorityFeePerGas); err != nil {
		return decimal.Zero, err
	}
	return decimal.NewFromBigInt(hexutil.MustDecodeBig(*result), 0), nil
}

func (e *ethNamespace) SendTransaction(sendingTx *SendingTx, wallet *Wallet) (string, error) {
	return e.sendTransaction(context.Background(), sendingTx, wallet)
}

func (e *ethNamespace) SendTransactionWithContext(
	ctx context.Context,
	sendingTx *SendingTx,
	wallet *Wallet,
) (string, error) {
	return e.sendTransaction(ctx, sendingTx, wallet)
}

func (e *ethNamespace) sendTransaction(
	ctx context.Context,
	sendingTx *SendingTx,
	wallet *Wallet,
) (string, error) {
	hash, rawTx, err := wallet.SignTx(sendingTx, e.c.ChainID())
	if err != nil {
		return "", err
	}
	txHash, err := e.sendRawTransaction(ctx, rawTx)
	if err != nil {
		return "", err
	}
	if hash != txHash {
		return "", errors.New("transaction hash mismatch")
	}
	return txHash, nil
}

func (e *ethNamespace) SendRawTransaction(rawTx string) (string, error) {
	return e.sendRawTransaction(context.Background(), rawTx)
}

func (e *ethNamespace) SendRawTransactionWithContext(ctx context.Context, rawTx string) (string, error) {
	return e.sendRawTransaction(ctx, rawTx)
}

func (e *ethNamespace) sendRawTransaction(ctx context.Context, rawTx string) (string, error) {
	result := new(string)
	if err := e.c.call(ctx, result, ethSendRawTransaction, rawTx); err != nil {
		return "", err
	}
	return *result, nil
}
