package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// Config - the config struct
type Config struct {
	configuration map[string]interface{}
	configFile    string
}

var defaultData = make(map[string]interface{})

// NewConfig - config Setup based the default Path
func NewConfig() *Config {
	path, ok := os.LookupEnv("CONFIG_PATH")
	if ok {
		return NewConfigWithFile(path)
	}
	return NewConfigWithFile("config.json")
}

// NewConfigWithFile - config Setup based on a file path
func NewConfigWithFile(configFile string) *Config {
	return &Config{
		configFile:    configFile,
		configuration: make(map[string]interface{}),
	}
}

func (c *Config) readConfig() error {
	var config map[string]interface{}
	jsonData, err := ioutil.ReadFile(c.configFile)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(jsonData), &config)
	if err != nil {
		return err
	}
	c.configuration = config
	return nil
}

func (c *Config) writeConfig() error {
	jsonData, err := json.MarshalIndent(c.configuration, "", "  ")
	if err != nil {
		return err
	}
	jsonFile, err := os.Create(c.configFile)

	if err != nil {
		return err
	}
	defer jsonFile.Close()

	jsonFile.Write(jsonData)
	return nil
}

// Check - checks wether a configuration file is already present, generate one otherwise
func (c *Config) Check() bool {
	_, err := os.Stat(c.configFile)
	if os.IsNotExist(err) {
		log.Printf("Could not find a configuration file at %s. Creating one!", c.configFile)
		c.writeConfig()
		c.readConfig()
		return false
	}
	c.readConfig()
	return true
}

// Get - read a config setting
func (c *Config) Get(key string) string {
	c.readConfig()
	value, present := c.configuration[key]
	if present {
		return value.(string)
	}
	return ""
}

// Set - writes a config setting
func (c *Config) Set(key string, value interface{}) {
	c.readConfig()
	c.configuration[key] = value
	c.writeConfig()
}

// GetAsArray - returns the config element as an string array
func (c *Config) GetAsArray(key string) []interface{} {
	c.readConfig()
	value, present := c.configuration[key].([]interface{})
	if present {
		return value
	}
	return nil
}

// GetAsMap - returns the config element as an map of string -> interface
func (c *Config) GetAsMap(key string) map[string]interface{} {
	c.readConfig()
	value, present := c.configuration[key].(map[string]interface{})
	if present {
		return value
	}
	return nil
}

// GetDefaultData - gets a prefilled default data map
func GetDefaultData() map[string]interface{} {
	return defaultData
}
