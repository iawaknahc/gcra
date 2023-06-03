package main

import (
	"time"
)

type GCRA struct {
	// EmissionInterval is the amount time 1 cell must wait in equally spaced schedule.
	EmissionInterval time.Duration
	// Tolerance is the capacity. It must be greater than or equal to 1.
	Tolerance int
}

type Result struct {
	// IsConforming is whether the cell is conforming.
	IsConforming bool
	// TAT is new TAT.
	// If IsConforming is true, TAT must be updated.
	TAT time.Time
	// AllowAt tells the earliest time the cell is conforming.
	// The value is meaningful only when IsConforming is false.
	AllowAt time.Time
}

func (a GCRA) IsConforming(tat time.Time, ta time.Time, quantity int) Result {
	// Determine TAT.
	// If TAT is not available, set TAT to the Actual Arrival Time (Ta).
	if tat.IsZero() {
		tat = ta
	}

	// Determine the new TAT.
	var newTAT time.Time
	// For 1 cell, we must wait EmissionInterval.
	// For n cells, we must wait n * EmissionInterval.
	increment := time.Duration(quantity) * a.EmissionInterval
	// If the cell arrives after TAT, then we calculate the new TAT with Ta.
	// This is equivalent to an empty bucket.
	// If the cell arrives before TAT, then we calculate the new TAT with TAT.
	// This is equivalent to an non-empty bucket.
	if ta.After(tat) {
		newTAT = ta.Add(increment)
	} else {
		newTAT = tat.Add(increment)
	}

	// Determine Delay Variation Tolerance.
	// We can think of it as how many cells do we allow before TAT.
	dvt := a.EmissionInterval * time.Duration(a.Tolerance)
	// Ta must be greater than or equal to allowAt in order to be conforming.
	allowAt := newTAT.Add(-dvt)
	if ta.Before(allowAt) {
		return Result{
			IsConforming: false,
			// TAT does not change if the cell is nonconforming.
			TAT:     tat,
			AllowAt: allowAt,
		}
	}
	return Result{
		IsConforming: true,
		TAT:          newTAT,
		AllowAt:      allowAt,
	}
}
