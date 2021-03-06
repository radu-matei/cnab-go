package manifest

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/viper"

	"github.com/radu-matei/cnab-go/pkg/duffle"
)

// Load opens the named file for reading. If successful, the manifest is returned.
func Load(name, dir string) (*Manifest, error) {
	v := viper.New()
	if name == "" {
		v.SetConfigName(duffle.DuffleFilename)
	} else {
		v.SetConfigFile(filepath.Join(dir, name))
	}
	v.AddConfigPath(dir)
	err := v.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("Error finding duffle config file: %s", err)
	}

	m := New()
	err = v.Unmarshal(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
