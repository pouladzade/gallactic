package rpc

import (
	"encoding/json"

	"time"

	"github.com/gallactic/gallactic/common/binary"
	"github.com/gallactic/gallactic/core/account"
	"github.com/gallactic/gallactic/core/proposal"
	"github.com/gallactic/gallactic/core/validator"
	"github.com/gallactic/gallactic/crypto"
	"github.com/gallactic/gallactic/txs"
	"github.com/tendermint/go-amino"
	consensusTypes "github.com/tendermint/tendermint/consensus/types"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/rpc/core/types"
	tmTypes "github.com/tendermint/tendermint/types"
)

// When using Tendermint types like Block and Vote we are forced to wrap the outer object and use amino marshalling
var aminoCodec = amino.NewCodec()

func init() {
	core_types.RegisterAmino(aminoCodec)
}

type StorageOutput struct {
	Key   binary.HexBytes
	Value binary.HexBytes
}

type AccountsOutput struct {
	BlockHeight uint64
	Accounts    []*account.Account
}

type DumpstorageOutput struct {
	StorageItems []StorageItem
}

type StorageItem struct {
	Key   binary.HexBytes
	Value binary.HexBytes
}

type BlocksOutput struct {
	LastHeight uint64
	BlockMetas []*tmTypes.BlockMeta
}

type BlockOutput struct {
	BlockMeta *BlockMeta
	Block     *Block
}

type BlockMeta struct {
	*tmTypes.BlockMeta
}

func (bm BlockMeta) MarshalJSON() ([]byte, error) {
	return aminoCodec.MarshalJSON(bm.BlockMeta)
}

func (bm *BlockMeta) UnmarshalJSON(data []byte) (err error) {
	return aminoCodec.UnmarshalJSON(data, &bm.BlockMeta)
}

// Needed for amino handling of interface types
type Block struct {
	*tmTypes.Block
}

func (b Block) MarshalJSON() ([]byte, error) {
	return aminoCodec.MarshalJSON(b.Block)
}

func (b *Block) UnmarshalJSON(data []byte) (err error) {
	return aminoCodec.UnmarshalJSON(data, &b.Block)
}

type StatusOutput struct {
	NodeInfo          p2p.NodeInfo
	GenesisHash       binary.HexBytes
	PubKey            crypto.PublicKey
	LatestBlockHash   binary.HexBytes
	LatestBlockHeight uint64
	LatestBlockTime   int64
	NodeVersion       string
}

type LastBlockInfoOutput struct {
	LastBlockHeight uint64
	LastBlockTime   time.Time
	LastBlockHash   binary.HexBytes
}

type ChainIdOutput struct {
	ChainName   string
	ChainId     string
	GenesisHash binary.HexBytes
}

type Peer struct {
	NodeInfo   p2p.NodeInfo
	IsOutbound bool
}

type NetInfoOutput struct {
	Listening bool
	Listeners []string
	Peers     []*Peer
}

type ValidatorsOutput struct {
	BlockHeight         uint64
	BondedValidators    []*validator.Validator
	UnbondingValidators []*validator.Validator
}

type DumpConsensusStateOutput struct {
	RoundState      consensusTypes.RoundStateSimple
	PeerRoundStates []*consensusTypes.PeerRoundState
}

type PeersOutput struct {
	Peers []*Peer
}

type AccountOutput struct {
	Account *account.Account
}

type BroadcastTxOutput struct {
	txs.Receipt
}

func (rbt BroadcastTxOutput) MarshalJSON() ([]byte, error) {
	return json.Marshal(rbt.Receipt)
}

func (rbt BroadcastTxOutput) UnmarshalJSON(data []byte) (err error) {
	return json.Unmarshal(data, &rbt.Receipt)
}

type UnconfirmedTxsOutput struct {
	Count int
	Txs   []*txs.Envelope
}

type GenesisOutput struct {
	Genesis *proposal.Genesis
}

type BlockTxsOutput struct {
	Count int
	Txs   []txs.Envelope
}
