package gen

import (
	"errors"
	"fmt"
	"github.com/RivenZoo/errcodegen/pkg/config"
	"github.com/RivenZoo/errcodegen/pkg/log"
	"github.com/dave/jennifer/jen"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

type CodeGenerator interface {
	// GenerateModuleErrorCode generate error code and save it to module file.
	GenerateModuleErrorCode(commonConfig *config.ErrCodeCommonConfig, moduleConfig *config.ErrCodeModuleConfig) (error)
}

func NewCodeGenerator() CodeGenerator {
	return moduleErrorCodeGenerator{}
}

type moduleErrorCodeGenerator struct {
}

func (g moduleErrorCodeGenerator) GenerateModuleErrorCode(commonConfig *config.ErrCodeCommonConfig,
	moduleConfig *config.ErrCodeModuleConfig) error {
	if commonConfig == nil || moduleConfig == nil {
		return errors.New("no config param")
	}
	var f *jen.File

	pkgName := path.Base(moduleConfig.OutputPath)
	if pkgName != "" && pkgName != "." {
		f = jen.NewFilePath(moduleConfig.OutputPath)
	} else {
		f = jen.NewFile(commonConfig.PkgName)
	}
	f.PackageComment("Code generated by errcodegen. DO NOT EDIT.")
	f.Commentf("This file contains module %s error codes.", moduleConfig.ModuleName)
	f.Line()

	importName := filepath.Base(commonConfig.NewErrorFuncPkg)
	f.ImportName(commonConfig.NewErrorFuncPkg, importName)

	f.Commentf("Define error code for errors from client side.")
	g.generateModuleVariables(f, commonConfig, commonConfig.ClientCodePrefix,
		moduleConfig.ModuleCode, moduleConfig.ClientCodes)

	f.Commentf("Define error code for errors from server side.")
	g.generateModuleVariables(f, commonConfig, commonConfig.ServerCodePrefix,
		moduleConfig.ModuleCode, moduleConfig.ServerCodes)

	fname := fixGoSuffix(moduleConfig.ModuleName)

	outputPath := fname
	if moduleConfig.OutputPath != "" {
		outputPath = path.Join(moduleConfig.OutputPath, fname)
		os.MkdirAll(moduleConfig.OutputPath, os.FileMode(0755))
	}

	err := f.Save(outputPath)
	if err != nil {
		log.Error("save file %s error %v", outputPath, err)
		return err
	}
	return nil
}

func (g moduleErrorCodeGenerator) generateModuleVariables(f *jen.File,
	commonConfig *config.ErrCodeCommonConfig, prefix string,
	moduleCode string, variablesConfig []config.ErrCodeVariableConfig) error {
	importName := filepath.Base(commonConfig.NewErrorFuncPkg)
	snip := jen.Empty()
	for i := range variablesConfig {
		codeConfig := &variablesConfig[i]
		variableName := capitalizeVariableName(codeConfig.Name)
		errCode, err := makeErrorCode(commonConfig.AppCode, prefix,
			moduleCode, codeConfig.ErrNumber)
		if err != nil {
			return err
		}
		args := makeNewErrorFuncParams(importName, commonConfig.NewErrorFunc,
			errCode, codeConfig.Msg)
		snip.Id(variableName).
			Op("=").Qual(commonConfig.NewErrorFuncPkg, commonConfig.NewErrorFunc).
			Call(args...).Line()
	}
	f.Var().Parens(snip)
	return nil
}

func fixGoSuffix(fname string) string {
	if filepath.Ext(fname) == ".go" {
		return fname[:len(fname)-3] + "_gen.go"
	}
	return fname + "_gen.go"
}

func capitalizeVariableName(name string) string {
	return strings.Title(name)
}

func makeErrorCode(appCode, prefix, moduleCode, errNum string) (int, error) {
	s := fmt.Sprintf("%s%s%s%s", prefix, appCode, moduleCode, errNum)
	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		log.Error("parse code %s error", s)
		return 0, err
	}
	return int(i), nil
}

func makeNewErrorFuncParams(importName string, newErrorFunc string, errCode int, msg string) []jen.Code {
	funcCall := importName + "." + newErrorFunc
	switch funcCall {
	case "errors.New":
		return []jen.Code{jen.Lit(fmt.Sprintf("[%d] %s", errCode, msg))}
	case "fmt.Errorf":
		return []jen.Code{jen.Lit(msg), jen.Lit(errCode)}
	default:
		// func funcNew(errCode int, msg string)
		return []jen.Code{jen.Lit(errCode), jen.Lit(msg)}
	}
}
