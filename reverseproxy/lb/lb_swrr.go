package lb

import (
	"sync"

	"github.com/ankur-anand/agesentry/reverseproxy/internal"
)

// weightedName is a weighted target.
type weightedName struct {
	Item            string
	Weight          int64
	CurrentWeight   int64
	EffectiveWeight int64
}

// SWRR Provided smooth weight round robin algorithm implementation
type SWRR struct {
	mu    sync.Mutex
	items []*weightedName
	n     int
}

func NewSWRR() *SWRR {
	return &SWRR{items: make([]*weightedName, 0)}
}

func (w *SWRR) DispatchTo(sd internal.StateDirector, sortedName []string) string {
	// weighted should be within the provided names only
	// and shouldn't loop forever.
	for i := 0; i < len(sortedName); i++ {
		name := w.nextWeighted()
		if sd.IsHealthy(name.Item) {
			return name.Item
		}
	}

	return ""
}

// Add a weighted Item to the SWRR.
// This should be done before calling DispatchTo.
func (w *SWRR) Add(item string, weight int64) {
	weighted := &weightedName{Item: item, Weight: weight, EffectiveWeight: weight}
	w.items = append(w.items, weighted)
	w.n++
}

// nextWeighted returns next selected weighted object.
func (w *SWRR) nextWeighted() *weightedName {
	if w.n == 0 {
		return nil
	}
	if w.n == 1 {
		return w.items[0]
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	return nextSmoothWeighted(w.items)
}

func nextSmoothWeighted(items []*weightedName) (best *weightedName) {
	total := int64(0)

	for i := 0; i < len(items); i++ {
		w := items[i]

		if w == nil {
			continue
		}

		w.CurrentWeight += w.EffectiveWeight
		total += w.EffectiveWeight
		if w.EffectiveWeight < w.Weight {
			w.EffectiveWeight++
		}

		if best == nil || w.CurrentWeight > best.CurrentWeight {
			best = w
		}

	}

	if best == nil {
		return nil
	}

	best.CurrentWeight -= total
	return best
}
