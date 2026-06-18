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
