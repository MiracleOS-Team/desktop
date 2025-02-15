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
