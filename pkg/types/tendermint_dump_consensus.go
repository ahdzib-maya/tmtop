package types

type DumpConsensusStateResponse struct {
	Result *DumpConsensusStateResult `json:"result"`
}

type DumpConsensusStateResult struct {
	RoundState *DumpConsensusStateRoundState `json:"round_state"`
}

type DumpConsensusStateRoundState struct {
	Validators    DumpConsensusStateRoundStateValidators `json:"validators"`
	ProposalBlock *DumpConsensusProposalBlock            `json:"proposal_block"`
}

type DumpConsensusStateRoundStateValidators struct {
	Validators []TendermintValidator `json:"validators"`
}

type DumpConsensusProposalBlock struct {
	Header DumpConsensusProposalBlockHeader `json:"header"`
}

type DumpConsensusProposalBlockHeader struct {
	Height  string `json:"height"`
	AppHash string `json:"app_hash"`
}
