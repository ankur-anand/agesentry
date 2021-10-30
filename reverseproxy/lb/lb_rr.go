package lb

import (
	"sync/atomic"

	"github.com/ankur-anand/agesentry/reverseproxy/internal"
)

// RoundRobinBalancer selects a proxy based on round-robin ordering.
type RoundRobinBalancer struct {
	idx uint64
}

func (r *RoundRobinBalancer) DispatchTo(sd internal.StateDirector, sortedName []string) string {
	n := uint64(len(sortedName))
	if n == 0 {
		return ""
	}

	for i := uint64(0); i < n; i++ {
		atomic.AddUint64(&r.idx, 1)
		rp := sortedName[r.idx%n]
		if sd.IsHealthy(rp) {
			return rp
		}
	}
	return ""
}

// NewRoundRobinBalancer returns a load Balancer implementing
// Round Robin Load Balancing Algorithm.
func NewRoundRobinBalancer() *RoundRobinBalancer {
	return &RoundRobinBalancer{idx: 0}
}
