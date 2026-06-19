package main

const (
	topbarHeight    = 2
	cellWidth       = 3
	cellHeight      = 1
	gutterWidth     = 12
	bottomBarHeight = 4
)

type layout struct {
	gardenWidth  int
	gardenHeight int
	gridStartX   int
	gridStartY   int
}

func (m model) View() string {
	switch m.currentScreen {
	case titleScreen:
		return m.title.View()
	case gardenScreen:
		return m.garden.View()
	}
	return ""
}
