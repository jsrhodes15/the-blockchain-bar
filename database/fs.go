package database

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func initDataDirIfNotExists(dataDir string) error {
	if fileExist(getGenesisJsonFilePath(dataDir)) {
		return nil
	}

	dbDir := getDatabaseDirPath(dataDir)
	err := os.MkdirAll(dbDir, os.ModePerm)
	if err != nil {
		return err
	}

	gen := getGenesisJsonFilePath(dataDir)
	err = writeGenesisToDisk(gen)
	if err != nil {
		return err
	}

	blocks := getBlocksDbFilePath(dataDir)
	err = writeEmptyBlocksDbToDisk(blocks)
	if err != nil {
		return err
	}

	return nil
}

func getDatabaseDirPath(dataDir string) string {
	return filepath.Join(dataDir, "database")
}

func getGenesisJsonFilePath(dataDir string) string {
	return filepath.Join(getDatabaseDirPath(dataDir), "genesis.json")
}

func getBlocksDbFilePath(dataDir string) string {
	return filepath.Join(getDatabaseDirPath(dataDir), "block.db")
}

func fileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		return false
	}

	return true
}

func dirExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	return true, err
}

func writeEmptyBlocksDbToDisk(path string) error {
	return ioutil.WriteFile(path, []byte(""), os.ModePerm)
}