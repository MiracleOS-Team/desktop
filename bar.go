package main

import (
	"fmt"
	"time"

	"github.com/MiracleOS-Team/desktoplib/batteryHandler"
	"github.com/MiracleOS-Team/desktoplib/foreignToplevel"
	"github.com/MiracleOS-Team/desktoplib/networkManagerHandler"
	"github.com/MiracleOS-Team/desktoplib/volumeHandler"
	"github.com/dlasky/gotk3-layershell/layershell"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func getDateInfo() (string, string) {
	hours, minutes, _ := time.Now().Clock()
	curTimeInString := fmt.Sprintf("%d:%02d", hours, minutes)

	curDay := time.Now().Day()
	curDayName := firstN(time.Now().Weekday().String(), 3)
	var curDayCal string

	if curDay <= 5 {
		curMonth := time.Now().Month()
		curDayCal = fmt.Sprintf("%s. %02d %s", curDayName, curDay, curMonth)
	} else {
		curDayCal = fmt.Sprintf("%s. %02d", curDayName, curDay)
	}
	return curDayCal, curTimeInString
}

func createSidestuff() *gtk.Box {
	sideBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	sideBox.SetHAlign(gtk.ALIGN_END)
	sc, _ := sideBox.GetStyleContext()
	sc.AddClass("sidestuff")

	otherIcons, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	keyboardImage, _ := gtk.ImageNewFromIconName("input-keyboard-symbolic", gtk.ICON_SIZE_BUTTON)
	sc, _ = keyboardImage.GetStyleContext()
	sc.AddClass("keyboard")

	otherIcons.PackStart(keyboardImage, false, false, 0)
	sc, _ = otherIcons.GetStyleContext()
	sc.AddClass("other-icons-wrapper")

	statusBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)

	volumeIcon, err := volumeHandler.GetAudioIcon()

	if err == nil {
		volumeImage, _ := gtk.ImageNewFromIconName(volumeIcon, gtk.ICON_SIZE_BUTTON)
		sc, _ = volumeImage.GetStyleContext()
		sc.AddClass("sound")

		statusBox.PackStart(volumeImage, false, false, 0)
	}

	networkIcon, err := networkManagerHandler.GetNetworkIcon()

	if err == nil {
		networkImage, _ := gtk.ImageNewFromIconName(networkIcon, gtk.ICON_SIZE_BUTTON)

		sc, _ = networkImage.GetStyleContext()
		sc.AddClass("network")

		statusBox.PackStart(networkImage, false, false, 0)
	}

	if batteryHandler.IsBattery() {
		batteryImage, _ := gtk.ImageNewFromIconName(batteryHandler.GetBatteryIcon(), gtk.ICON_SIZE_BUTTON)

		sc, _ = batteryImage.GetStyleContext()
		sc.AddClass("power")

		statusBox.PackStart(batteryImage, false, false, 0)
	}

	sc, _ = statusBox.GetStyleContext()
	sc.AddClass("status-icons-wrapper")

	curDayCal, currUTCTimeInString := getDateInfo()

	clock, _ := gtk.LabelNew(currUTCTimeInString)
	sc, _ = clock.GetStyleContext()
	sc.AddClass("clock-text")

	dayText, _ := gtk.LabelNew(curDayCal)
	sc, _ = dayText.GetStyleContext()
	sc.AddClass("day-text")

	glib.TimeoutAdd(uint(500), func() bool {
		// Get new date/time info.
		newDayCal, newTime := getDateInfo()

		// Update the labels.
		clock.SetText(newTime)
		dayText.SetText(newDayCal)

		// Return true to keep the timeout active.
		return true
	})

	notificationButton, _ := gtk.ButtonNew()
	notificationBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	sc, _ = notificationBox.GetStyleContext()
	sc.AddClass("notification-bell-wrapper")

	notificationImage, _ := gtk.ImageNewFromIconName("preferences-system-notifications-symbolic", gtk.ICON_SIZE_BUTTON)
	sc, _ = notificationImage.GetStyleContext()
	sc.AddClass("notification-bell")

	notificationBox.PackStart(notificationImage, false, false, 0)
	notificationButton.Add(notificationBox)

	sideBox.PackStart(otherIcons, false, false, 0)
	sideBox.PackStart(statusBox, false, false, 0)
	sideBox.PackStart(clock, false, false, 0)
	sideBox.PackStart(dayText, false, false, 0)
	sideBox.PackStart(notificationButton, false, false, 0)

	return sideBox
}

