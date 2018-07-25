package pbft

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/opentracing/opentracing-go"

	"github.com/yejiayu/go-cita/common/merkle"
	"github.com/yejiayu/go-cita/database"
	"github.com/yejiayu/go-cita/database/block"
	"github.com/yejiayu/go-cita/log"
	"github.com/yejiayu/go-cita/types"
)

type Interface interface {
	Run() error
}

func New(
	dbFatory database.Factory,
	authClient types.AuthClient,
	chainClient types.ChainClient,
) Interface {
	return &pbft{
		blockDB:     dbFatory.BlockDB(),
		authClient:  authClient,
		chainClient: chainClient,
	}
}

type pbft struct {
	blockDB     block.Interface
	authClient  types.AuthClient
	chainClient types.ChainClient

	mu      sync.RWMutex
	latestH *types.BlockHeader
}

func (p *pbft) Run() error {
	h, err := p.blockDB.GetHeaderByLatest(context.Background())
	if err != nil {
		return err
	}

	p.latestH = h

	errCh := make(chan error)
	go func() {
		for {
			if err := p.consensus(); err != nil {
				log.Error(err)
			}

			p.updateHeader()
			time.Sleep(3 * time.Second)
		}
	}()

	return <-errCh
}

func (p *pbft) consensus() error {
	span := opentracing.GlobalTracer().StartSpan("consensus-pbft")
	defer span.Finish()

	p.mu.RLock()
	latestH := p.latestH
	p.mu.RUnlock()

	res, err := p.authClient.PackTransactions(
		opentracing.ContextWithSpan(context.Background(), span),
		&types.PackTransactionsReq{
			Height: latestH.GetHeight(),
		})
	if err != nil {
		return err
	}

	newBlock := &types.Block{
		Header: &types.BlockHeader{
			Prevhash:  latestH.GetPrevhash(),
			Timestamp: uint64(time.Now().Unix()),
			Height:    latestH.GetHeight() + 1,
		},
		Body: &types.BlockBody{
			TxHashes: res.GetTxHashes(),
		},
	}

	if len(res.GetTxHashes()) > 0 {
		txTree, err := merkle.New(res.GetTxHashes())
		if err != nil {
			return err
		}

		newBlock.Header.TransactionsRoot = txTree.MerkleRoot()
	}

	_, err = p.chainClient.NewBlock(
		opentracing.ContextWithSpan(context.Background(), span),
		&types.NewBlockReq{
			Block: newBlock,
		},
	)
	return err
}

func (p *pbft) updateHeader() error {
	h, err := p.blockDB.GetHeaderByLatest(context.Background())
	if err != nil {
		return err
	}

	p.mu.Lock()
	p.latestH = h
	p.mu.Unlock()

	log.Infof("update header, current height is %d", h.GetHeight())
	body, err := p.blockDB.GetBodyByHeight(context.Background(), h.GetHeight())
	if err != nil {
		return err
	}

	hashes := []string{}
	for _, hash := range body.GetTxHashes() {
		hashes = append(hashes, common.ToHex(hash))
	}

	log.Infof("tx hashes %s", strings.Join(hashes, ","))
	return nil
}
