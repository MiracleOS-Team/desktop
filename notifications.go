/*
	MiracleOS's desktop - A taskbar, with a main menu and a notification daemon
    Copyright (C) 2025 MiracleOS Contributors

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.

*/

package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/MiracleOS-Team/libxdg-go/notificationDaemon"
	"github.com/dlasky/gotk3-layershell/layershell"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func createNotification(notification *notificationDaemon.Notification, nDaemon *notificationDaemon.Daemon) *gtk.Box {
	notificationBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 15)
	sc, _ := notificationBox.GetStyleContext()
	sc.AddClass("ntf_main_div")

	ntfTopBar, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 15)

	sc, _ = ntfTopBar.GetStyleContext()
	sc.AddClass("ntf_top_bar")

	ntfTopBarText, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 15)

	sc, _ = ntfTopBarText.GetStyleContext()
	sc.AddClass("nf_topbar_text")

	ntfTopBarImage, _ := gtk.ImageNewFromIconName(notification.AppIcon, gtk.ICON_SIZE_BUTTON)

	ntfTopBarTextLabel, _ := gtk.LabelNew(notification.AppName)

	ntfTopBarDeleteButtonBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 15)

	ntfTopBarDeleteButton, _ := gtk.ButtonNew()

	sc, _ = ntfTopBarDeleteButton.GetStyleContext()
	sc.AddClass("button")

	ntfTopBarDeleteButtonLabel, _ := gtk.LabelNew("Close")

	ntfTopBarDeleteButton.Connect("clicked", func() {
		nDaemon.CloseNotificationAsUser(notification.ID)
	})

	ntfTopBarDeleteButton.Add(ntfTopBarDeleteButtonLabel)
	ntfTopBarDeleteButtonBox.PackEnd(ntfTopBarDeleteButton, false, false, 0)
	ntfTopBar.PackEnd(ntfTopBarDeleteButtonBox, false, false, 0)
	ntfTopBarText.PackStart(ntfTopBarImage, false, false, 0)
	ntfTopBarText.PackStart(ntfTopBarTextLabel, false, false, 0)
	ntfTopBar.PackStart(ntfTopBarText, false, false, 0)
	notificationBox.PackStart(ntfTopBar, false, false, 0)

	notificationBody, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 15)

	//TODO: Handle notification Image

	notificationTextBody, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 15)

	sc, _ = notificationTextBody.GetStyleContext()
	sc.AddClass("ntf_text_contents")

	if notification.Summary != "" {
		notificationSummary, _ := gtk.LabelNew(notification.Summary)
		notificationSummary.SetXAlign(0)
		sc, _ = notificationSummary.GetStyleContext()
		sc.AddClass("h2")
		notificationTextBody.PackStart(notificationSummary, false, false, 0)
	}

	if notification.Body != "" {
		notificationBody, _ := gtk.LabelNew(notification.Body)
		notificationBody.SetXAlign(0)
		notificationTextBody.PackStart(notificationBody, false, false, 0)

	}

	hours, minutes, _ := notification.Timestamp.Clock()

	timeLabel, _ := gtk.LabelNew(fmt.Sprintf("%d:%02d", hours, minutes))
	timeLabel.SetXAlign(0)
	sc, _ = timeLabel.GetStyleContext()
	sc.AddClass("h4")

	notificationTextBody.PackEnd(timeLabel, false, false, 0)

	notificationBody.PackStart(notificationTextBody, false, false, 0)
	notificationBox.PackStart(notificationBody, false, false, 0)

	return notificationBox
}

func createNotificationBarTitle(nDaemon *notificationDaemon.Daemon) *gtk.Box {
	tBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)

	title, _ := gtk.LabelNew(strings.Join([]string{strconv.Itoa(len(nDaemon.Notifications)), " Notifications"}, ""))
	sc, _ := title.GetStyleContext()
	sc.AddClass("h1")

	glib.TimeoutAdd(uint(100), func() bool {
		// Get new date/time info.
		title.SetText(strings.Join([]string{strconv.Itoa(len(nDaemon.Notifications)), " Notifications"}, ""))

		// Return true to keep the timeout active.
		return true
	})

	closeAllButtonText, _ := gtk.LabelNew("Clear all")
	sc, _ = closeAllButtonText.GetStyleContext()
	sc.AddClass("button")

	closeAllButton, _ := gtk.ButtonNew()
	closeAllButton.Add(closeAllButtonText)

	closeAllButton.Connect("clicked", func() {
		for _, elem := range nDaemon.Notifications {
			nDaemon.CloseNotificationAsUser(elem.ID)
		}
	})

	tBox.PackStart(title, false, false, 0)
	tBox.PackEnd(closeAllButton, false, false, 0)
	return tBox
}

func createNotificationBar(nDaemon *notificationDaemon.Daemon) *gtk.Window {
	win, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	win.SetTitle("Notification Bar")
	win.SetDecorated(false)
	win.SetResizable(false)
	win.SetKeepAbove(true)
	win.SetSkipTaskbarHint(true)
	win.SetSkipPagerHint(true)
	win.SetAppPaintable(true)

	layershell.InitForWindow(win)
	layershell.SetNamespace(win, "miracleos")

	layershell.SetAnchor(win, layershell.LAYER_SHELL_EDGE_BOTTOM, true)
	layershell.SetAnchor(win, layershell.LAYER_SHELL_EDGE_RIGHT, true)

	layershell.SetLayer(win, layershell.LAYER_SHELL_LAYER_TOP)
	layershell.SetMargin(win, layershell.LAYER_SHELL_EDGE_TOP, 0)
	layershell.SetMargin(win, layershell.LAYER_SHELL_EDGE_LEFT, 0)
	layershell.SetMargin(win, layershell.LAYER_SHELL_EDGE_RIGHT, 0)

	layershell.SetKeyboardMode(win, layershell.LAYER_SHELL_KEYBOARD_MODE_NONE)
	disp, _ := gdk.DisplayGetDefault()
	mon, _ := disp.GetMonitor(0)
	layershell.SetMonitor(win, mon)

	mBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)

	for _, nt := range nDaemon.Notifications {
		mBox.PackEnd(createNotification(&nt, nDaemon), false, false, 0)
	}

	mBox.PackStart(createNotificationBarTitle(nDaemon), false, false, 0)

	glib.TimeoutAdd(uint(500), func() bool {
		chil := mBox.GetChildren()

		cLength := int(chil.Length())

		for i := 0; i < cLength; i++ {
			chil.NthData(uint(i)).(*gtk.Widget).Destroy()
		}

		for _, nt := range nDaemon.Notifications {
			mBox.PackEnd(createNotification(&nt, nDaemon), false, false, 0)
		}

		mBox.PackStart(createNotificationBarTitle(nDaemon), false, false, 0)

		mBox.ShowAll()
		if len(nDaemon.Notifications) == 0 {
			win.Hide()
		}
		// Return true to keep the timeout active.
		return true
	})

	win.Add(mBox)

	return win
}

func listenNotifications() *notificationDaemon.Daemon {
	daemon := notificationDaemon.NewDaemon(notificationDaemon.Config{
		Capabilities: []string{"body", "actions", "actions-ions", "icon-static"},
	})
	if err := daemon.Start(); err != nil {
		log.Fatalf("Failed to start notification daemon: %v", err)
	}

	return daemon
}
