package config

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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
    DockerUrl       string                      `json:"dockerUrl"`
    RegisterUrl     string                      `json:"registerUrl"`
	HostIp          string                      `json:"hostIp,omitpempty"`
	Templates       map[string][]*ConfigTemplate  `json:"templates"`
    filename        string                      // non persistent
}



// NewConfigFile initializes an empty configuration file for the given filename 'fn'
func NewConfigFile(fn string) ConfigFile {
	return ConfigFile{
		Templates: make(map[string][]*ConfigTemplate),
		filename:    fn,
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
	log.Println("Read config file: ", filePath)
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
	fmt.Fprintf(os.Stderr, "Save file %v/%v", configFile.filename, configFile.filename )
	
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