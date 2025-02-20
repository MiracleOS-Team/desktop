package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/MiracleOS-Team/desktoplib/batteryHandler"
	"github.com/MiracleOS-Team/desktoplib/foreignToplevel"
	"github.com/MiracleOS-Team/desktoplib/networkManagerHandler"
	"github.com/MiracleOS-Team/desktoplib/volumeHandler"
	"github.com/MiracleOS-Team/libxdg-go/notificationDaemon"
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

func createSidestuff(nDaemon *notificationDaemon.Daemon) *gtk.Box {
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
		glib.TimeoutAdd(uint(500), func() bool {

			newVolumeIcon, err := volumeHandler.GetAudioIcon()
			if err == nil {
				volumeImage.SetFromIconName(newVolumeIcon, gtk.ICON_SIZE_BUTTON)
			}

			// Return true to keep the timeout active.
			return true
		})

		statusBox.PackStart(volumeImage, false, false, 0)
	}

	networkIcon, err := networkManagerHandler.GetNetworkIcon()

	if err == nil {
		networkImage, _ := gtk.ImageNewFromIconName(networkIcon, gtk.ICON_SIZE_BUTTON)

		sc, _ = networkImage.GetStyleContext()
		sc.AddClass("network")

		glib.TimeoutAdd(uint(500), func() bool {

			networkIcon, err := networkManagerHandler.GetNetworkIcon()
			if err == nil {
				networkImage.SetFromIconName(networkIcon, gtk.ICON_SIZE_BUTTON)
			}

			// Return true to keep the timeout active.
			return true
		})

		statusBox.PackStart(networkImage, false, false, 0)
	}

	if batteryHandler.IsBattery() {
		batteryImage, _ := gtk.ImageNewFromIconName(batteryHandler.GetBatteryIcon(), gtk.ICON_SIZE_BUTTON)

		sc, _ = batteryImage.GetStyleContext()
		sc.AddClass("power")

		glib.TimeoutAdd(uint(500), func() bool {

			if batteryHandler.IsBattery() {
				newBatteryIcon := batteryHandler.GetBatteryIcon()
				batteryImage.SetFromIconName(newBatteryIcon, gtk.ICON_SIZE_BUTTON)
			}

			// Return true to keep the timeout active.
			return true
		})

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

	notificationBar := createNotificationBar(nDaemon)
	notificationButton.Connect("clicked", func() {
		if notificationBar.IsVisible() {
			notificationBar.Hide()
		} else {
			if len(nDaemon.Notifications) != 0 {
				notificationBar.ShowAll()
			}

		}
	})

	ntStack, _ := gtk.StackNew()

	notificationImage, _ := gtk.ImageNewFromIconName("preferences-system-notifications-symbolic", gtk.ICON_SIZE_BUTTON)
	sc, _ = notificationImage.GetStyleContext()
	sc.AddClass("notification-bell")

	notificationText, _ := gtk.LabelNew(strconv.Itoa(len(nDaemon.Notifications)))
	sc, _ = notificationText.GetStyleContext()
	sc.AddClass("h2")

	ntStack.Add(notificationImage)
	ntStack.Add(notificationText)

	notificationBox.PackStart(ntStack, false, false, 0)
	notificationButton.Add(notificationBox)

	glib.TimeoutAdd(uint(100), func() bool {
		// Get new date/time info.
		notificationText.SetText(strconv.Itoa(len(nDaemon.Notifications)))
		if len(nDaemon.Notifications) == 0 {
			ntStack.SetVisibleChild(notificationImage)
			ntStack.SetTransitionType(gtk.STACK_TRANSITION_TYPE_SLIDE_LEFT)
		} else {
			ntStack.SetVisibleChild(notificationText)
			ntStack.SetTransitionType(gtk.STACK_TRANSITION_TYPE_SLIDE_RIGHT)
		}

		// Return true to keep the timeout active.
		return true
	})

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

		pathn, err := foreignToplevel.GetIconFromToplevel(k, 16, 1)
		if err == nil {
			img, _ := gtk.ImageNewFromFile(pathn)
			imgButton.Add(img)
		}

		imgButton.Connect("clicked", func() {
			foreignToplevel.SelectToplevel(k)
		})
		box.PackStart(imgButton, false, false, 0)
	}

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

	mm := createMainMenu()

	customButton.Connect("clicked", func() {
		if mm.IsVisible() {
			mm.Hide()
		} else {
			mm.ShowAll()
		}
		//box.Hide()
	})

	box.PackStart(desktopImage, false, false, 0)
	box.PackStart(customButton, false, false, 0)
	box.PackStart(searchImage, false, false, 0)

	return box
}

func createBar(nDaemon *notificationDaemon.Daemon) *gtk.Window {
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
	box.PackEnd(createSidestuff(nDaemon), false, false, 0)

	glib.TimeoutAdd(uint(500), func() bool {
		chil := box.GetChildren()
		chil.NthData(uint(0)).(*gtk.Widget).Destroy()

		wspaces := createWorkspaces()
		box.PackStart(wspaces, false, false, 0)

		box.ShowAll()
		// Return true to keep the timeout active.
		return true
	})

	win.Add(box)
	return win
}
