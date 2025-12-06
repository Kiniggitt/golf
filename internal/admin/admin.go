// Package admin provides UI components for the admin dashboard.
// It handles displaying user information in table format and integrating
// with the Fyne UI framework.
package admin

import (
	"golf/internal/membership"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// ListUsers converts all registered users into a 2D string array for display.
// Each row contains: [Username, Full Name, Address, Balance, Days Since Last Fee].
// Returns an empty slice if there are no users. Nil users are skipped.
func ListUsers() [][]string {
	allUsers := membership.GetUsers()
	if len(allUsers) == 0 {
		return [][]string{}
	}
	users := make([][]string, len(allUsers))
	for i, user := range allUsers {
		if user == nil {
			continue
		}
		users[i] = make([]string, 5)
		users[i][0] = user.Username
		users[i][1] = user.GetName()
		users[i][2] = user.Address
		users[i][3] = user.GetBalance()
		users[i][4] = user.GetDaysSinceLastFee()
	}
	return users
}

// GetUserTable creates and returns a configured Fyne table widget for displaying users.
// The table has five columns: Username (120px), Name (180px), Address (250px),
// Balance (100px), and Days Since Last Fee (120px). The table automatically
// refreshes its data when the underlying user list changes.
func GetUserTable() *widget.Table {
	headers := []string{"Username", "Name", "Address", "Balance", "Days Since Last Fee"}

	table := widget.NewTableWithHeaders(
		func() (int, int) {
			data := ListUsers()
			if len(data) == 0 {
				return 0, len(headers)
			}
			return len(data), len(headers) // Number of rows and columns in the table body
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Template") // Template for table body cells
		},
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)
			data := ListUsers()
			if id.Row >= 0 && id.Row < len(data) && data[id.Row] != nil && id.Col >= 0 && id.Col < len(data[id.Row]) {
				label.SetText(data[id.Row][id.Col]) // Update table body cell content
			} else {
				label.SetText("")
			}
		},
	)

	table.UpdateHeader = func(id widget.TableCellID, obj fyne.CanvasObject) {
		label := obj.(*widget.Label)
		if id.Col >= 0 && id.Col < len(headers) {
			label.SetText(headers[id.Col])
		}
	}

	table.SetColumnWidth(0, 120) // Username
	table.SetColumnWidth(1, 180) // Name
	table.SetColumnWidth(2, 250) // Address
	table.SetColumnWidth(3, 100) // Balance
	table.SetColumnWidth(4, 120) // Days Since Last Fee

	return table
}
