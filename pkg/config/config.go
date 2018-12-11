package config

import "github.com/BurntSushi/toml"

var (
	config ErrCodeConfig
)

type ErrCodeCommonConfig struct {
	PkgName          string `json:"pkg_name"`
	ClientCodePrefix string `json:"client_code_prefix"`
	ServerCodePrefix string `json:"server_code_prefix"`
	AppCode          string `json:"app_code"`
	NewErrorFuncPkg  string `json:"new_error_func_pkg"` // pkg like fmt or errors
	NewErrorFunc     string `json:"new_error_func"`     // func like `New` from errors or `Errorf` from fmt
}

type ErrCodeModuleConfig struct {
	ModuleName  string                  `json:"module_name"`
	ModuleCode  string                  `json:"module_code"`
	ClientCodes []ErrCodeVariableConfig `json:"client_codes"`
	ServerCodes []ErrCodeVariableConfig `json:"server_codes"`
}

type ErrCodeVariableConfig struct {
	Name      string `json:"name"`
	ErrNumber string `json:"err_number"`
	Msg       string `json:"msg"`
}

type ErrCodeConfig struct {
	ErrCodeCommonConfig
	Modules []ErrCodeModuleConfig `json:"modules"`
}

func GetConfig() *ErrCodeConfig {
	return &config
}

func MustInit(configFile string) {
	meta, err := toml.DecodeFile(configFile, &config)
	if err != nil {
		panic(err)
	}
	_ = meta
}

func ParseModuleConfig() error {
	return nil
}
