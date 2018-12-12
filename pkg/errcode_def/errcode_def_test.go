package errcode_def

import (
	"bytes"
	"errcodegen/pkg/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestModuleErrorCodeParser_ParseErrorCodeDefinition(t *testing.T) {
	testDefinition := []byte(`
[moduleA(module_code=01)] # moduleA is module name; (module_code=01) is optional
[[client_error]] 		  # client_error is keyword
ErrModuleAVar1 = "error msg 1"
ErrModuleAVar2 = "error msg 2"

[[server_error]] 		 # server_error is keyword

[moduleB]
[[server_error]] 		 # server_error is keyword
ErrModuleBVar1 = "error msg 2"`)
	expect := []config.ErrCodeModuleConfig{
		config.ErrCodeModuleConfig{
			ModuleName: "moduleA",
			ModuleCode: "01",
			ClientCodes: []config.ErrCodeVariableConfig{
				config.ErrCodeVariableConfig{
					Name:      "ErrModuleAVar1",
					ErrNumber: "01",
					Msg:       "error msg 1",
				},
				config.ErrCodeVariableConfig{
					Name:      "ErrModuleAVar2",
					ErrNumber: "02",
					Msg:       "error msg 2",
				},
			},
		},
		config.ErrCodeModuleConfig{
			ModuleName: "moduleB",
			ModuleCode: "02",
			ServerCodes: []config.ErrCodeVariableConfig{
				config.ErrCodeVariableConfig{
					Name:      "ErrModuleBVar1",
					ErrNumber: "01",
					Msg:       "error msg 2",
				},
			},
		},
	}
	p := NewModuleErrorCodeParser()
	modulesConfig, err := p.ParseErrorCodeDefinition(bytes.NewReader(testDefinition))
	assert.Nil(t, err)
	assert.EqualValues(t, expect, modulesConfig)
}
