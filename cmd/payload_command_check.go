package cmd

import (
	"github.com/fe3dback/go-arch-lint/checker"
	"github.com/fe3dback/go-arch-lint/spec/annotated_validator"
)

type payloadTypeCommandCheck struct {
	ExecutionWarnings      []annotated_validator.YamlAnnotatedWarning `json:"execution_warnings"`
	ExecutionError         string                                     `json:"execution_error"`
	ArchHasWarnings        bool                                       `json:"arch_has_warnings"`
	ArchWarningsDeps       []checker.WarningDep                       `json:"arch_warnings_deps"`
	ArchWarningsNotMatched []checker.WarningNotMatched                `json:"arch_warnings_not_matched"`
}