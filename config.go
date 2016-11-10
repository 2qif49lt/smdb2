package main

import (
	"github.com/2qif49lt/toml"

	"fmt"
	"io"
	"os"
	"path/filepath"
)

var (
	ErrConfigFileMiss = fmt.Errorf(`config file miss`)
)

type ConfigFile struct {
	Admin []string
	Smtp  struct {
		Email string
		Pwd   string
		Srv   string
		Port  int
	}
	Mas struct {
		Workdir string
		Restart string
		Appxml  struct {
			Run     string
			Normal  string
			Standby string
		}
	}
	Filename string `toml:"-"`
}

func (configFile *ConfigFile) Load() error {
	if _, err := os.Stat(configFile.Filename); err == nil {
		file, err := os.Open(configFile.Filename)
		if err != nil {
			return fmt.Errorf("%s - %v", configFile.Filename, err)
		}
		defer file.Close()

		err = configFile.LoadFromReader(file)
		if err != nil {
			err = fmt.Errorf("%s - %v", configFile.Filename, err)
			return err
		}
		return nil
	} else {
		return ErrConfigFileMiss
	}
}

// LoadFromReader reads the configuration data given
func (configFile *ConfigFile) LoadFromReader(configData io.Reader) error {
	if err := toml.NewDecoder(configData).Decode(configFile); err != nil {
		return err
	}

	return nil
}

func (configFile *ConfigFile) SaveToWriter(writer io.Writer) error {

	data, err := toml.Marshal(*configFile)
	if err != nil {
		return err
	}
	_, err = writer.Write(data)
	return err
}

// Save encodes and writes out all information
func (configFile *ConfigFile) Save() error {
	if configFile.Filename == "" {
		return fmt.Errorf("Can't save config with empty filename")
	}

	if err := os.MkdirAll(filepath.Dir(configFile.Filename), 0700); err != nil {
		return err
	}
	f, err := os.OpenFile(configFile.Filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	return configFile.SaveToWriter(f)
}
