package main

import (
	"fmt"
	"golf/membership"
	"golf/registrar"
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func getRegisterContainer() (con *fyne.Container) {
	form := widget.NewForm()
	nameEntry := widget.NewEntry()
	addressEntry := widget.NewEntry()
	form.Append("Name:", nameEntry)
	form.Append("Address:", addressEntry)
	con = container.New(layout.NewGridLayout(1), form)
	form.OnSubmit = func() {
		newUser := membership.NewUser(nameEntry.Text, addressEntry.Text)
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
		} else if slices.Contains(keys, c) {
			con.Hidden = false
		} else {
			con.Hidden = true
		}
		con.Refresh()
	}
	appContainers["content"].Refresh()
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
		membership.UpdateUser(usernameLabel.Text, nameEntry.Text, addressEntry.Text)
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

func getHomeContainer() (con *fyne.Container) {
	signInButton := widget.NewButton("Sign In", func() {
		showAndHide("signin")
	})
	registerButton := widget.NewButton("Register", func() {
		showAndHide("register")
	})
	con = container.New(layout.NewGridLayout(2), registerButton, signInButton)
	return
}

var appContainers map[string]*fyne.Container = make(map[string]*fyne.Container)

func main() {
	a := app.New()
	w := a.NewWindow("Solid Rock Golf Course")
	appContainers["register"] = getRegisterContainer()
	appContainers["signin"] = getSignInContainer()
	appContainers["home"] = getHomeContainer()

	appContainers["content"] = container.New(layout.NewGridLayout(1))
	for key, con := range appContainers {
		if key != "content" {
			appContainers["content"].Add(con)
		}

	}

	w.Resize(fyne.NewSize(400, 200))

	membership.AddUsers(registrar.ReadJSON("test2.json"))
	fmt.Println(membership.GetUsers())

	w.SetContent(appContainers["content"])
	w.ShowAndRun()

	registrar.WriteJSON("test3.json", membership.GetUsers())

}
