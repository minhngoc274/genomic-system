package jobs

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	golibcron "github.com/golibs-starter/golib-cron"
	"github.com/golibs-starter/golib/log"
	"github.com/minhngoc274/genomic-system/genomic-service/adapters/blockchain"
	"github.com/minhngoc274/genomic-system/genomic-service/adapters/repositories"
	"github.com/minhngoc274/genomic-system/genomic-service/adapters/tee"
	"github.com/minhngoc274/genomic-system/genomic-service/config"
	"math/big"
	"strings"
)

const (
	uploadDataEventABI = `[{
		"anonymous":false,
		"inputs":[
			{"indexed":false,"internalType": "string","name":"docId","type":"string"},
			{"indexed":false,"internalType": "uint256","name":"sessionId","type":"uint256"}
		],
		"name":"UploadData",
		"type":"event"
	}]`
)

// FilterLogsJob is responsible for fetching the latest block and filtering logs for the UploadData event
type FilterLogsJob struct {
	props                 *config.ChainProperties
	client                *ethclient.Client
	abi                   abi.ABI
	geneticDataRepository *repositories.GeneticDataRepository
	blockchainService     *blockchain.BlockchainService
	teeService            *tee.TEEService
}

// NewFilterLogsJob ...
func NewFilterLogsJob(
	props *config.ChainProperties,
	geneticDataRepository *repositories.GeneticDataRepository,
	blockchainService *blockchain.BlockchainService,
	teeService *tee.TEEService,
) golibcron.Job {
	parsedABI, err := abi.JSON(strings.NewReader(uploadDataEventABI))
	if err != nil {
		log.Fatalf("Failed to parse contract ABI: %v", err)
	}

	client, err := ethclient.Dial(props.RPC)
	if err != nil {
		log.Fatalf("Failed to connect to blockchain: %v", err)
	}

	return &FilterLogsJob{
		props:                 props,
		client:                client,
		abi:                   parsedABI,
		geneticDataRepository: geneticDataRepository,
		blockchainService:     blockchainService,
		teeService:            teeService,
	}
}

// Run job handler
func (j *FilterLogsJob) Run(ctx context.Context) {
	log.Infof("[FilterLogsJob] job start")

	if err := j.fetchAndProcessLogs(ctx); err != nil {
		log.Errorf("[FilterLogsJob] Error processing logs %v", err)
	}

	log.Info("[FilterLogsJob] job stop")
}

// fetchAndProcessLogs fetches logs for the given block and processes UploadData events
func (j *FilterLogsJob) fetchAndProcessLogs(ctx context.Context) error {
	blockNumber, err := j.client.BlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("[FilterLogsJob] Failed to fetch latest block number: %v", err)
	}

	log.Infof("[FilterLogsJob] start filter from block %d to block %d", blockNumber-20, blockNumber)
	fromBlock := big.NewInt(int64(blockNumber - 20))
	query := ethereum.FilterQuery{
		FromBlock: fromBlock,
		ToBlock:   big.NewInt(int64(blockNumber)),
		Addresses: []common.Address{common.HexToAddress(j.props.ControllerAddress)},
	}

	logs, err := j.client.FilterLogs(ctx, query)
	if err != nil {
		return err
	}

	for _, vLog := range logs {
		events, err := j.abi.Unpack("UploadData", vLog.Data)
		if err != nil {
			log.Errorf("[FilterLogsJob] Error unpacking event data: %v\n", err)
			continue
		}

		if len(events) != 2 {
			log.Errorf("[FilterLogsJob] Unexpected number of parameters in event data")
			continue
		}

		docID, ok1 := events[0].(string)
		sessionID, ok2 := events[1].(*big.Int)

		if !ok1 || !ok2 {
			log.Errorf("[FilterLogsJob] Failed to cast event data to expected types")
			continue
		}

		log.Info("[FilterLogsJob] Decoded UploadData Event: DocID=%s, SessionID=%d\n", docID, sessionID)

		if err := j.processEvent(ctx, docID, sessionID); err != nil {
			log.Errorf("[FilterLogsJob] fail process event: %v", err)
			continue
		}
	}

	return nil
}

// processEvent process UploadData events
func (j *FilterLogsJob) processEvent(ctx context.Context, docID string, sessionID *big.Int) error {
	geneticData, exists := j.geneticDataRepository.Retrieve(docID)
	if !exists {
		return fmt.Errorf("no genetic data found for DocID=%s", docID)
	}

	if geneticData.IsConfirmed == true {
		log.Infof("[FilterLogsJob] docs %s already been processed", geneticData.FileID)
		return nil
	}

	riskScore, err := j.teeService.ProcessEncryptedData(geneticData.EncryptedData)
	if err != nil {
		return err
	}

	if err := j.blockchainService.Confirm(sessionID, big.NewInt(riskScore), docID, string(geneticData.DataHash), ""); err != nil {
		return err
	}

	geneticData.IsConfirmed = true
	if err := j.geneticDataRepository.Update(geneticData); err != nil {
		return err
	}

	return nil
}
