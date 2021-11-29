package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go/parser"
	"go/token"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/armantarkhanian/gotype"

	"github.com/iancoleman/strcase"
	"golang.org/x/mod/modfile"
	"golang.org/x/tools/imports"
)

var (
	configFile = "gomocker.json"
	version    = "v1.1.0"
)

type Config struct {
	Package string              `json:"package"`
	Output  string              `json:"output"`
	Mocks   map[string][]string `json:"mocks"`
	Imports map[string]string   `json:"imports"`
}

var configExample string = `{
    "package": "golang_package_name",

    "output": "gomocker_output.go",

    "mocks": {
        "io": [
            "Reader",
            "Writer",
            "ReadWriter"
        ],
        "math/rand": [ "Source" ]
    },

    "imports": {
        "io": "io",
        "rand": "math/rand"
    }
}
`

var helpMessage string = `Gomocker is a tool for generating mocked interfaces in Go.

Usage:
	gomocker              run gomocker
	gomocker <command>    run command

The commands are:

	touch       create example "gomocker.json" config file
	version     print Gomocker version
	help        print this help message
`

func main() {
	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "version":
			fmt.Println("gomocker version", version)
		case "touch":
			f, err := os.OpenFile(configFile, os.O_RDONLY, 0666)
			if err != nil {
				if errors.Is(err, fs.ErrNotExist) {
					if err := ioutil.WriteFile(configFile, []byte(configExample), 0666); err != nil {
						fmt.Println(err)
						return
					}
					fmt.Println("gomocker.json file was created.")
				} else {
					fmt.Println(err)
				}
				return
			}
			defer f.Close()
			fmt.Println("gomocker.json file already existed")
		case "help", "--help", "-h":
			fmt.Println(helpMessage)
			fmt.Println("")
			fmt.Println(`For mor information, see https://github.com/gomocker/gomocker`)
		}
		return
	}

	packageName, err := packageName()
	if err != nil {
		fmt.Println(err)
		return
	}
	moduleName, err := moduleName()
	if err != nil {
		fmt.Println(err)
		return
	}

	b, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	var config Config

	if err := json.Unmarshal(b, &config); err != nil {
		fmt.Println(err)
		return
	}
	t, err := template.New("").Parse(output)
	if err != nil {
		fmt.Println(err)
		return
	}
	var data Data

	data.Package = packageName
	data.Imports = make(map[string]string)

	aliasInterfacePackage := ""
	aliasInterfaceName := ""

	for pkg := range config.Mocks {
		for i := 0; i < len(config.Mocks[pkg]); i++ {
			var interfacePackage string
			var interfaceName string

			if aliasInterfaceName != "" {
				interfaceName = aliasInterfaceName
				interfacePackage = strings.Trim(aliasInterfacePackage, "/")
			} else {
				interfaceName = config.Mocks[pkg][i]
				interfacePackage = strings.Trim(pkg, "/")
			}

			searchList, err := gotype.GenerateTypesFromSpecs(gotype.TypeSpec{
				PackagePath: interfacePackage,
				Name:        interfaceName,
			})
			if err != nil {
				fmt.Println(err)
				return
			}
			if len(searchList) == 0 {
				fmt.Println("Nothing found")
				return
			}

			interfaceType := searchList[0]

			if interfaceType.InterfaceType == nil && interfaceType.QualType == nil {
				continue
			}

			if interfaceType.QualType != nil {
				// if interface is an alias to another interface
				aliasInterfacePackage = interfaceType.QualType.Package
				aliasInterfaceName = interfaceType.QualType.Name
				i--
				continue
			}

			// if this is true interface
			if aliasInterfacePackage != "" && i+1 != len(config.Mocks[pkg]) {
				i++
			}

			aliasInterfacePackage = ""
			aliasInterfaceName = ""

			interfaceName = config.Mocks[pkg][i]
			interfacePackage = strings.Trim(pkg, "/")

			tmplData := TemplateStruct{
				StructFields: []string{},
			}

			if strings.HasPrefix(interfacePackage, moduleName+"/") {
				interfacePackage = strings.TrimPrefix(interfacePackage, moduleName+"/")
			}

			packageBase := filepath.Base(interfacePackage)
			if packageBase == packageName || packageBase == moduleName {
				packageBase = ""
			}

			tmplData.OriginalInterface += packageBase
			if packageBase != "" {
				tmplData.OriginalInterface = packageBase + "."
			}
			tmplData.OriginalInterface += interfaceName

			tmplData.MockName = strcase.ToCamel(interfaceName)
			tmplData.StructName = strcase.ToLowerCamel(tmplData.MockName)

			for _, m := range interfaceType.InterfaceType.Methods {
				for _, t := range m.Func.Inputs {
					short, long := t.Type.GetImportString()
					if short != "" {
						data.Imports[short] = fmt.Sprintf("%q", long)
					}
				}

				defaultReturnStmt := ""
				for i, t := range m.Func.Outputs {
					short, long := t.Type.GetImportString()
					if short != "" {
						data.Imports[short] = fmt.Sprintf("%q", long)
					}
					_, defValue := t.Type.Default(short)

					if defValue == "struct{}" {
						if short != packageName {
							defValue = short + "."
						}
						defValue += t.Type.String(short) + "{}"
					}
					defaultReturnStmt += defValue
					if i+1 != len(m.Func.Outputs) {
						defaultReturnStmt += ", "
					}
				}

				anonFunc := m.Func.String(moduleName)

				structField := m.Name + " " + anonFunc
				tmplData.StructFields = append(tmplData.StructFields, structField)

				tmplData.Methods = append(tmplData.Methods, Method{
					Name:              m.Name,
					WithTypes:         m.Name + anonFunc[4:],
					WithoutTypes:      m.Name + m.Func.StringWithoutTypes(moduleName)[4:],
					DefaultReturnStmt: defaultReturnStmt,
					DoesReturn:        len(m.Func.Outputs) > 0,
				})
			}
			data.Tmpls = append(data.Tmpls, tmplData)
		}
	}
	if len(data.Tmpls) == 0 {
		fmt.Println("Nothing to mock")
		return
	}

	if config.Package != "" {
		data.Package = config.Package
	}

	for short, long := range config.Imports {
		data.Imports[short] = fmt.Sprintf("%q", long)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		fmt.Println(err)
		return
	}
	b, err = imports.Process("", buf.Bytes(), nil)
	if err != nil {
		fmt.Println("Go imports organize error:", err)
		return
	}

	if config.Output != "" {
		err = ioutil.WriteFile(config.Output, b, 0666)
		if err != nil {
			fmt.Println("\nPlease organize imports by yourself")
		}
		fmt.Println("ok")
	} else {
		fmt.Print(string(b))
	}
}

