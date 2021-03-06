/*
Copyright (c) 2015, Brian Hummer (brian@redq.me)
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package speciater

import (
	"github.com/rqme/neat"
)

const (
	CompatibilityThresholdFloor = 0.3
)

type DynamicSettings interface {
	SetCompatibilityThreshold(float64) // Threshold above which two genomes are not compatible
	TargetNumberOfSpecies() int        // The desired number of species
	CompatibilityModifier() float64    // Amount to change the compatibility threshold for next iteration
}

type Dynamic struct {
	DynamicSettings
	Classic
}

func NewDynamic(ds DynamicSettings, cs ClassicSettings) *Dynamic {
	return &Dynamic{
		DynamicSettings: ds,
		Classic:         Classic{ClassicSettings: cs},
	}
}

func (s *Dynamic) Speciate(curr []neat.Species, genomes []neat.Genome) (next []neat.Species, err error) {

	// Speciate using the internal speciater
	next, err = s.Classic.Speciate(curr, genomes)
	if err != nil {
		return
	}

	// Adjust the compatibily theshold as necessary
	ct := s.CompatibilityThreshold()
	if len(next) < s.TargetNumberOfSpecies() {
		ct -= s.CompatibilityModifier()
		if ct < CompatibilityThresholdFloor {
			ct = CompatibilityThresholdFloor
		}
	} else if len(next) > s.TargetNumberOfSpecies() {
		ct += s.CompatibilityModifier()
	}
	s.SetCompatibilityThreshold(ct)
	return
}
