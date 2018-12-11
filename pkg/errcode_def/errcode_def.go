// errcode_def contain self defined error code definition format.
package errcode_def

import "errcodegen/pkg/config"

type ModuleErrorCodeParser interface {
	ParseDefinitionFile(f string) ([]config.ErrCodeModuleConfig, error)
}

func NewModuleErrorCodeParser() ModuleErrorCodeParser {
	return moduleErrorCodeParser{}
}

type moduleErrorCodeParser struct {
}

// ParseDefinitionFile parse self defined module error code config format.
// Format:
// [moduleA(module_code=01)] # moduleA is module name; (module_code=01) is optional
// [[client_error]] 		 # client_error is keyword
// ErrModuleAVar1 = "error msg"
//
// [[server_error]] 		 # server_error is keyword
//
// [moduleB]
// [[client_error]] 		 # client_error is keyword
// ErrModuleBVar1 = "error msg"
//
// Return:
// []config.ErrCodeModuleConfig{
//	config.ErrCodeModuleConfig{
//		ModuleCode: "01",
//		ModuleName: "moduleA",
//		ClientCodes: []config.ErrCodeVariableConfig{
//			config.ErrCodeVariableConfig{
//				Name:      "ErrModuleAVar1",
//				ErrNumber: "01",
//				Msg:       "error msg",
//			},
//		},
//	},
//	config.ErrCodeModuleConfig{
//		ModuleCode: "02",
//		ModuleName: "moduleB",
//		ClientCodes: []config.ErrCodeVariableConfig{
//			config.ErrCodeVariableConfig{
//				Name:      "ErrModuleBVar1",
//				ErrNumber: "01",
//				Msg:       "error msg",
//			},
//		},
//	},
//}
func (p moduleErrorCodeParser) ParseDefinitionFile(f string) ([]config.ErrCodeModuleConfig, error) {
	return nil, nil
}
