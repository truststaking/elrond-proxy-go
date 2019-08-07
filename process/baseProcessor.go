package process

import (
	"bytes"
	"errors"
	"net/http"
	"sync"

	"github.com/ElrondNetwork/elrond-go/core/logger"
	"github.com/ElrondNetwork/elrond-go/data/state"
	"github.com/ElrondNetwork/elrond-go/sharding"
	"github.com/ElrondNetwork/elrond-proxy-go/config"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/gin-gonic/gin/json"
)

var log = logger.DefaultLogger()

// NodeStatusPath represents the path to send a request in order to see if a node is online
const NodeStatusPath = "/node/status"

// BaseProcessor represents an implementation of CoreProcessor that helps
// processing requests
type BaseProcessor struct {
	addressConverter state.AddressConverter
	lastConfig       *config.Config
	mutState         sync.RWMutex
	shardCoordinator sharding.Coordinator
	observers        map[uint32][]*data.Observer

	httpClient *http.Client
}

// NewBaseProcessor creates a new instance of BaseProcessor struct
func NewBaseProcessor(addressConverter state.AddressConverter) (*BaseProcessor, error) {
	if addressConverter == nil {
		return nil, ErrNilAddressConverter
	}

	return &BaseProcessor{
		observers:        make(map[uint32][]*data.Observer),
		httpClient:       http.DefaultClient,
		addressConverter: addressConverter,
	}, nil
}

// ApplyConfig applies a config on a base processor
func (bp *BaseProcessor) ApplyConfig(cfg *config.Config) error {
	if cfg == nil {
		return ErrNilConfig
	}
	if len(cfg.Observers) == 0 {
		return ErrEmptyObserversList
	}

	newObservers := make(map[uint32][]*data.Observer)
	maxShardId := uint32(0)
	for _, observer := range cfg.Observers {
		shardId := observer.ShardId
		if maxShardId < shardId {
			maxShardId = shardId
		}

		newObservers[shardId] = append(newObservers[shardId], observer)
	}

	newShardCoordinator, err := sharding.NewMultiShardCoordinator(maxShardId+1, 0)
	if err != nil {
		return err
	}

	bp.mutState.Lock()
	bp.shardCoordinator = newShardCoordinator
	bp.observers = newObservers
	bp.mutState.Unlock()

	return nil
}

// GetObservers returns the registered observers on a shard
func (bp *BaseProcessor) GetObservers(shardId uint32) ([]*data.Observer, error) {
	bp.mutState.RLock()
	defer bp.mutState.RUnlock()

	observers := bp.observers[shardId]
	if len(observers) == 0 {
		return nil, ErrMissingObserver
	}

	return observers, nil
}

// GetFirstAvailableObserver will return the first observer which will be found online
func (bp *BaseProcessor) GetFirstAvailableObserver() (*data.Observer, error) {
	for _, observerByShard := range bp.observers {
		for _, observer := range observerByShard {
			statusResponse := &data.StatusResponse{}
			err := bp.CallGetRestEndPoint(observer.Address, NodeStatusPath, statusResponse)
			if err != nil {
				continue
			}
			if statusResponse.Running == true {
				return observer, nil
			}
		}
	}
	return nil, errors.New("no observer online")
}

// ComputeShardId computes the shard id in which the account resides
func (bp *BaseProcessor) ComputeShardId(addressBuff []byte) (uint32, error) {
	bp.mutState.RLock()
	defer bp.mutState.RUnlock()

	address, err := bp.addressConverter.CreateAddressFromPublicKeyBytes(addressBuff)
	if err != nil {
		return 0, err
	}

	return bp.shardCoordinator.ComputeId(address), nil
}

// CallGetRestEndPoint calls an external end point (sends a request on a node)
func (bp *BaseProcessor) CallGetRestEndPoint(
	address string,
	path string,
	value interface{},
) error {

	req, err := http.NewRequest("GET", address+path, nil)
	if err != nil {
		return err
	}

	userAgent := "Elrond Proxy / 1.0.0 <Requesting data from nodes>"
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", userAgent)

	resp, err := bp.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		errNotCritical := resp.Body.Close()
		log.LogIfError(errNotCritical)
	}()

	return json.NewDecoder(resp.Body).Decode(value)
}

// CallPostRestEndPoint calls an external end point (sends a request on a node)
func (bp *BaseProcessor) CallPostRestEndPoint(
	address string,
	path string,
	data interface{},
	response interface{},
) error {

	buff, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", address+path, bytes.NewReader(buff))
	if err != nil {
		return err
	}

	userAgent := "Elrond Proxy / 1.0.0 <Posting to nodes>"
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)

	resp, err := bp.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		errNotCritical := resp.Body.Close()
		log.LogIfError(errNotCritical)
	}()

	return json.NewDecoder(resp.Body).Decode(response)
}
