package main

import (
	"testing"
	"time"
)

type event struct {
	t            time.Time
	quantity     int
	isConforming bool
	allowAt      time.Time
}

func sec(s int64) time.Time {
	return time.Unix(s, 0)
}

func TestGCRA(t *testing.T) {
	test := func(gcra GCRA, events []event) {
		var tat time.Time
		for _, e := range events {
			result := gcra.IsConforming(tat, e.t, e.quantity)
			if result.IsConforming != e.isConforming {
				t.Errorf("t = %v: %v != %v", e.t, result.IsConforming, e.isConforming)
			}
			if result.AllowAt != e.allowAt {
				t.Errorf("t = %v: %v != %v", e.t, result.AllowAt, e.allowAt)
			}
			if result.IsConforming {
				tat = result.TAT
			}
		}
	}

	test(GCRA{
		EmissionInterval: 1 * time.Second,
		Tolerance:        3,
	}, []event{
		{sec(0), 0, true, sec(-3)},
		{sec(0), 0, true, sec(-3)},

		{sec(0), 1, true, sec(-2)},
		{sec(0), 1, true, sec(-1)},
		{sec(0), 1, true, sec(0)},
		{sec(0), 1, false, sec(1)},
		{sec(0), 2, false, sec(2)},

		{sec(1), 2, false, sec(2)},
		{sec(1), 3, false, sec(3)},
		{sec(1), 1, true, sec(1)},
		{sec(1), 1, false, sec(2)},

		{sec(4), 0, true, sec(1)},
		{sec(4), 1, true, sec(2)},
		{sec(4), 1, true, sec(3)},
		{sec(4), 1, true, sec(4)},
		{sec(4), 1, false, sec(5)},
	})
}
