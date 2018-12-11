package gen

import (
	"errcodegen/pkg/config"
	"github.com/dave/jennifer/jen"
)

type CodeGenerator interface {
	GenerateModuleErrorCode(commonConfig *config.ErrCodeCommonConfig, moduleConfig *config.ErrCodeModuleConfig) ([]byte, error)
}

func NewCodeGenerator() CodeGenerator {
	return moduleErrorCodeGenerator{}
}

type moduleErrorCodeGenerator struct {
}

func (g moduleErrorCodeGenerator) GenerateModuleErrorCode(commonConfig *config.ErrCodeCommonConfig,
	moduleConfig *config.ErrCodeModuleConfig) ([]byte, error) {
	f := jen.NewFilePathName(commonConfig.PkgName, moduleConfig.ModuleName)
	f.Func().Id("main").Params().Block(
		jen.Qual("a.b/c", "Foo").Call(),
	)
	f.Save(moduleConfig.ModuleName)
	return nil, nil
}
