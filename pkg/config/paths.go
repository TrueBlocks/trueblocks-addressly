package config

import (
	"os"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
)

func GetCacheFolder(chain string /* may be empty */) string {
	dirName, _ := os.UserCacheDir()
	dirName += "/TrueBlocks/browse/"
	if chain != "" {
		dirName += chain + "/"
	}
	if !file.FolderExists(dirName) {
		_ = file.EstablishFolder(dirName)
	}
	// logger.Info("Cache dir: ", dirName, file.FolderExists(dirName))
	return dirName
}
