
package config

import (
	"testing"
	"fmt"

)

func TestConfigGeneration(t *testing.T) {
	//polarisConfig := newPolarisConfig()
	//assert.Equal(t, polarisConfig, Parameters)
	//defaultConfig := newDefaultConfig()
	//assert.NotEqual(t, defaultConfig, polarisConfig)

	config := DefConfig
	fmt.Printf("%+v",config.CrossChain)
}