func createWorkspaces() *gtk.Box {
	box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	box.SetHAlign(gtk.ALIGN_START)
	sc, _ := box.GetStyleContext()
	sc.AddClass("workspaces")

	toplevels, err := foreignToplevel.ListToplevels()
	if err != nil {
		fmt.Println("Error getting toplevels:", err)
		return box
	}

	for _, k := range toplevels {
		imgButton, _ := gtk.ButtonNew()
		sc, _ := imgButton.GetStyleContext()
		sc.AddClass("app")

		img, _ := gtk.ImageNewFromIconName(k.AppID, gtk.ICON_SIZE_BUTTON)
		imgButton.Add(img)
		box.PackStart(imgButton, false, false, 0)
	}
	// Placeholder for dynamic window list
	imgButton1, _ := gtk.ButtonNew()
	sc, _ = imgButton1.GetStyleContext()
	sc.AddClass("app")
	img1, _ := gtk.ImageNewFromIconName("preferences-desktop", gtk.ICON_SIZE_BUTTON)
	imgButton1.Add(img1)
	box.PackStart(imgButton1, false, false, 0)

	return box
}

func createMainIcons() *gtk.Box {
	box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	box.SetHAlign(gtk.ALIGN_CENTER)

	desktopImage, _ := gtk.ImageNewFromIconName("preferences-system-windows-symbolic", gtk.ICON_SIZE_LARGE_TOOLBAR)
	searchImage, _ := gtk.ImageNewFromIconName("system-search-symbolic", gtk.ICON_SIZE_LARGE_TOOLBAR)
	customIcon, _ := gtk.ImageNewFromFile("images/pp.png")
	customButton, _ := gtk.ButtonNew()
	customButton.Add(customIcon)

	customButton.Connect("clicked", func() {
		createMainMenu().ShowAll()
		box.Hide()
	})

	box.PackStart(desktopImage, false, false, 0)
	box.PackStart(customButton, false, false, 0)
	box.PackStart(searchImage, false, false, 0)

	return box
}

func createBar() *gtk.Window {
	win, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	win.SetTitle("Main Bar")
	win.SetDecorated(false)
	win.SetResizable(false)
	win.SetTypeHint(gdk.WINDOW_TYPE_HINT_DOCK)

	layershell.InitForWindow(win)
	layershell.SetNamespace(win, "miracleos")
	layershell.SetAnchor(win, layershell.LAYER_SHELL_EDGE_LEFT, true)
	layershell.SetAnchor(win, layershell.LAYER_SHELL_EDGE_BOTTOM, true)
	layershell.SetAnchor(win, layershell.LAYER_SHELL_EDGE_RIGHT, true)

	layershell.SetLayer(win, layershell.LAYER_SHELL_LAYER_TOP)
	layershell.SetMargin(win, layershell.LAYER_SHELL_EDGE_TOP, 0)
	layershell.SetMargin(win, layershell.LAYER_SHELL_EDGE_LEFT, 0)
	layershell.SetMargin(win, layershell.LAYER_SHELL_EDGE_RIGHT, 0)

	layershell.SetExclusiveZone(win, 75)
	layershell.SetKeyboardMode(win, layershell.LAYER_SHELL_KEYBOARD_MODE_NONE)
	disp, _ := gdk.DisplayGetDefault()
	mon, _ := disp.GetMonitor(0)
	layershell.SetMonitor(win, mon)

	box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	sc, _ := box.GetStyleContext()
	sc.AddClass("bar")
	box.PackStart(createWorkspaces(), false, false, 0)
	box.SetCenterWidget(createMainIcons())
	box.PackEnd(createSidestuff(), false, false, 0)
	win.Add(box)
	return win
}
