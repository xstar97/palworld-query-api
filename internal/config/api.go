package config

import (
	"io/ioutil"
	"os"
	"sync"
	"time"
)

var (
	mutex       sync.RWMutex
	cachedData  []byte
	lastUpdated time.Time
	cacheFile   = Config.CachePath+"cached_data.json" // Name of the file to store cached data
)

func init() {
	// Load cached data from file when the package is initialized
	loadCachedDataFromFile()
}

func loadCachedDataFromFile() {
	// Check if the cache file exists
	if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
		return // Cache file does not exist, no data to load
	}

	// Read cached data from file
	data, err := ioutil.ReadFile(cacheFile)
	if err != nil {
		// Error reading cache file
		// Log the error or handle it accordingly
		return
	}

	// Update cached data and lastUpdated time
	mutex.Lock()
	defer mutex.Unlock()
	cachedData = data
	lastUpdated = time.Now()
}

func SetCachedData(data []byte) {
	mutex.Lock()
	defer mutex.Unlock()
	cachedData = data
	lastUpdated = time.Now()

	// Save cached data to file
	err := ioutil.WriteFile(cacheFile, data, 0644)
	if err != nil {
		// Error saving cache file
		// Log the error or handle it accordingly
	}
}

func GetCachedData() []byte {
	mutex.RLock()
	defer mutex.RUnlock()
	return cachedData
}

func NeedToUpdateCache() bool {
	// Update cache every hour
	return time.Since(lastUpdated) >= time.Hour
}
