package initconfig

import (
	"embed"
	"io"
	"os"
)

//go:embed config-example.toml
var configFile embed.FS

// writeDefaultConfigFile write default configuration file
func writeDefaultConfigFile(filename string) error {
	src, err := configFile.Open("config-example.toml")
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}
	return nil
}
