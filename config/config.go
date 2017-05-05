package config

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
)

const (
	// ConfigFileName is the name of config file
	ConfigFileName = "config.json"
	configFileDir  = ".docking"
)

var (
	configDir = os.Getenv("DOCKING_CONFIG")
)

func init() {
	if configDir == "" {
		configDir = filepath.Join(os.Getenv("HOME"), configFileDir)
	}
}

// ConfigDir returns the directory the configuration file is stored in
func ConfigDir() string {
	return configDir
}

type ConfigFile struct {
	LogLevel string          `json:"logLevel,omitpempty"`
	HostIp   string          `json:"hostIp,omitpempty"`
	Targets  []*ConfigTarget `json:"targets"`
	filename string          // non persistent
}
type ConfigTarget struct {
	Name        string                       `json:"name"`
	UrlTemplate string                       `json:"url"`
	Url         string                       // Note: not serialized - for internal use only
	HttpHeaders map[string]string            `json:"httpHeaders,omitpempty"`
	Templates   map[string][]*ConfigTemplate `json:"templates"`
}

// NewConfigFile initializes an empty configuration file for the given filename 'fn'
func NewConfigFile(fn string) ConfigFile {
	return ConfigFile{
		LogLevel: "debug",
		Targets:  make([]*ConfigTarget, 0),
		filename: fn,
	}
}

// Filename returns the name of the configuration file
func (configFile *ConfigFile) Filename() string {
	return configFile.filename
}

// LoadFromReader reads the configuration data given
// information with given directory and populates the receiver object
func (configFile *ConfigFile) LoadFromReader(configData io.Reader) error {
	if err := json.NewDecoder(configData).Decode(&configFile); err != nil {
		return err
	}
	return nil
}

func Load(configDir string) (*ConfigFile, error) {
	if configDir == "" {
		configDir = ConfigDir()
	}

	filePath := filepath.Join(configDir, ConfigFileName)
	log.Info("Read config file: ", filePath)
	configFile := NewConfigFile(filePath)

	if _, err := os.Stat(configFile.filename); err == nil {
		file, err := os.Open(configFile.filename)
		if err != nil {
			return &configFile, fmt.Errorf("%s - %v", configFile.filename, err)
		}
		defer file.Close()
		err = configFile.LoadFromReader(file)
		if err != nil {
			err = fmt.Errorf("%s - %v", configFile.filename, err)
		}
		level, err := log.ParseLevel(configFile.LogLevel)
		if err != nil {
			err = fmt.Errorf("%s - %v", configFile.filename, err)
		}
		log.SetLevel(level)
		return &configFile, err
	} else {
		// if file is there but we can't stat it for any reason other
		// than it doesn't exist then stop
		return &configFile, fmt.Errorf("%s - %v", configFile.filename, err)
	}

}

func (configFile *ConfigFile) SaveToWriter(writer io.Writer) error {
	data, err := json.MarshalIndent(configFile, "", "\t")
	if err != nil {
		return err
	}
	_, err = writer.Write(data)
	return err
}

func (configFile *ConfigFile) Save() error {
	if configFile.Filename() == "" {
		configFile.filename = filepath.Join(configDir, ConfigFileName)
	}
	fmt.Fprintf(os.Stderr, "Save file %v/%v", configFile.filename, configFile.filename)

	if err := os.MkdirAll(filepath.Dir(configFile.filename), 0700); err != nil {
		return err
	}
	f, err := os.OpenFile(configFile.filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	return configFile.SaveToWriter(f)
}
