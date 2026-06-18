package main

import "time"

type plot struct {
	plant       *plant
	seededAt    time.Time
	lastWatered time.Time
}

func (p plot) getGrowthVisual() rune {
	if p.plant == nil {
		return ' '
	}

	elapsed := time.Since(p.seededAt)
	fraction := float64(elapsed) / float64(p.plant.growthTime)
	index := int(fraction * float64(len(p.plant.growthRunes)))

	index = min(index, len(p.plant.growthRunes)-1)

	return p.plant.growthRunes[index]
}

func (p plot) isOccupied() bool {
	return p.plant != nil
}
