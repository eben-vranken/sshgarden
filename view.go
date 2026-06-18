package main

const (
	topbarHeight     = 2
	sidebarFullWidth = 20
	cellWidth        = 3
	cellHeight       = 1
	gutterWidth      = 12
)

type layout struct {
	sidebarWidth  int
	gardenWidth   int
	gardenHeight  int
	sidebarHeight int
	gridStartX    int
	gridStartY    int
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
