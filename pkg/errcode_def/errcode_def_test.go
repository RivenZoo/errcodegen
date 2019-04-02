package errcode_def

import (
	"bytes"
	"github.com/RivenZoo/errcodegen/pkg/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestModuleErrorCodeParser_ParseErrorCodeDefinition(t *testing.T) {
	testDefinition := []byte(`
[moduleA(module_code=06 output_path=./errors)] # moduleA is module name; (module_code=01 output_path=./errors) is optional
[[client_error]] 		  # client_error is keyword
ErrModuleAVar1 = "error msg 1"
ErrModuleAVar2 = "error msg 2"

[[server_error]] 		 # server_error is keyword

[moduleB]
[[server_error]] 		 # server_error is keyword
ErrModuleBVar1 = "error msg 2"`)
	expect := []config.ErrCodeModuleConfig{
		config.ErrCodeModuleConfig{
			OutputPath: "./errors",
			ModuleName: "moduleA",
			ModuleCode: "06",
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
			ModuleCode: "07",
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

func TestParseModuleArgs(t *testing.T) {
	line := "[moduleA(module_code=06)] # moduleA is module name; (module_code=01) is optional"
	expect := config.ErrCodeModuleConfig{
		ModuleName: "moduleA",
		ModuleCode: "06",
	}
	p := moduleErrorCodeParser{}
	name, args := p.parseModuleArgs(line)
	assert.Equal(t, expect.ModuleName, name)
	assert.Equal(t, expect.ModuleCode, args["module_code"])
}
