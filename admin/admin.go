package admin

import (
	"fmt"
	"golf/membership"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func ListUsers() [][]string {
	users := make([][]string, len(membership.GetUsers()))
	for i := range users {
		users[i] = make([]string, 3)
		for j := range users[i] {
			users[i][j] = fmt.Sprintf("Cell %d,%d", i, j)
		}
	}
	return users
}

func GetUserTable() *widget.Table {
	data := ListUsers()
	table := widget.NewTableWithHeaders(
		func() (int, int) {
			return len(data), len(data[0]) // Number of rows and columns in the table body
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Template") // Template for table body cells
		},
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)
			label.SetText(data[id.Row][id.Col]) // Update table body cell content
		},
	)
	return table
}
