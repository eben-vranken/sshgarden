package main

import (
	"time"

	"github.com/charmbracelet/lipgloss"
)

type plant struct {
	name             string
	growthTime       time.Duration
	waterRequirement time.Duration
	thirstTolerance  time.Duration
	growthRunes      []rune
	color            lipgloss.Color
}

var plantRegistry = map[string]*plant{
	"carrot": {
		name:             "Carrot",
		growthTime:       time.Second * 20,
		waterRequirement: time.Second * 10,
		thirstTolerance:  time.Second * 5,
		growthRunes:      []rune{'.', ',', 'i', 'Y'},
		color:            colorRedOrange,
	},
}
