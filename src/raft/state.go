package raft

type State struct {
	nodeNum     int
	myself      int64
	currentTerm int64
	voteFor     int64
	logIndex    int64
	commitIndex int64
	lastApplied int64
	log         log
}
