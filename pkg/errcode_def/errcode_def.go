// errcode_def contain self defined error code definition format.
package errcode_def

import (
	"bufio"
	"github.com/RivenZoo/errcodegen/pkg/config"
	"fmt"
	"io"
	re "regexp"
	"strings"
)

type ModuleErrorCodeParser interface {
	ParseErrorCodeDefinition(r io.Reader) ([]config.ErrCodeModuleConfig, error)
}

func NewModuleErrorCodeParser() ModuleErrorCodeParser {
	return moduleErrorCodeParser{}
}

type moduleErrorCodeParser struct {
}

// ParseErrorCodeDefinition parse self defined module error code config format.
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
func (p moduleErrorCodeParser) ParseErrorCodeDefinition(r io.Reader) ([]config.ErrCodeModuleConfig, error) {
	ret := []config.ErrCodeModuleConfig{}
	scanner := bufio.NewScanner(r)

	lctx := newLineScannerContext(scanner)
	next := true
	var err error

	for {
		if next {
			if !lctx.Iterator().Next() {
				break
			}
		}
		line := lctx.Line()
		if line == "" {
			next = true
			continue
		}
		ltype := lctx.LineType()
		if ltype == moduleDefine {
			var modules []config.ErrCodeModuleConfig
			modules, next, err = p.parseModuleConfig(lctx)
			if err != nil {
				return nil, err
			}
			if len(modules) > 0 {
				ret = append(ret, modules...)
			}
		} else {
			next = true
		}
	}
	if err = lctx.Iterator().Err(); err != nil {
		return nil, err
	}
	config.PatchModules(ret)
	return ret, nil
}

// parseModuleConfig parse module config from first module definition line.
func (p moduleErrorCodeParser) parseModuleConfig(lctx lineContext) ([]config.ErrCodeModuleConfig, bool, error) {
	modulesConfig := []config.ErrCodeModuleConfig{
		config.ErrCodeModuleConfig{},
	}
	line := lctx.Line()

	var args map[string]string

	module := &modulesConfig[0]
	module.ModuleName, args = p.parseModuleArgs(line)
	if module.ModuleName == "" {
		return nil, true, fmt.Errorf("line %s wrong module config format", line)
	}
	if v, ok := args["module_code"]; ok {
		module.ModuleCode = v
	}

	var err error
	next := true

	for {
		if next {
			if !lctx.Iterator().Next() {
				break
			}
		}
		line := lctx.Line()
		if line == "" {
			next = true
			continue
		}
		ltype := lctx.LineType()
		switch ltype {
		case moduleDefine:
			return modulesConfig, false, nil
		case clientVariableDefine:
			var clientErrConfig []config.ErrCodeVariableConfig
			clientErrConfig, next, err = p.parseVariableConfig(lctx)
			if err != nil {
				return nil, true, err
			}
			if len(clientErrConfig) > 0 {
				module.ClientCodes = append(module.ClientCodes, clientErrConfig...)
			}
		case serverVariableDefine:
			var serverErrConfig []config.ErrCodeVariableConfig
			serverErrConfig, next, err = p.parseVariableConfig(lctx)
			if err != nil {
				return nil, true, err
			}
			if len(serverErrConfig) > 0 {
				module.ServerCodes = append(module.ServerCodes, serverErrConfig...)
			}
		default:
			// skip
			next = true
		}
	}
	return modulesConfig, false, nil
}

// parseVariableConfig parse error code definition from line client_error or server_error.
func (p moduleErrorCodeParser) parseVariableConfig(lctx lineContext) ([]config.ErrCodeVariableConfig, bool, error) {
	variablesConfig := []config.ErrCodeVariableConfig{}

	next := true

	// skip line client_error or server_error
	for {
		if next {
			if !lctx.Iterator().Next() {
				break
			}
		}
		line := lctx.Line()
		if line == "" {
			next = true
			continue
		}
		ltype := lctx.LineType()
		switch ltype {
		case moduleDefine, serverVariableDefine, clientVariableDefine:
			return variablesConfig, false, nil
		case otherDefine:
			args := p.parseSingleArgs(line)
			for k, v := range args {
				variablesConfig = append(variablesConfig, config.ErrCodeVariableConfig{
					Name: k,
					Msg:  v,
				})
			}
		default:
			next = true
		}
	}
	return variablesConfig, false, nil
}

// parseModuleArgs parse module definition [A(module_code=01)].
func (p moduleErrorCodeParser) parseModuleArgs(line string) (moduleName string, args map[string]string) {
	pattern := re.MustCompile(`\[([a-zA-Z0-9_]+)(\(.*\))?\]`)
	matches := pattern.FindStringSubmatch(line)
	if len(matches) < 2 {
		return
	}
	moduleName = trimRawLine(matches[1])
	if len(matches) == 3 && matches[2] != "" {
		args = p.parseSingleArgs(matches[2])
	}
	return
}

func (p moduleErrorCodeParser) parseSingleArgs(line string) (args map[string]string) {
	args = map[string]string{}

	kvs := strings.Split(line, "=")
	if len(kvs) != 2 {
		return
	}
	k := trimRawLine(kvs[0])
	v := trimRawLine(kvs[1])
	if k != "" && v != "" {
		args[trimQuotation(k)] = trimQuotation(v)
	}
	return
}

func trimQuotation(line string) string {
	return strings.Trim(line, `'"`)
}
