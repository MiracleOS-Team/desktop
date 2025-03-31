package main

import (
	"io"
	"net/http"
	"os"

	"github.com/MiracleOS-Team/desktoplib/foreignToplevel"
	basedir "github.com/MiracleOS-Team/libxdg-go/baseDir"
	"github.com/go-git/go-git/v5"
)

func checkRequiredDirectoriesAndCreate() error {
	_, mosSoft := os.Stat("/opt/miracleos-software")

	if os.IsNotExist(mosSoft) {
		err := os.Mkdir("/opt/miracleos-software", os.ModePerm)
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
		return err
	}

	// download labwc data

	_, autostartFile := os.Stat(basedir.GetXDGDirectory("config").(string) + "/labwc/autostart")
	if os.IsNotExist(autostartFile) {
		downloadFile("https://github.com/MiracleOS-Team/Dotfiles/raw/refs/heads/main/labwc/autostart", basedir.GetXDGDirectory("config").(string)+"/labwc/autostart")
	}

	_, environmentFile := os.Stat(basedir.GetXDGDirectory("config").(string) + "/labwc/environment")
	if os.IsNotExist(environmentFile) {
		downloadFile("https://github.com/MiracleOS-Team/Dotfiles/raw/refs/heads/main/labwc/environment", basedir.GetXDGDirectory("config").(string)+"/labwc/environment")
	}

	_, menuFile := os.Stat(basedir.GetXDGDirectory("config").(string) + "/labwc/menu.xml")
	if os.IsNotExist(menuFile) {
		downloadFile("https://github.com/MiracleOS-Team/Dotfiles/raw/refs/heads/main/labwc/menu.xml", basedir.GetXDGDirectory("config").(string)+"/labwc/menu.xml")
	}

	_, rcFile := os.Stat(basedir.GetXDGDirectory("config").(string) + "/labwc/rc.xml")
	if os.IsNotExist(rcFile) {
		downloadFile("https://github.com/MiracleOS-Team/Dotfiles/raw/refs/heads/main/labwc/rc.xml", basedir.GetXDGDirectory("config").(string)+"/labwc/rc.xml")
	}

	return nil
}

func checkRequiredSoftwareAndInstall() error {
	_, err := foreignToplevel.ListToplevels()
	return err
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
