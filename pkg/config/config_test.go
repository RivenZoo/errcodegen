package config

import (
	"github.com/BurntSushi/toml"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

const testErrCodeConfig = `
PkgName = "errorcodes"
ClientCodePrefix = "4"
ServerCodePrefix = "5"
AppCode = "101"

[[Modules]]
	ModuleName = "model"
	ModuleCode = "01"
	[[Modules.ClientCodes]]
		Name = "Err01"
		ErrNumber = "01"
		Msg = "test error code"

	[[Modules.ServerCodes]]
		Name = "Err01"
		ErrNumber = "01"
		Msg = "test error code"

[[Modules]]
	ModuleCode = "02"
	[[Modules.ClientCodes]]
		Name = "Err01"
		ErrNumber = "01"
		Msg = "test error code"

	[[Modules.ServerCodes]]
		Name = "Err01"
		ErrNumber = "01"
		Msg = "test error code"
`

func TestMustInit(t *testing.T) {
	conf := ErrCodeConfig{}
	meta, err := toml.Decode(testErrCodeConfig, &conf)
	assert.Nil(t, err)
	_ = meta
	t.Logf("%v", viper.Get("AppCode"))

	assert.NotNil(t, conf)
	assert.NotEmpty(t, conf.Modules)

	t.Logf("config: %#v", conf)
}
