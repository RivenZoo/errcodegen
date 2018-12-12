// This file is used to trigger go generate. DO NOT ADD OTHER CODE HERE.
package example

//go:generate errcodegen --err_func=Errorf --err_func_pkg=fmt --pkg=example --err_def=error_code_def.conf
