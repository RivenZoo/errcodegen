package gen

import (
	"github.com/RivenZoo/errcodegen/pkg/config"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestModuleErrorCodeGenerator_GenerateModuleErrorCode(t *testing.T) {
	g := moduleErrorCodeGenerator{}
	commonConfig := &config.ErrCodeCommonConfig{
		PkgName:         "test",
		NewErrorFuncPkg: "golang.org/errors",
		NewErrorFunc:    "New",
	}
	moduleConfig := &config.ErrCodeModuleConfig{
		ModuleName: "test_module",
		ClientCodes: []config.ErrCodeVariableConfig{
			config.ErrCodeVariableConfig{
				Name:      "ErrNo1",
				ErrNumber: "01",
				Msg:       "error msg",
			},
		},
	}
	err := g.GenerateModuleErrorCode(commonConfig, moduleConfig)
	assert.Nil(t, err)

	fname := fixGoSuffix(moduleConfig.ModuleName)
	if _, err := os.Stat(fname); os.IsNotExist(err) {
		t.Fatalf("generate file %s fail", fname)
	}
	defer os.Remove(fname)
}
