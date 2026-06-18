package main

func (g gardenModel) renderSidebarBody() string {
	selectedPlot := g.gardenGrid[g.selectedPlot.y][g.selectedPlot.x]

	if selectedPlot.isOccupied() {
		return selectedPlot.plant.name
	} else {
		return "Empty plot\n\n[c] Plant Carrot"
	}
}
