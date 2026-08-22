package video

import (
	"testing"
	"time"
)

func TestSeqTrackerFirstSequence(t *testing.T) {
	var tr seqTracker
	if tr.advance(1000) {
		t.Fatal("first sequence number must not be reported as loss")
	}
}

func TestSeqTrackerContiguous(t *testing.T) {
	var tr seqTracker
	tr.advance(100)
	for seq := 101; seq <= 110; seq++ {
		if tr.advance(uint16(seq)) {
			t.Fatalf("contiguous sequence %d must not be reported as loss", seq)
		}
	}
}

func TestSeqTrackerDetectsGap(t *testing.T) {
	var tr seqTracker
	tr.advance(100)
	if !tr.advance(103) {
		t.Fatal("jump from 100 to 103 must be reported as loss")
	}
	// Continuing from the new position is contiguous again.
	if tr.advance(104) {
		t.Fatal("contiguous sequence after gap must not be reported as loss")
	}
}

func TestSeqTrackerWraparound(t *testing.T) {
	var tr seqTracker
	tr.advance(65534)
	if tr.advance(65535) {
		t.Fatal("65534 -> 65535 is contiguous")
	}
	if tr.advance(0) {
		t.Fatal("65535 -> 0 is contiguous (wraparound)")
	}
	if !tr.advance(2) {
		t.Fatal("0 -> 2 must be reported as loss across wraparound")
	}
}

func TestSeqTrackerReorderAndDuplicate(t *testing.T) {
	var tr seqTracker
	tr.advance(10)
	if tr.advance(9) {
		t.Fatal("late-arriving packet must not be reported as loss")
	}
	if tr.advance(10) {
		t.Fatal("duplicate packet must not be reported as loss")
	}
	if !tr.advance(13) {
		t.Fatal("jump from 10 to 13 must be reported as loss")
	}
}

func TestRateLimiterAllowsThenBlocks(t *testing.T) {
	l := newRateLimiter(500 * time.Millisecond)
	if !l.allow() {
		t.Fatal("first call must be allowed")
	}
	if l.allow() {
		t.Fatal("immediate second call must be blocked")
	}
}

func TestRateLimiterAllowsAfterInterval(t *testing.T) {
	l := newRateLimiter(10 * time.Millisecond)
	if !l.allow() {
		t.Fatal("first call must be allowed")
	}
	time.Sleep(15 * time.Millisecond)
	if !l.allow() {
		t.Fatal("call after interval must be allowed")
	}
}
