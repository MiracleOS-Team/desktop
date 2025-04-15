package main

import (
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xprop"
	"github.com/MiracleOS-Team/desktoplib/wallpaper"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

func loadCSS() {
	// Load CSS into GTK
	provider, _ := gtk.CssProviderNew()
	//err = provider.LoadFromData(css)

	err := provider.LoadFromPath("desktop.css")
	if err != nil {
		err = provider.LoadFromPath("/opt/miracleos-software/desktop/desktop.css")
		if err != nil {
			log.Println("Failed to load CSS into GTK:", err)
			return
		}
	}

	display, err := gdk.DisplayGetDefault()
	if err != nil {
		log.Println("Failed to get default display:", err)
		return
	}

	screen, err := display.GetDefaultScreen()
	if err != nil {
		log.Println("Failed to get default screen:", err)
		return
	}

	gtk.AddProviderForScreen(screen, provider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
}

// scalePixbuf scales a Pixbuf while maintaining aspect ratio
func scalePixbuf(pixbuf *gdk.Pixbuf, maxWidth, maxHeight int) *gdk.Pixbuf {
	width := pixbuf.GetWidth()
	height := pixbuf.GetHeight()

	// Maintain aspect ratio
	scaleFactor := float64(maxWidth) / float64(width)
	if float64(maxHeight)/float64(height) < scaleFactor {
		scaleFactor = float64(maxHeight) / float64(height)
	}

	newWidth := int(float64(width) * scaleFactor)
	newHeight := int(float64(height) * scaleFactor)

	scaledPixbuf, err := pixbuf.ScaleSimple(newWidth, newHeight, gdk.INTERP_BILINEAR)
	if err != nil {
		log.Fatal("Could not scale image:", err)
	}
	return scaledPixbuf
}

func firstN(s string, n int) string {
	i := 0
	for j := range s {
		if i == n {
			return s[:j]
		}
		i++
	}
	return s
}

func setStrutPartial(xu *xgbutil.XUtil, win xproto.Window, height uint, screenWidth uint) error {
	// Reserve space at the bottom of the screen
	strutPartial := []uint{
		0, 0, 0, height, // left, right, top, bottom
		0, 0, 0, 0, // left_start, left_end, right_start, right_end
		0, 0, // top_start, top_end
		0, screenWidth - 1, // bottom_start, bottom_end
	}
	err := xprop.ChangeProp32(xu, win, "_NET_WM_STRUT_PARTIAL", "CARDINAL", strutPartial...)

	if err != nil {
		return fmt.Errorf("failed to set _NET_WM_STRUT_PARTIAL: %s", err)
	}

	xprop.ChangeProp32(xu, win, "_NET_WM_STRUT", "CARDINAL", 0, 0, 0, height)
	return nil
}

func main() {
	err := EnsureInstallation()
	if err != nil {
		panic(err)
	}

	wallpaper.SetImageWallpaper("/usr/share/backgrounds/miracleos_dark_default.jpg", "")

	gtk.Init(&os.Args)
	loadCSS()

	daemon := listenNotifications()
	defer daemon.Stop()

	bar := createBar(daemon)
	bar.ShowAll()

	gtk.Main()
}
