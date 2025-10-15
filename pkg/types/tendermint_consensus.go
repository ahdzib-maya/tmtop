package types

import (
	"time"
)

type ConsensusStateResponse struct {
	Result *ConsensusStateResult `json:"result"`
}

type ConsensusStateResult struct {
	RoundState *ConsensusStateRoundState `json:"round_state"`
}

type ConsensusStateRoundState struct {
	HeightRoundStep string                   `json:"height/round/step"`
	StartTime       time.Time                `json:"start_time"`
	HeightVoteSet   []ConsensusHeightVoteSet `json:"height_vote_set"`
	Proposer        ConsensusStateProposer   `json:"proposer"`
	Proposal        *ConsensusProposal       `json:"proposal"`
}

type ConsensusHeightVoteSet struct {
	Round              int                   `json:"round"`
	Prevotes           []ConsensusVote       `json:"prevotes"`
	Precommits         []ConsensusVote       `json:"precommits"`
	PrevotesBitArray   ConsensusVoteBitArray `json:"prevotes_bit_array"`
	PrecommitsBitArray ConsensusVoteBitArray `json:"precommits_bit_array"`
}

type ConsensusStateProposer struct {
	Address string `json:"address"`
	Index   int    `json:"index"`
}

type ConsensusProposal struct {
	Height  string       `json:"height"`
	Round   int          `json:"round"`
	BlockID *BlockID     `json:"block_id"`
}

type BlockID struct {
	Hash  string        `json:"hash"`
	Parts *BlockIDParts `json:"parts"`
}

type BlockIDParts struct {
	Total int    `json:"total"`
	Hash  string `json:"hash"`
}

type ConsensusVote string
type ConsensusVoteBitArray string