type Data struct {
	Package string
	Imports map[string]string
	Tmpls   []TemplateStruct
}

type Method struct {
	Name              string
	WithTypes         string
	WithoutTypes      string
	DefaultReturnStmt string
	DoesReturn        bool
}

type TemplateStruct struct {
	MockName          string
	StructFields      []string
	Methods           []Method
	OriginalInterface string
	StructName        string
	DoesReturn        map[string]bool
}

func packageName() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	var packageName string

	if err := filepath.Walk(wd, func(path string, info fs.FileInfo, err error) error {
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		if packageName != "" {
			return nil
		}

		astFile, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.PackageClauseOnly)
		if err != nil {
			return err
		}
		if astFile.Name == nil {
			return fmt.Errorf("no package name found")
		}
		packageName = astFile.Name.Name
		return nil
	}); err != nil {
		return "", err
	}
	return packageName, nil
}

func moduleName() (string, error) {
	modFile, err := findGoModFile()
	if err != nil {
		return "", err
	}
	f, err := os.Open(modFile)
	if err != nil {
		return "", errors.New("can't open go.mod file")
	}
	defer f.Close()

	goModFile, err := os.Open(modFile)
	if err != nil {
		return "", fmt.Errorf("cannot open go.mod file: %w", err)
	}
	defer goModFile.Close()

	goModBytes, err := ioutil.ReadAll(goModFile)
	if err != nil {
		return "", fmt.Errorf("cannot read go.mod file: %w", err)
	}

	moduleFile, err := modfile.Parse("go.mod", goModBytes, nil)
	if err != nil {
		return "", fmt.Errorf("cannot parse go.mod file: %w", err)
	}

	return moduleFile.Module.Mod.Path, nil
}

func findGoModFile() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("cannot get current working dir: %w", err)
	}

	maxDepth := 15
	for i := 0; i < maxDepth; i++ {
		goModPath := path.Join(wd, "go.mod")
		if _, err := os.Stat(goModPath); errors.Is(err, os.ErrNotExist) {
			if wd == "/" {
				break
			}
			wd = path.Join(wd, "..")
			continue
		}

		return goModPath, nil
	}

	return "", fmt.Errorf("no go.mod file found")
}
