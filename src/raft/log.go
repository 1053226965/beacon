package raft

type LogEntry struct {
	term    int64
	index   int64
	content any

	prev *LogEntry
	next *LogEntry
}

type log struct {
	head *LogEntry
	tail *LogEntry
}

func (l *log) getLastTermAndIndex() (term int64, index int64) {
	if l.tail == nil {
		return 0, 0
	}
	return l.tail.term, l.tail.index
}

func (l *log) insertLog(term int64, index int64, context any) {
	if l.tail == nil {
		l.tail = &LogEntry{term: term, index: index, content: context}
		l.head = l.tail
		return
	}
	l.tail.next = &LogEntry{term: term, index: index, content: context}
	l.tail.next.prev = l.tail
	l.tail = l.tail.next
}
