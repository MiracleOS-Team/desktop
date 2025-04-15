package main

import (
	"io"
	"net/http"
	"os"
	"os/exec"

	basedir "github.com/MiracleOS-Team/libxdg-go/baseDir"
	"github.com/go-git/go-git/v5"
)

func checkRequiredDirectoriesAndCreate() error {
	_, mosSoft := os.Stat("/opt/miracleos-software")

	if os.IsNotExist(mosSoft) {
		err := os.Mkdir("/opt/miracleos-software", os.ModePerm)
		return err
	}

	_, deskData := os.Stat("/opt/miracleos-software/desk-data")

	if os.IsNotExist(deskData) {
		err := os.Mkdir("/opt/miracleos-software/desk-data", os.ModePerm)
		return err
	}

	_, labWCConf := os.Stat(basedir.GetXDGDirectory("config").(string) + "/labwc")
	if os.IsNotExist(labWCConf) {
		err := os.Mkdir(basedir.GetXDGDirectory("config").(string)+"/labwc", os.ModePerm)
		return err
	}
	return nil
}

func downloadFile(downloadURL string, saveDir string) error {
	// Create blank file
	file, err := os.Create(saveDir)
	if err != nil {
		return err
	}
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	// Put content on file
	resp, err := client.Get(downloadURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	defer file.Close()
	return nil
}

func checkRequiredReposAndDownload() error {
	_, mosDesktop := os.Stat("/opt/miracleos-software/desktop")

	if os.IsNotExist(mosDesktop) {
		_, err := git.PlainClone("/opt/miracleos-software/desktop", false, &git.CloneOptions{
			URL:      "https://github.com/MiracleOS-Team/desktop",
			Progress: os.Stdout,
		})
		if err != nil {
			return err
		}
	}

	_, mosIcons := os.Stat("/opt/miracleos-software/Icons")

	if os.IsNotExist(mosIcons) {
		_, err := git.PlainClone("/opt/miracleos-software/Icons", false, &git.CloneOptions{
			URL:      "https://github.com/MiracleOS-Team/Icons",
			Progress: os.Stdout,
		})
		if err != nil {
			return err
		}
	}

	_, mosIconsInstall := os.Stat("/usr/share/icons/MiracleOSIcons")

	if os.IsNotExist(mosIconsInstall) {
		err := os.CopyFS("/usr/share/icons/MiracleOSIcons", os.DirFS("/opt/miracleos-software/Icons/MiracleOSIcons"))
		if err != nil {
			return err
		}
	}

	_, appIconsAlias := os.Stat("/opt/miracleos-software/desk-data/app-icons-alias.json")
	if os.IsNotExist(appIconsAlias) {
		err := downloadFile("https://github.com/MiracleOS-Team/desktoplib/raw/refs/heads/main/app-icons-alias.json", "/opt/miracleos-software/desk-data/app-icons-alias.json")
		if err != nil {
			return err
		}
	}

	_, wallpaperDark := os.Stat("/usr/share/backgrounds/miracleos_dark_default.jpg")
	if os.IsNotExist(wallpaperDark) {
		err := downloadFile("https://raw.githubusercontent.com/MiracleOS-Team/miracleos-team.github.io/refs/heads/main/brand/wallpaper.jpg", "/usr/share/backgrounds/miracleos_dark_default.jpg")
		if err != nil {
			return err
		}
	}

	// download labwc data

	_, autostartFile := os.Stat(basedir.GetXDGDirectory("config").(string) + "/labwc/autostart")
	if os.IsNotExist(autostartFile) {
		err := downloadFile("https://github.com/MiracleOS-Team/Dotfiles/raw/refs/heads/main/labwc/autostart", basedir.GetXDGDirectory("config").(string)+"/labwc/autostart")
		if err != nil {
			return err
		}
	}

	_, environmentFile := os.Stat(basedir.GetXDGDirectory("config").(string) + "/labwc/environment")
	if os.IsNotExist(environmentFile) {
		err := downloadFile("https://github.com/MiracleOS-Team/Dotfiles/raw/refs/heads/main/labwc/environment", basedir.GetXDGDirectory("config").(string)+"/labwc/environment")
		if err != nil {
			return err
		}
	}

	_, menuFile := os.Stat(basedir.GetXDGDirectory("config").(string) + "/labwc/menu.xml")
	if os.IsNotExist(menuFile) {
		err := downloadFile("https://github.com/MiracleOS-Team/Dotfiles/raw/refs/heads/main/labwc/menu.xml", basedir.GetXDGDirectory("config").(string)+"/labwc/menu.xml")
		if err != nil {
			return err
		}
	}

	_, rcFile := os.Stat(basedir.GetXDGDirectory("config").(string) + "/labwc/rc.xml")
	if os.IsNotExist(rcFile) {
		err := downloadFile("https://github.com/MiracleOS-Team/Dotfiles/raw/refs/heads/main/labwc/rc.xml", basedir.GetXDGDirectory("config").(string)+"/labwc/rc.xml")
		if err != nil {
			return err
		}
	}

	return nil
}

func checkRequiredSoftwareAndInstall() error {
	_, err := exec.LookPath("wlrctl")
	if err != nil {
		return err
	}

	_, err = exec.LookPath("swww")
	if err != nil {
		return err
	}

	_, err = exec.LookPath("mpvpaper")
	if err != nil {
		return err
	}

	return nil
}

func EnsureInstallation() error {

	err := checkRequiredDirectoriesAndCreate()
	if err != nil {
		return err
	}
	err = checkRequiredReposAndDownload()
	if err != nil {
		return err
	}
	err = checkRequiredSoftwareAndInstall()
	if err != nil {
		return err
	}
	return nil
}
