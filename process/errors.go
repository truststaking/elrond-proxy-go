package process

import "errors"

// ErrMissingObserver signals that no observers have been provided for provided shard ID
var ErrMissingObserver = errors.New("missing observer")

// ErrSendingRequest signals that sending the request failed on all observers
var ErrSendingRequest = errors.New("sending request error")

// ErrNilShardCoordinator signals that a nil shard coordinator has been provided
var ErrNilShardCoordinator = errors.New("nil shard coordinator")

// ErrInvalidRequestTimeout signals that the provided number of seconds before timeout is invalid
var ErrInvalidRequestTimeout = errors.New("invalid duration until timeout for requests")

// ErrNilCoreProcessor signals that a nil core processor has been provided
var ErrNilCoreProcessor = errors.New("nil core processor")

// ErrNilPrivateKeysLoader signals that a nil private keys loader has been provided
var ErrNilPrivateKeysLoader = errors.New("nil private keys loader")

// ErrEmptyMapOfAccountsFromPem signals that an empty map of accounts was received
var ErrEmptyMapOfAccountsFromPem = errors.New("empty map of accounts read from the pem file")

// ErrInvalidEconomicsConfig signals that the provided economics config cannot be parsed
var ErrInvalidEconomicsConfig = errors.New("cannot parse economics config")

// ErrHeartbeatNotAvailable signals that the heartbeat status is not found
var ErrHeartbeatNotAvailable = errors.New("heartbeat status not found at any observer")

// ErrNilHeartbeatCacher signals that the provided heartbeat cacher is nil
var ErrNilHeartbeatCacher = errors.New("nil heartbeat cacher")

// ErrNilValidatorStatisticsCacher signals that the provided validator statistics cacher is nil
var ErrNilValidatorStatisticsCacher = errors.New("nil validator statistics cacher")

// ErrValidatorStatisticsNotAvailable signals that the validator statistics data is not found
var ErrValidatorStatisticsNotAvailable = errors.New("validator statistics data not found at any observer")

// ErrInvalidCacheValidityDuration signals that the given validity duration for cache data is invalid
var ErrInvalidCacheValidityDuration = errors.New("invalid cache validity duration")

// ErrNilDefaultFaucetValue signals that a nil default faucet value has been provided
var ErrNilDefaultFaucetValue = errors.New("nil default faucet value provided")

// ErrInvalidDefaultFaucetValue signals that the provided faucet value is not strictly positive
var ErrInvalidDefaultFaucetValue = errors.New("default faucet value is not strictly positive")

// ErrNilObserversProvider signals that a nil observers provider has been provided
var ErrNilObserversProvider = errors.New("the observers provider is nil")

// ErrInvalidShardId signals that a invalid shard id has been provided
var ErrInvalidShardId = errors.New("invalid shard id")

// ErrNilPubKeyConverter signals that a nil pub key converter has been provided
var ErrNilPubKeyConverter = errors.New("nil pub key converter provided")

// ErrNoValidTransactionToSend signals that no valid transaction were received
var ErrNoValidTransactionToSend = errors.New("no valid transaction to send")

// ErrCannotGetTransactionStatus -
var ErrCannotGetTransactionStatus = errors.New("cannot get transaction status")

// ErrNilDatabaseReader signals that a nil database reader was provided
var ErrNilDatabaseReader = errors.New("not valid database reader")
