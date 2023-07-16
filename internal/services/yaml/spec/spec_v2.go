package spec

import (
	"fmt"
	"path"

	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/arch"
	"github.com/fe3dback/go-arch-lint/internal/models/common"
)

type (
	ArchV2Document struct {
		filePath common.Referable[string]

		reference                  common.Reference
		internalVersion            common.Referable[int]
		internalWorkingDir         common.Referable[string]
		internalVendors            archV2InternalVendors
		internalExclude            archV2InternalExclude
		internalExcludeFilesRegExp archV2InternalExcludeFilesRegExp
		internalComponents         archV2InternalComponents
		internalDependencies       archV2InternalDependencies
		internalCommonComponents   archV2InternalCommonComponents
		internalCommonVendors      archV2InternalCommonVendors

		V2Version            int                                    `yaml:"version" json:"version"`
		V2WorkDir            string                                 `yaml:"workdir" json:"workdir"`
		V2Allow              ArchV2Allow                            `yaml:"allow" json:"allow"`
		V2Vendors            map[arch.VendorName]ArchV2Vendor       `yaml:"vendors" json:"vendors"`
		V2Exclude            []string                               `yaml:"exclude" json:"exclude"`
		V2ExcludeFilesRegExp []string                               `yaml:"excludeFiles" json:"excludeFiles"`
		V2Components         map[arch.ComponentName]ArchV2Component `yaml:"components" json:"components"`
		V2Dependencies       map[arch.ComponentName]ArchV2Rules     `yaml:"deps" json:"deps"`
		V2CommonComponents   []string                               `yaml:"commonComponents" json:"commonComponents"`
		V2CommonVendors      []string                               `yaml:"commonVendors" json:"commonVendors"`
	}

	ArchV2Allow struct {
		reference              common.Reference
		internalDepOnAnyVendor common.Referable[bool]

		V2DepOnAnyVendor bool `yaml:"depOnAnyVendor" json:"depOnAnyVendor"`
	}

	ArchV2Vendor struct {
		reference           common.Reference
		internalImportPaths []common.Referable[models.Glob]

		V2ImportPaths stringsList `yaml:"in" json:"in"`
	}

	ArchV2Component struct {
		reference          common.Reference
		internalLocalPaths []common.Referable[models.Glob]

		V2LocalPaths stringsList `yaml:"in" json:"in"`
	}

	ArchV2Rules struct {
		reference              common.Reference
		internalMayDependOn    []common.Referable[string]
		internalCanUse         []common.Referable[string]
		internalAnyProjectDeps common.Referable[bool]
		internalAnyVendorDeps  common.Referable[bool]

		V2MayDependOn    []string `yaml:"mayDependOn" json:"mayDependOn"`
		V2CanUse         []string `yaml:"canUse" json:"canUse"`
		V2AnyProjectDeps bool     `yaml:"anyProjectDeps" json:"anyProjectDeps"`
		V2AnyVendorDeps  bool     `yaml:"anyVendorDeps" json:"anyVendorDeps"`
	}
)

type (
	archV2InternalVendors struct {
		reference common.Reference
		data      map[arch.VendorName]ArchV2Vendor
	}

	archV2InternalComponents struct {
		reference common.Reference
		data      map[arch.ComponentName]ArchV2Component
	}

	archV2InternalExclude struct {
		reference common.Reference
		data      []common.Referable[string]
	}

	archV2InternalExcludeFilesRegExp struct {
		reference common.Reference
		data      []common.Referable[string]
	}

	archV2InternalCommonVendors struct {
		reference common.Reference
		data      []common.Referable[string]
	}

	archV2InternalCommonComponents struct {
		reference common.Reference
		data      []common.Referable[string]
	}

	archV2InternalDependencies struct {
		reference common.Reference
		data      map[arch.ComponentName]ArchV2Rules
	}
)

// -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --

func (doc ArchV2Document) FilePath() common.Referable[string] {
	return doc.filePath
}

func (doc ArchV2Document) Reference() common.Reference {
	return doc.reference
}

func (doc ArchV2Document) Version() common.Referable[int] {
	return doc.internalVersion
}

func (doc ArchV2Document) WorkingDirectory() common.Referable[string] {
	return doc.internalWorkingDir
}

func (doc ArchV2Document) Options() arch.Options {
	return doc.V2Allow
}

func (doc ArchV2Document) ExcludedDirectories() arch.ExcludedDirectories {
	return doc.internalExclude
}

func (doc ArchV2Document) ExcludedFilesRegExp() arch.ExcludedFilesRegExp {
	return doc.internalExcludeFilesRegExp
}

func (doc ArchV2Document) Vendors() arch.Vendors {
	return doc.internalVendors
}

func (doc ArchV2Document) Components() arch.Components {
	return doc.internalComponents
}

func (doc ArchV2Document) CommonComponents() arch.CommonComponents {
	return doc.internalCommonComponents
}

func (doc ArchV2Document) CommonVendors() arch.CommonVendors {
	return doc.internalCommonVendors
}

