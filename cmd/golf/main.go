package main

import (
	"fmt"
	"strconv"

	"golf/internal/admin"
	"golf/internal/config"
	"golf/internal/membership"
	"golf/internal/registrar"
	"golf/internal/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var appContainers = make(map[string]*fyne.Container)
var selectedUserIndex = -1
var adminTable *widget.Table
var appWindow fyne.Window
var appConfig *config.Config

func getRegisterContainer() (con *fyne.Container) {
	form := widget.NewForm()
	nameEntry := widget.NewEntry()
	addressEntry := widget.NewEntry()
	form.Append("Name:", nameEntry)
	form.Append("Address:", addressEntry)
	con = container.New(layout.NewGridLayout(1), form)
	form.OnSubmit = func() {
		newUser, err := membership.NewUser(nameEntry.Text, addressEntry.Text)
		if err != nil {
			dialog.ShowError(err, appWindow)
			return
		}
		ShowProfile(newUser.Username)
		showAndHide("profile")
		nameEntry.Text = ""
		addressEntry.Text = ""
	}
	form.OnCancel = func() {
		showAndHide("home")
	}
	con.Hidden = true
	return
}

func getSignInContainer() (con *fyne.Container) {
	form := widget.NewForm()
	usernameEntry := widget.NewEntry()
	form.Append("Username:", usernameEntry)
	con = container.New(layout.NewGridLayout(1), form)
	form.OnSubmit = func() {
		if membership.UsernameExists(usernameEntry.Text) {
			ShowProfile(usernameEntry.Text)
			showAndHide("profile")
			usernameEntry.Text = ""
		}
	}
	form.OnCancel = func() {
		showAndHide("home")
	}
	con.Hidden = true
	return
}

func showAndHide(keys ...string) {
	fmt.Println(keys)
	for c, con := range appContainers {
		if c == "content" {
			continue
		} else if contains(keys, c) {
			con.Hidden = false
		} else {
			con.Hidden = true
		}
		con.Refresh()
	}
	appContainers["content"].Refresh()
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func newProfileContainer(user *membership.User) {
	form := widget.NewForm()

	nameEntry := widget.NewEntry()
	nameEntry.SetText(user.GetName())

	addressEntry := widget.NewMultiLineEntry()
	addressEntry.SetText(user.Address)

	usernameLabel := widget.NewLabel(user.Username)
	form.Append("Username:", usernameLabel)
	form.Append("Name:", nameEntry)
	form.Append("Address:", addressEntry)
	form.Append("Balance:", widget.NewLabel(user.GetBalance()))
	form.Append("Fees:", widget.NewLabel(user.GetFees()))

	form.OnSubmit = func() {
		err := membership.UpdateUser(usernameLabel.Text, nameEntry.Text, addressEntry.Text)
		if err != nil {
			dialog.ShowError(err, appWindow)
			return
		}
		showAndHide("home")
	}
	form.OnCancel = func() {
		showAndHide("home")
	}
	con := container.New(layout.NewGridLayout(1), form)
	appContainers["profile"] = con
	appContainers["content"].Add(appContainers["profile"])
}

func ShowProfile(user string) {
	profile, ok := membership.GetUser(user)
	if ok {
		newProfileContainer(profile)
	}
}

func getAdminPasswordContainer() (con *fyne.Container) {
	form := widget.NewForm()
	passwordEntry := widget.NewPasswordEntry()
	form.Append("Admin Password:", passwordEntry)
	con = container.New(layout.NewGridLayout(1), form)
	form.OnSubmit = func() {
		if passwordEntry.Text == appConfig.AdminPassword {
			appWindow.Resize(fyne.NewSize(appConfig.Window.AdminWidth, appConfig.Window.AdminHeight))
			showAndHide("admin")
			passwordEntry.Text = ""
		} else {
			passwordEntry.Text = ""
		}
	}
	form.OnCancel = func() {
		showAndHide("home")
		passwordEntry.Text = ""
	}
	con.Hidden = true
	return
}

func getHomeContainer() (con *fyne.Container) {
	// Main buttons in center - make them larger with bigger text
	signInButton := widget.NewButton("Sign In", func() {
		showAndHide("signin")
	})
	signInButton.Importance = widget.HighImportance

	registerButton := widget.NewButton("Register", func() {
		showAndHide("register")
	})
	registerButton.Importance = widget.HighImportance

	// Small admin button for lower left (normal importance gives it an outline)
	adminButton := widget.NewButton("Admin", func() {
		showAndHide("adminpassword")
	})

	// Create center content with main buttons - horizontal layout with padding
	centerButtons := container.NewPadded(
		container.New(layout.NewGridLayout(2), registerButton, signInButton),
	)

	// Use border layout to place admin button in bottom-left corner
	bottomLeft := container.NewHBox(adminButton, layout.NewSpacer())

	con = container.New(
		layout.NewBorderLayout(
			nil,        // Top
			bottomLeft, // Bottom (with admin button on left)
			nil,        // Left
			nil,        // Right
		),
		bottomLeft,    // Add the bottom content
		centerButtons, // Center the main buttons
	)
	return
}

func showAddFeeDialog() {
	if selectedUserIndex < 0 {
		dialog.ShowInformation("No Selection", "Please select a user from the table first.", appWindow)
		return
	}

	users := membership.GetUsers()
	if selectedUserIndex >= len(users) {
		return
	}
	selectedUser := users[selectedUserIndex]

	// Get standard fees for dropdown
	standardFees := membership.GetStandardFees()
	feeOptions := []string{"Custom Fee"}
	for _, fee := range standardFees {
		feeOptions = append(feeOptions, fee.Name)
	}

	amountEntry := widget.NewEntry()
	amountEntry.SetPlaceHolder("0.00")
	reasonEntry := widget.NewEntry()
	reasonEntry.SetPlaceHolder("Fee reason")

	// Fee type selector
	feeSelect := widget.NewSelect(feeOptions, func(selected string) {
		if selected != "Custom Fee" {
			// Find the selected standard fee and populate fields
			for _, fee := range standardFees {
				if fee.Name == selected {
					amountEntry.SetText(fmt.Sprintf("%.2f", fee.Amount))
					reasonEntry.SetText(fee.Name)
					break
				}
			}
		} else {
			// Clear for custom entry
			amountEntry.SetText("")
			reasonEntry.SetText("")
		}
	})
	feeSelect.Selected = "Custom Fee"

	items := []*widget.FormItem{
		widget.NewFormItem("User", widget.NewLabel(selectedUser.Username)),
		widget.NewFormItem("Fee Type", feeSelect),
		widget.NewFormItem("Amount ($)", amountEntry),
		widget.NewFormItem("Reason", reasonEntry),
	}

	d := dialog.NewForm("Assign Fee", "Add Fee", "Cancel", items, func(submitted bool) {
		if submitted {
			amount, err := strconv.ParseFloat(amountEntry.Text, 64)
			if err != nil || amount <= 0 {
				dialog.ShowError(fmt.Errorf("invalid amount"), appWindow)
				return
			}
			if reasonEntry.Text == "" {
				dialog.ShowError(fmt.Errorf("reason is required"), appWindow)
				return
			}
			selectedUser.AddFee(amount, reasonEntry.Text)
			refreshAdminTable()
		}
	}, appWindow)
	d.Resize(fyne.NewSize(400, 350))
	d.Show()
}

func showEditUserDialog() {
	if selectedUserIndex < 0 {
		dialog.ShowInformation("No Selection", "Please select a user from the table first.", appWindow)
		return
	}

	users := membership.GetUsers()
	if selectedUserIndex >= len(users) {
		return
	}
	selectedUser := users[selectedUserIndex]

	nameEntry := widget.NewEntry()
	nameEntry.SetText(selectedUser.GetName())
	addressEntry := widget.NewMultiLineEntry()
	addressEntry.SetText(selectedUser.Address)

	items := []*widget.FormItem{
		widget.NewFormItem("Username", widget.NewLabel(selectedUser.Username)),
		widget.NewFormItem("Name", nameEntry),
		widget.NewFormItem("Address", addressEntry),
	}

	d := dialog.NewForm("Edit User", "Save", "Cancel", items, func(submitted bool) {
		if submitted {
			if nameEntry.Text == "" {
				dialog.ShowError(fmt.Errorf("name is required"), appWindow)
				return
			}
			err := membership.UpdateUser(selectedUser.Username, nameEntry.Text, addressEntry.Text)
			if err != nil {
				dialog.ShowError(err, appWindow)
				return
			}
			refreshAdminTable()
		}
	}, appWindow)
	d.Resize(fyne.NewSize(400, 300))
	d.Show()
}

func showUserDetailsDialog() {
	if selectedUserIndex < 0 {
		dialog.ShowInformation("No Selection", "Please select a user from the table first.", appWindow)
		return
	}

	users := membership.GetUsers()
	if selectedUserIndex >= len(users) {
		return
	}
	selectedUser := users[selectedUserIndex]

	// Build detailed fee list
	feeDetails := "Fees:\n"
	if len(selectedUser.Fees) == 0 {
		feeDetails += "No fees"
	} else {
		for i, fee := range selectedUser.Fees {
			status := "UNPAID"
			if fee.Paid {
				status = "PAID"
			}
			feeDetails += fmt.Sprintf("%d. $%.2f - %s [%s]\n", i+1, fee.Amount, fee.Reason, status)
		}
	}

	content := widget.NewLabel(fmt.Sprintf(
		"Username: %s\nName: %s\nAddress: %s\n\nBalance: %s\n\n%s",
		selectedUser.Username,
		selectedUser.GetName(),
		selectedUser.Address,
		selectedUser.GetBalance(),
		feeDetails,
	))
	content.Wrapping = fyne.TextWrapWord

	d := dialog.NewCustom("User Details", "Close", content, appWindow)
	d.Resize(fyne.NewSize(500, 400))
	d.Show()
}

func refreshAdminTable() {
	if adminTable != nil {
		adminTable.Refresh()
	}
}

func showManageStandardFeesDialog() {
	var selectedFeeID int = -1

	// Create list of fees
	feeList := widget.NewList(
		func() int {
			return len(membership.GetStandardFees())
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Template")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			fees := membership.GetStandardFees()
			if id < len(fees) {
				obj.(*widget.Label).SetText(fmt.Sprintf("%s - $%.2f", fees[id].Name, fees[id].Amount))
			}
		},
	)

	feeList.OnSelected = func(id widget.ListItemID) {
		selectedFeeID = int(id)
	}

	addButton := widget.NewButton("Add New Fee", func() {
		showAddStandardFeeDialog(feeList)
	})

	deleteButton := widget.NewButton("Delete Selected", func() {
		if selectedFeeID >= 0 {
			fees := membership.GetStandardFees()
			if selectedFeeID < len(fees) {
				membership.DeleteStandardFee(fees[selectedFeeID].ID)
				selectedFeeID = -1
				feeList.Refresh()
			}
		}
	})

	buttonBox := container.NewHBox(addButton, deleteButton)
	content := container.NewBorder(nil, buttonBox, nil, nil, feeList)

	d := dialog.NewCustom("Manage Standard Fees", "Close", content, appWindow)
	d.Resize(fyne.NewSize(400, 500))
	d.Show()
}

func showAddStandardFeeDialog(feeList *widget.List) {
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Fee name")
	amountEntry := widget.NewEntry()
	amountEntry.SetPlaceHolder("0.00")

	items := []*widget.FormItem{
		widget.NewFormItem("Fee Name", nameEntry),
		widget.NewFormItem("Amount ($)", amountEntry),
	}

	d := dialog.NewForm("Add Standard Fee", "Create", "Cancel", items, func(submitted bool) {
		if submitted {
			amount, err := strconv.ParseFloat(amountEntry.Text, 64)
			if err != nil || amount <= 0 {
				dialog.ShowError(fmt.Errorf("invalid amount"), appWindow)
				return
			}
			if nameEntry.Text == "" {
				dialog.ShowError(fmt.Errorf("name is required"), appWindow)
				return
			}
			_, err = membership.NewStandardFee(nameEntry.Text, amount)
			if err != nil {
				dialog.ShowError(err, appWindow)
				return
			}
			if feeList != nil {
				feeList.Refresh()
			}
		}
	}, appWindow)
	d.Resize(fyne.NewSize(350, 250))
	d.Show()
}

func getAdminContainer() (con *fyne.Container) {
	adminTable = admin.GetUserTable()
	adminTable.OnSelected = func(id widget.TableCellID) {
		selectedUserIndex = id.Row
	}

	viewDetailsButton := widget.NewButton("View Details", showUserDetailsDialog)
	assignFeeButton := widget.NewButton("Assign Fee", showAddFeeDialog)
	editUserButton := widget.NewButton("Edit User", showEditUserDialog)
	manageFeesButton := widget.NewButton("Manage Standard Fees", showManageStandardFeesDialog)
	backButton := widget.NewButton("Back to Home", func() {
		selectedUserIndex = -1
		appWindow.Resize(fyne.NewSize(appConfig.Window.DefaultWidth, appConfig.Window.DefaultHeight))
		showAndHide("home")
	})

	buttonBox := container.New(layout.NewGridLayout(5), viewDetailsButton, assignFeeButton, editUserButton, manageFeesButton, backButton)
	con = container.New(layout.NewBorderLayout(nil, buttonBox, nil, nil), buttonBox, adminTable)
	con.Hidden = true
	return
}

func main() {
	// Load configuration
	var err error
	appConfig, err = config.Load("config.json")
	if err != nil {
		fmt.Printf("Error loading config: %v\nUsing defaults.\n", err)
		appConfig = config.Default()
	}

	a := app.New()
	a.SetIcon(ui.CreateGolfBallIcon())
	w := a.NewWindow("Solid Rock Golf Course")
	appWindow = w

	appContainers["register"] = getRegisterContainer()
	appContainers["signin"] = getSignInContainer()
	appContainers["home"] = getHomeContainer()
	appContainers["adminpassword"] = getAdminPasswordContainer()
	appContainers["admin"] = getAdminContainer()

	appContainers["content"] = container.New(layout.NewGridLayout(1))
	for key, con := range appContainers {
		if key != "content" {
			appContainers["content"].Add(con)
		}
	}

	w.Resize(fyne.NewSize(appConfig.Window.DefaultWidth, appConfig.Window.DefaultHeight))

	// Load users and standard fees
	users, err := registrar.ReadJSON(appConfig.GetUsersPath())
	if err != nil {
		fmt.Printf("Warning: Failed to load users: %v\n", err)
	} else {
		membership.AddUsers(users)
	}

	fees, err := registrar.ReadStandardFeesJSON(appConfig.GetStandardFeesPath())
	if err != nil {
		fmt.Printf("Warning: Failed to load standard fees: %v\n", err)
	} else {
		membership.AddStandardFees(fees)
	}
	fmt.Println(membership.GetUsers())

	w.SetContent(appContainers["content"])
	w.ShowAndRun()

	// Save users and standard fees
	if err := registrar.WriteJSON(appConfig.GetUsersPath(), membership.GetUsers()); err != nil {
		fmt.Printf("Error: Failed to save users: %v\n", err)
	}
	if err := registrar.WriteStandardFeesJSON(appConfig.GetStandardFeesPath(), membership.GetStandardFees()); err != nil {
		fmt.Printf("Error: Failed to save standard fees: %v\n", err)
	}
}
