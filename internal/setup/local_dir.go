package setup

import (
	"encoding/json"
	"log"
	"os"
	"runtime"
)

type Config struct {
	ConfigPath   string `json:"config_path"`
	LogPath      string `json:"log_path"`
	DownloadPath string `json:"download_path"`
}

const (
	CONFIG_FILE   = "config.json"
	MAIN_DIR      = "download4/"
	LOG_DIR       = "logs"
	DOWNLOADS_DIR = "downloads"
	CONFIG_DIR    = "config"

	UNIX_CONFIG    = "/.config/"
	UNIX_LOGS      = "/.cache/"
	UNIX_DOWNLOADS = "/Documents/"

	WIN_CONFIG_LOGS = "/AppData/Roaming/"
	WIN_DOWNLOADS   = "/Documents/"
)

func ConfigSetup() Config {
	var cfg Config = detectOS()
	return cfg
}

func detectOS() Config {
	var os string = runtime.GOOS

	homeDir, err := getHomeDir()
	if err != nil {
		log.Fatalln("Unable to locate home dir", err)
	}

	var cfg Config

	switch os {
	case "linux", "darwin":
		cfg = unixSetup(homeDir, cfg)
	case "windows":
		cfg = winSetup(homeDir, cfg)
	default:
		log.Fatalln("This isn't Doom. You can't just run it on whatever this thing is... Detected OS:", os)
	}

	return cfg
}

func getHomeDir() (string, error) {
	return os.UserHomeDir()
}

func winSetup(homeDir string, cfg Config) Config {
	cfg.ConfigPath = homeDir + WIN_CONFIG_LOGS + MAIN_DIR + CONFIG_DIR
	cfg.LogPath = homeDir + WIN_CONFIG_LOGS + MAIN_DIR + LOG_DIR
	cfg.DownloadPath = homeDir + WIN_DOWNLOADS + MAIN_DIR + DOWNLOADS_DIR

	cfg = dirCheck(cfg)

	return cfg
}

func unixSetup(homeDir string, cfg Config) Config {
	cfg.ConfigPath = homeDir + UNIX_CONFIG + MAIN_DIR + CONFIG_DIR
	cfg.LogPath = homeDir + UNIX_LOGS + MAIN_DIR + LOG_DIR
	cfg.DownloadPath = homeDir + UNIX_DOWNLOADS + MAIN_DIR + DOWNLOADS_DIR

	cfg = dirCheck(cfg)

	return cfg
}

func readConfig(cfgFile string) (string, string) {
	c, err := os.ReadFile(cfgFile)
	if err != nil {
		log.Fatalln("Could not open config file.", err)
	}

	var t Config
	err = json.Unmarshal(c, &t)
	if err != nil {
		log.Fatalln("Could not read config.json.", err)
	}

	return t.LogPath, t.DownloadPath
}

func createDir(path, dirName string) {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		log.Fatalf("Error creating %s directory: %v\n", dirName, err)
	}
	log.Printf("%s directory created at %s\n", dirName, path)
}

func dirCheck(cfg Config) Config {
	var pathToCfgFile = cfg.ConfigPath + "/" + CONFIG_FILE

	_, err := os.Stat(cfg.ConfigPath)
	if os.IsNotExist(err) {
		createDir(cfg.ConfigPath, CONFIG_DIR)
	}

	_, err = os.Stat(pathToCfgFile)
	if !os.IsNotExist(err) {
		cfg.LogPath, cfg.DownloadPath = readConfig(pathToCfgFile)

		_, err := os.Stat(cfg.LogPath)
		if os.IsNotExist(err) {
			createDir(cfg.LogPath, LOG_DIR)
		}

		_, err = os.Stat(cfg.DownloadPath)
		if os.IsNotExist(err) {
			createDir(cfg.DownloadPath, DOWNLOADS_DIR)
		}
	} else {
		_, err = os.Stat(cfg.LogPath)
		if os.IsNotExist(err) {
			createDir(cfg.LogPath, LOG_DIR)
		}

		_, err = os.Stat(cfg.DownloadPath)
		if os.IsNotExist(err) {
			createDir(cfg.DownloadPath, DOWNLOADS_DIR)
		}
	}

	return cfg
}