func (doc ArchV2Document) Dependencies() arch.Dependencies {
	return doc.internalDependencies
}

func (doc ArchV2Document) applyReferences(resolve yamlDocumentPathResolver) ArchV2Document {
	doc.reference = resolve("$.version")

	// Version
	doc.internalVersion = common.NewReferable(
		doc.V2Version,
		resolve("$.version"),
	)

	// Working Directory
	actualWorkDirectory := "./" // fallback from version 1
	if doc.V2WorkDir != "" {
		actualWorkDirectory = doc.V2WorkDir
	}

	doc.internalWorkingDir = common.NewReferable(
		actualWorkDirectory,
		resolve("$.workdir"),
	)

	// Allow
	doc.V2Allow = doc.V2Allow.applyReferences(resolve)

	// Vendors
	vendors := make(map[string]ArchV2Vendor)
	for name, vendor := range doc.V2Vendors {
		vendors[name] = vendor.applyReferences(name, resolve)
	}
	doc.internalVendors = archV2InternalVendors{
		reference: resolve("$.vendors"),
		data:      vendors,
	}

	// Exclude
	excludedDirectories := make([]common.Referable[string], len(doc.V2Exclude))
	for ind, item := range doc.V2Exclude {
		excludedDirectories[ind] = common.NewReferable(
			item,
			resolve(fmt.Sprintf("$.exclude[%d]", ind)),
		)
	}

	doc.internalExclude = archV2InternalExclude{
		reference: resolve("$.exclude"),
		data:      excludedDirectories,
	}

	// ExcludeFilesRegExp
	excludedFiles := make([]common.Referable[string], len(doc.V2ExcludeFilesRegExp))
	for ind, item := range doc.V2ExcludeFilesRegExp {
		excludedFiles[ind] = common.NewReferable(
			item,
			resolve(fmt.Sprintf("$.excludeFiles[%d]", ind)),
		)
	}

	doc.internalExcludeFilesRegExp = archV2InternalExcludeFilesRegExp{
		reference: resolve("$.excludeFiles"),
		data:      excludedFiles,
	}

	// Components
	components := make(map[string]ArchV2Component)
	for name, component := range doc.V2Components {
		components[name] = component.applyReferences(name, doc.internalWorkingDir.Value, resolve)
	}
	doc.internalComponents = archV2InternalComponents{
		reference: resolve("$.components"),
		data:      components,
	}

	// Dependencies
	dependencies := make(map[string]ArchV2Rules)
	for name, rules := range doc.V2Dependencies {
		dependencies[name] = rules.applyReferences(name, resolve)
	}
	doc.internalDependencies = archV2InternalDependencies{
		reference: resolve("$.deps"),
		data:      dependencies,
	}

	// CommonComponents
	commonComponents := make([]common.Referable[string], len(doc.V2CommonComponents))
	for ind, item := range doc.V2CommonComponents {
		commonComponents[ind] = common.NewReferable(
			item,
			resolve(fmt.Sprintf("$.commonComponents[%d]", ind)),
		)
	}
	doc.internalCommonComponents = archV2InternalCommonComponents{
		reference: resolve("$.commonComponents"),
		data:      commonComponents,
	}

	// CommonVendors
	commonVendors := make([]common.Referable[string], len(doc.V2CommonVendors))
	for ind, item := range doc.V2CommonVendors {
		commonVendors[ind] = common.NewReferable(
			item,
			resolve(fmt.Sprintf("$.commonVendors[%d]", ind)),
		)
	}
	doc.internalCommonVendors = archV2InternalCommonVendors{
		reference: resolve("$.commonVendors"),
		data:      commonVendors,
	}

	return doc
}

// -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --

func (opt ArchV2Allow) Reference() common.Reference {
	return opt.reference
}

func (opt ArchV2Allow) IsDependOnAnyVendor() common.Referable[bool] {
	return opt.internalDepOnAnyVendor
}

func (opt ArchV2Allow) DeepScan() common.Referable[bool] {
	return common.NewEmptyReferable(false)
}

func (opt ArchV2Allow) applyReferences(resolve yamlDocumentPathResolver) ArchV2Allow {
	opt.reference = resolve("$.allow")

	opt.internalDepOnAnyVendor = common.NewReferable(
		opt.V2DepOnAnyVendor,
		resolve("$.allow.depOnAnyVendor"),
	)

	return opt
}

// -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --

func (v ArchV2Vendor) Reference() common.Reference {
	return v.reference
}

func (v ArchV2Vendor) ImportPaths() []common.Referable[models.Glob] {
	return v.internalImportPaths
}

func (v ArchV2Vendor) applyReferences(name arch.VendorName, resolve yamlDocumentPathResolver) ArchV2Vendor {
	v.reference = resolve(fmt.Sprintf("$.vendors.%s", name))

	for ind, importPath := range v.V2ImportPaths.list {
		yamlPath := fmt.Sprintf("$.vendors.%s.in", name)
		if v.V2ImportPaths.definedAsList {
			yamlPath = fmt.Sprintf("%s[%d]", yamlPath, ind)
		}

		v.internalImportPaths = append(v.internalImportPaths, common.NewReferable(
			models.Glob(importPath),
			resolve(yamlPath),
		))
	}

	return v
}

