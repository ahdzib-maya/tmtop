package types

import (
	"errors"
	"math/big"
	"strings"
)

func ValidatorsWithLatestRoundFromTendermintResponse(
	consensus *ConsensusStateResponse,
	tendermintValidators []TendermintValidator,
	round int64,
) (ValidatorsWithRoundVote, error) {
	lastHeightVoteSet := consensus.Result.RoundState.HeightVoteSet[round]
	validators := make(ValidatorsWithRoundVote, len(lastHeightVoteSet.Prevotes))

	for index, prevote := range lastHeightVoteSet.Prevotes {
		precommit := lastHeightVoteSet.Precommits[index]
		validator := tendermintValidators[index]

		vp := new(big.Int)
		vp, ok := vp.SetString(validator.VotingPower, 10)
		if !ok {
			return nil, errors.New("error setting string")
		}

		prevoteVote, prevoteHash := VoteAndHashFromString(prevote)
		precommitVote, precommitHash := VoteAndHashFromString(precommit)

		validators[index] = ValidatorWithRoundVote{
			Validator: Validator{
				Address:     validator.Address,
				VotingPower: vp,
			},
			RoundVote: RoundVote{
				Address:             validator.Address,
				Prevote:             prevoteVote,
				PrevoteBlockHash:    prevoteHash,
				Precommit:           precommitVote,
				PrecommitBlockHash:  precommitHash,
				IsProposer:          validator.Address == consensus.Result.RoundState.Proposer.Address,
			},
		}
	}

	totalVP := validators.GetTotalVotingPower()

	for index, validator := range validators {
		validators[index].Validator.Index = index

		votingPowerPercent := big.NewFloat(0).SetInt(validator.Validator.VotingPower)
		votingPowerPercent = votingPowerPercent.Quo(votingPowerPercent, big.NewFloat(0).SetInt(totalVP))
		votingPowerPercent = votingPowerPercent.Mul(votingPowerPercent, big.NewFloat(100))

		validators[index].Validator.VotingPowerPercent = votingPowerPercent
	}

	return validators, nil
}

func ValidatorsWithAllRoundsFromTendermintResponse(
	consensus *ConsensusStateResponse,
	tendermintValidators []TendermintValidator,
) (ValidatorsWithAllRoundsVotes, error) {
	validators := make(Validators, len(tendermintValidators))
	for index, validator := range tendermintValidators {
		vp := new(big.Int)
		vp, ok := vp.SetString(validator.VotingPower, 10)
		if !ok {
			return ValidatorsWithAllRoundsVotes{}, errors.New("error setting string")
		}

		validators[index] = Validator{
			Address:     validator.Address,
			VotingPower: vp,
		}
	}

	totalVP := validators.GetTotalVotingPower()

	for index, validator := range validators {
		validators[index].Index = index

		votingPowerPercent := big.NewFloat(0).SetInt(validator.VotingPower)
		votingPowerPercent = votingPowerPercent.Quo(votingPowerPercent, big.NewFloat(0).SetInt(totalVP))
		votingPowerPercent = votingPowerPercent.Mul(votingPowerPercent, big.NewFloat(100))

		validators[index].VotingPowerPercent = votingPowerPercent
	}

	roundsVotes := make([]RoundVotes, len(consensus.Result.RoundState.HeightVoteSet))

	for round, roundHeightVoteSet := range consensus.Result.RoundState.HeightVoteSet {
		currentRoundVotes := make(RoundVotes, len(roundHeightVoteSet.Prevotes))

		for index, prevote := range roundHeightVoteSet.Prevotes {
			precommit := roundHeightVoteSet.Precommits[index]
			validator := tendermintValidators[index]

			prevoteVote, prevoteHash := VoteAndHashFromString(prevote)
			precommitVote, precommitHash := VoteAndHashFromString(precommit)

			currentRoundVotes[index] = RoundVote{
				Address:            validator.Address,
				Prevote:            prevoteVote,
				PrevoteBlockHash:   prevoteHash,
				Precommit:          precommitVote,
				PrecommitBlockHash: precommitHash,
				IsProposer:         validator.Address == consensus.Result.RoundState.Proposer.Address,
			}
		}

		roundsVotes[round] = currentRoundVotes
	}

	return ValidatorsWithAllRoundsVotes{
		Validators:  validators,
		RoundsVotes: roundsVotes,
	}, nil
}

// VoteAndHashFromString parses a vote string and returns the vote type and block hash.
// Vote format: Vote{idx:addr height/round/type(typeStr) BLOCKHASH signature @ timestamp}
// Returns (vote type, block hash)
func VoteAndHashFromString(source ConsensusVote) (Vote, string) {
	sourceStr := string(source)

	if sourceStr == "nil-Vote" {
		return VotedNil, ""
	}

	// Extract block hash: find the part after ") " and take next 12 chars
	// Format: ...SIGNED_MSG_TYPE_PREVOTE(Prevote) BLOCKHASH SIGNATURE...
	closingParenIdx := strings.Index(sourceStr, ") ")
	if closingParenIdx == -1 {
		// Malformed vote, return as voted with empty hash
		return Voted, ""
	}

	// Skip ") " to get to the block hash
	hashStart := closingParenIdx + 2
	if hashStart+12 > len(sourceStr) {
		// Not enough characters for hash
		return Voted, ""
	}

	blockHash := sourceStr[hashStart : hashStart+12]

	// Determine vote type
	if blockHash == "000000000000" {
		return VotedZero, blockHash
	}

	return Voted, blockHash
}

// VoteFromString returns just the vote type (for backward compatibility)
func VoteFromString(source ConsensusVote) Vote {
	vote, _ := VoteAndHashFromString(source)
	return vote
}
