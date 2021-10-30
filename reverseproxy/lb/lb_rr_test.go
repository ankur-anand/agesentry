package lb_test

import (
	"reflect"
	"testing"

	"github.com/ankur-anand/agesentry/reverseproxy/lb"
)

type sdMock struct {
	db map[string]bool
}

func (s sdMock) IsHealthy(name string) bool {
	return s.db[name]
}

func (s sdMock) LoadCount(name string) (int, bool) {
	panic("implement me")
}

func TestRoundRobinBalancer_DispatchTo(t *testing.T) {
	names := []string{
		"A",
		"B",
		"C",
		"D",
		"E",
		"F",
	}

	db := map[string]bool{
		"A": true,
		"B": true,
		"C": false,
		"D": true,
		"E": false,
		"F": false,
	}

	sdm := sdMock{db: db}

	ld := lb.NewRoundRobinBalancer()
	output := make([]string, 0)
	for range names {
		name := ld.DispatchTo(sdm, names)
		output = append(output, name)
	}
	if !reflect.DeepEqual([]string{"B", "D", "A", "B", "D", "A"}, output) {
		t.Errorf("round robin load balancer wanted order didn't matched with got order")
	}
}
