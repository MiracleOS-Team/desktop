package main

import (
	"log"
	"os"

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
