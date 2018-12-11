package config

import (
	"errcodegen/pkg/log"
	"fmt"
	"os"
	"strconv"
)

var (
	defaultErrCodeCommonConfig = ErrCodeCommonConfig{
		PkgName:          "errorcodes",
		ClientCodePrefix: "4",
		ServerCodePrefix: "5",
		AppCode:          "100",
		NewErrorFuncPkg:  "errors",
		NewErrorFunc:     "New",
	}
	moduleCodePattern   = "%02d"
	variableCodePattern = "%02d"
)

func PatchConfig(conf *ErrCodeConfig) {
	if conf.PkgName == "" {
		conf.PkgName = defaultErrCodeCommonConfig.PkgName
	}
	if conf.ClientCodePrefix == "" {
		conf.ClientCodePrefix = defaultErrCodeCommonConfig.ClientCodePrefix
	}
	if conf.ServerCodePrefix == "" {
		conf.ServerCodePrefix = defaultErrCodeCommonConfig.ServerCodePrefix
	}
	if conf.AppCode == "" {
		conf.AppCode = defaultErrCodeCommonConfig.AppCode
	}
	if conf.NewErrorFunc == "" && conf.NewErrorFuncPkg == "" {
		conf.NewErrorFunc = defaultErrCodeCommonConfig.NewErrorFunc
		conf.NewErrorFuncPkg = defaultErrCodeCommonConfig.NewErrorFuncPkg
	}
	patchModules(conf.Modules)
}

func patchModules(modules []ErrCodeModuleConfig) {
	lastCode := 0
	for i := range modules {
		if modules[i].ModuleCode == "" {
			lastCode += 1
			modules[i].ModuleCode = fmt.Sprintf(moduleCodePattern, lastCode)
		} else {
			i, err := strconv.ParseInt(modules[i].ModuleCode, 10, 32)
			if err != nil {
				log.Error("module code %s wrong format", modules[i].ModuleCode)
				os.Exit(-1)
			}
			lastCode = int(i)
		}
		patchCodeVariables(modules[i].ClientCodes)
		patchCodeVariables(modules[i].ServerCodes)
	}
}

func patchCodeVariables(codes []ErrCodeVariableConfig) {
	lastNum := 0
	for i := range codes {
		if codes[i].ErrNumber == "" {
			lastNum += 1
			codes[i].ErrNumber = fmt.Sprint(variableCodePattern, lastNum)
		} else {
			i, err := strconv.ParseInt(codes[i].ErrNumber, 10, 32)
			if err != nil {
				log.Error("variable code %s wrong format", codes[i].ErrNumber)
				os.Exit(-1)
			}
			lastNum = int(i)
		}
	}
}

func CheckConfig(conf *ErrCodeConfig) {
	if conf.ClientCodePrefix == "" || conf.ServerCodePrefix == "" || conf.AppCode == "" || conf.NewErrorFunc == "" {
		log.Error("wrong config %v", conf)
		os.Exit(-1)
	}
	if len(conf.Modules) == 0 {
		log.Info("no modules error code config")
		return
	}
	checkModules(conf.Modules)
}

func checkModules(modules []ErrCodeModuleConfig) {
	existsModule := map[string]struct{}{}
	existsVariables := map[string]struct{}{}

	lastCode := modules[0].ModuleCode
	for i := range modules {
		if modules[i].ModuleName == "" {
			log.Error("no module name")
			os.Exit(-1)
		}
		if _, ok := existsModule[modules[i].ModuleName]; ok {
			log.Error("module %s already exist", modules[i].ModuleName)
			os.Exit(-1)
		}
		existsModule[modules[i].ModuleName] = struct{}{}

		if i > 0 && modules[i].ModuleCode <= lastCode {
			log.Error("module %v code less than last one", modules[i])
			os.Exit(-1)
		}
		lastCode = modules[i].ModuleCode
		checkVariableCodes(modules[i].ClientCodes, existsVariables)
		checkVariableCodes(modules[i].ServerCodes, existsVariables)
	}
}

func checkVariableCodes(codes []ErrCodeVariableConfig, existsVariables map[string]struct{}) {
	lastCode := codes[0].ErrNumber
	for i := range codes {
		if codes[i].Name == "" {
			log.Error("no variable name")
			os.Exit(-1)
		}
		if _, ok := existsVariables[codes[i].Name]; ok {
			log.Error("variable %s already exist", codes[i].Name)
			os.Exit(-1)
		}
		existsVariables[codes[i].Name] = struct{}{}

		if i > 0 && codes[i].ErrNumber <= lastCode {
			log.Error("variable %v code less than last one", codes[i])
			os.Exit(-1)
		}
		lastCode = codes[i].ErrNumber
	}
}