// -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --

func (c ArchV2Component) Reference() common.Reference {
	return c.reference
}

func (c ArchV2Component) RelativePaths() []common.Referable[models.Glob] {
	return c.internalLocalPaths
}

func (c ArchV2Component) applyReferences(
	name arch.ComponentName,
	workDirectory string,
	resolve yamlDocumentPathResolver,
) ArchV2Component {
	c.reference = resolve(fmt.Sprintf("$.components.%s", name))

	for ind, importPath := range c.V2LocalPaths.list {
		yamlPath := fmt.Sprintf("$.components.%s.in", name)
		if c.V2LocalPaths.definedAsList {
			yamlPath = fmt.Sprintf("%s[%d]", yamlPath, ind)
		}

		c.internalLocalPaths = append(c.internalLocalPaths, common.NewReferable(
			models.Glob(
				path.Clean(fmt.Sprintf("%s/%s",
					workDirectory,
					importPath,
				)),
			),
			resolve(yamlPath),
		))
	}

	return c
}

// -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --

func (rule ArchV2Rules) Reference() common.Reference {
	return rule.reference
}

func (rule ArchV2Rules) MayDependOn() []common.Referable[string] {
	return rule.internalMayDependOn
}

func (rule ArchV2Rules) CanUse() []common.Referable[string] {
	return rule.internalCanUse
}

func (rule ArchV2Rules) AnyProjectDeps() common.Referable[bool] {
	return rule.internalAnyProjectDeps
}

func (rule ArchV2Rules) AnyVendorDeps() common.Referable[bool] {
	return rule.internalAnyVendorDeps
}

func (rule ArchV2Rules) DeepScan() common.Referable[bool] {
	return common.NewEmptyReferable(false)
}

func (rule ArchV2Rules) applyReferences(name arch.ComponentName, resolve yamlDocumentPathResolver) ArchV2Rules {
	rule.reference = resolve(fmt.Sprintf("$.deps.%s", name))

	// --
	rule.internalAnyVendorDeps = common.NewReferable(
		rule.V2AnyVendorDeps,
		resolve(fmt.Sprintf("$.deps.%s.anyVendorDeps", name)),
	)

	// --
	rule.internalAnyProjectDeps = common.NewReferable(
		rule.V2AnyProjectDeps,
		resolve(fmt.Sprintf("$.deps.%s.anyProjectDeps", name)),
	)

	// --
	canUse := make([]common.Referable[string], len(rule.V2CanUse))
	for ind, item := range rule.V2CanUse {
		canUse[ind] = common.NewReferable(
			item,
			resolve(fmt.Sprintf("$.deps.%s.canUse[%d]", name, ind)),
		)
	}
	rule.internalCanUse = canUse

	// --
	mayDependOn := make([]common.Referable[string], len(rule.V2MayDependOn))
	for ind, item := range rule.V2MayDependOn {
		mayDependOn[ind] = common.NewReferable(
			item,
			resolve(fmt.Sprintf("$.deps.%s.mayDependOn[%d]", name, ind)),
		)
	}
	rule.internalMayDependOn = mayDependOn

	// --
	return rule
}

// -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --

func (a archV2InternalDependencies) Reference() common.Reference {
	return a.reference
}

func (a archV2InternalDependencies) Map() map[arch.ComponentName]arch.DependencyRule {
	res := make(map[arch.ComponentName]arch.DependencyRule)
	for name, rules := range a.data {
		res[name] = rules
	}
	return res
}

func (a archV2InternalCommonComponents) Reference() common.Reference {
	return a.reference
}

func (a archV2InternalCommonComponents) List() []common.Referable[string] {
	return a.data
}

func (a archV2InternalCommonVendors) Reference() common.Reference {
	return a.reference
}

func (a archV2InternalCommonVendors) List() []common.Referable[string] {
	return a.data
}

func (a archV2InternalExcludeFilesRegExp) Reference() common.Reference {
	return a.reference
}

func (a archV2InternalExcludeFilesRegExp) List() []common.Referable[string] {
	return a.data
}

func (a archV2InternalExclude) Reference() common.Reference {
	return a.reference
}

func (a archV2InternalExclude) List() []common.Referable[string] {
	return a.data
}

func (a archV2InternalComponents) Reference() common.Reference {
	return a.reference
}

func (a archV2InternalComponents) Map() map[arch.ComponentName]arch.Component {
	res := make(map[arch.ComponentName]arch.Component)
	for name, component := range a.data {
		res[name] = component
	}
	return res
}

func (a archV2InternalVendors) Reference() common.Reference {
	return a.reference
}

func (a archV2InternalVendors) Map() map[arch.VendorName]arch.Vendor {
	res := make(map[arch.VendorName]arch.Vendor)
	for name, vendor := range a.data {
		res[name] = vendor
	}
	return res
}
