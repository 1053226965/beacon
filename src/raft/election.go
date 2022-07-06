package raft

type Vote struct {
	Term         int64
	CandidateId  int64
	LastLogIndex int64
	LastLogTerm  int64
}

type VoteResult struct {
	Term        int64
	whoVoted    int64
	votedFor    int64
	VoteGranted bool
}

type IElection interface {
	ResetTerm(term int64)
	NewTerm()
	RequestVote(v *Vote) (vr *VoteResult)
	Grant(vr *VoteResult) (lead bool)
}

type election struct {
	st           *State
	grantedNodes map[int64]bool
}

const (
	NotVoted = -1
)

func NewElection(st *State) *election {
	if st == nil {
		panic("got nil state")
	}
	return &election{st: st, grantedNodes: map[int64]bool{}}
}

func (e *election) ResetTerm(term int64) {
	if term <= e.st.currentTerm || term < 0 {
		return
	}
	e.grantedNodes = make(map[int64]bool)
	e.st.currentTerm = term
	e.st.voteFor = NotVoted
}

func (e *election) RequestVote(v *Vote) (vr *VoteResult) {
	vr = &VoteResult{}
	vr.Term = v.Term
	vr.whoVoted = e.st.myself
	vr.votedFor = v.CandidateId
	vr.VoteGranted = false
	if v.Term < e.st.currentTerm {
		return vr
	}
	if v.Term > e.st.currentTerm {
		e.ResetTerm(v.Term)
	}
	lastLogTerm, lastLogIndex := e.st.log.getLastTermAndIndex()
	if e.st.voteFor != NotVoted ||
		v.LastLogTerm < lastLogTerm ||
		v.LastLogTerm == lastLogTerm && v.LastLogIndex < lastLogIndex {
		return vr
	}
	vr.VoteGranted = true
	return vr
}

func (e *election) Grant(vr *VoteResult) (lead bool) {
	if vr.Term != e.st.currentTerm {
		return false
	}
	if vr.votedFor != e.st.myself {
		return false
	}
	e.grantedNodes[vr.whoVoted] = true
	return len(e.grantedNodes) > e.st.nodeNum/2
}
