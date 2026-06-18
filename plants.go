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
		growthTime:       time.Hour * 2,
		waterRequirement: time.Minute * 30,
		thirstTolerance:  time.Minute * 30,
		growthRunes:      []rune{'.', ',', 'i', 'Y'},
		color:            colorOrange,
	},
}
