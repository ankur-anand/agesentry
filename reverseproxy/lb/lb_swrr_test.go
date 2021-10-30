package lb_test

import (
	"testing"

	"github.com/ankur-anand/agesentry/reverseproxy/lb"
)

func TestSWRR_DispatchTo(t *testing.T) {
	db := map[string]bool{
		"A": true,
		"B": true,
		"C": false,
		"D": true,
		"E": false,
		"F": true,
	}
	sdm := sdMock{db: db}
	// index gives the weight
	names := []string{
		"A",
		"B",
		"C",
		"D",
		"E",
		"F",
	}
	ld := lb.NewSWRR()
	for i, v := range names {
		ld.Add(v, int64(i+1))
	}

	result := make(map[string]int)
	for i := 0; i < 100000; i++ {
		s := ld.DispatchTo(sdm, names)
		result[s]++
	}
	if len(result) != 4 {
		t.Errorf("expected only four healthy node to be returned.")
	}
	// ratio "1:2:4:6"
	count := 0
	for _, v := range []string{"B", "D", "F"} {
		count = count + (result[v] / result["A"])
	}
	if count != 12 {
		t.Errorf("expected ratio of call failed from algorithm")
	}
}

func TestSWRR_DispatchToPattern(t *testing.T) {
	db := map[string]bool{
		"A": true,
		"B": true,
		"C": false,
		"D": true,
		"E": false,
		"F": true,
	}
	sdm := sdMock{db: db}
	// index gives the weight
	names := []string{
		"A",
		"B",
		"C",
		"D",
		"E",
		"F",
	}
	ld := lb.NewSWRR()
	for i, v := range names {
		ld.Add(v, int64(i+1))
	}

	pattern := "FDBFDFAFDFBDFFDBFDFA"
	result := ""
	for i := 0; i < 20; i++ {
		s := ld.DispatchTo(sdm, names)
		result = result + s
	}
	if result != pattern {
		t.Errorf("smooth weight pattern didn't match")
	}
}
