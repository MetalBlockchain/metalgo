// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package poll

import (
	"fmt"

	"github.com/MetalBlockchain/metalgo/ids"
	"github.com/MetalBlockchain/metalgo/utils/bag"
)

type earlyTermNoTraversalFactory struct {
	alpha int
}

// NewEarlyTermNoTraversalFactory returns a factory that returns polls with
// early termination, without doing DAG traversals
func NewEarlyTermNoTraversalFactory(alpha int) Factory {
	return &earlyTermNoTraversalFactory{alpha: alpha}
}

func (f *earlyTermNoTraversalFactory) New(vdrs bag.Bag[ids.NodeID]) Poll {
	return &earlyTermNoTraversalPoll{
		polled: vdrs,
		alpha:  f.alpha,
	}
}

// earlyTermNoTraversalPoll finishes when any remaining validators can't change
// the result of the poll. However, does not terminate tightly with this bound.
// It terminates as quickly as it can without performing any DAG traversals.
type earlyTermNoTraversalPoll struct {
	votes  bag.Bag[ids.ID]
	polled bag.Bag[ids.NodeID]
	alpha  int
}

// Vote registers a response for this poll
func (p *earlyTermNoTraversalPoll) Vote(vdr ids.NodeID, vote ids.ID) {
	count := p.polled.Count(vdr)
	// make sure that a validator can't respond multiple times
	p.polled.Remove(vdr)

	// track the votes the validator responded with
	p.votes.AddCount(vote, count)
}

// Drop any future response for this poll
func (p *earlyTermNoTraversalPoll) Drop(vdr ids.NodeID) {
	p.polled.Remove(vdr)
}

// Finished returns true when one of the following conditions is met.
//
//  1. There are no outstanding votes.
//  2. It is impossible for the poll to achieve an alpha majority after applying
//     transitive voting.
//  3. A single element has achieved an alpha majority.
func (p *earlyTermNoTraversalPoll) Finished() bool {
	remaining := p.polled.Len()
	if remaining == 0 {
		return true // Case 1
	}

	received := p.votes.Len()
	maxPossibleVotes := received + remaining
	if maxPossibleVotes < p.alpha {
		return true // Case 2
	}

	_, freq := p.votes.Mode()
	return freq >= p.alpha // Case 3
}

// Result returns the result of this poll
func (p *earlyTermNoTraversalPoll) Result() bag.Bag[ids.ID] {
	return p.votes
}

func (p *earlyTermNoTraversalPoll) PrefixedString(prefix string) string {
	return fmt.Sprintf(
		"waiting on %s\n%sreceived %s",
		p.polled.PrefixedString(prefix),
		prefix,
		p.votes.PrefixedString(prefix),
	)
}

func (p *earlyTermNoTraversalPoll) String() string {
	return p.PrefixedString("")
}
