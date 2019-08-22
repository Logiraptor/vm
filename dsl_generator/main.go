package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"html/template"
	"os"
	"os/exec"
	"strings"
)

const tmpl = `package {{.Package}}

{{- $name := .InstructionName | printf "%sBuilder" }}

type {{$name}} struct {
  Instructions []{{.InstructionType}}
}

{{range .Instructions -}}
func (b *{{$name}}) {{.Name}}({{range .Args -}}
{{.Name | lower}} {{.Type}},
{{- end -}}
) {
  b.Instructions = append(b.Instructions, {{ .Type }}{
{{- range .Args}}
	{{.Name}}: {{.Name | lower}},
{{- end}}
  })
}

{{end}}
`

type Package struct {
	Package         string
	InstructionName string
	InstructionType string
	Instructions    []Instruction
}

type Instruction struct {
	Name string
	Type string
	Args []Arg
}

type Arg struct {
	Name string
	Type string
}

func main() {

	typ := flag.String("type", "", "Type to generate a dsl for")

	dir, _ := os.Getwd()
	flag.Parse()

	if *typ == "" {
		flag.Usage()
		os.Exit(1)
	}

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, func(info os.FileInfo) bool {
		return !strings.HasSuffix(info.Name(), "_test.go")
	}, 0)
	if err != nil {
		panic(err)
	}
	if len(pkgs) != 1 {
		panic(fmt.Sprintf("Cannot handle multiple packages declared in the same directory"))
	}

	var instructions []Instruction
	var pkgName string
	for name, pkg := range pkgs {
		pkgName = name
		instructions = processPackage(*typ, fset, pkg)
	}

	generateDsl(typ, pkgName, instructions)

	cmd := exec.Command("goimports", "-w", strings.ToLower(*typ)+"_dsl.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func generateDsl(typ *string, pkgName string, instructions []Instruction) {
	outFile, err := os.Create(strings.ToLower(*typ) + "_dsl.go")
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	tmpl, err := template.New("dsl").Funcs(template.FuncMap{
		"lower": strings.ToLower,
	}).Parse(tmpl)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(outFile, Package{
		Package:         pkgName,
		InstructionName: *typ,
		InstructionType: *typ,
		Instructions:    instructions,
	})
	if err != nil {
		panic(err)
	}
}

func processPackage(typ string, fset *token.FileSet, pkg *ast.Package) []Instruction {
	var instructions []Instruction

	for _, f := range pkg.Files {
		for _, decl := range f.Decls {
			if v, ok := decl.(*ast.GenDecl); ok {
				for _, spec := range v.Specs {
					if ts, ok := spec.(*ast.TypeSpec); ok {
						if strct, ok := ts.Type.(*ast.StructType); ok {
							var isInstruction bool
							var args []Arg
							for _, field := range strct.Fields.List {
								if len(field.Names) == 0 && tts(fset, field.Type) == typ {
									isInstruction = true
								}
								for _, name := range field.Names {
									args = append(args, Arg{
										Name: name.String(),
										Type: tts(fset, field.Type),
									})
								}

							}
							if isInstruction {
								instructions = append(instructions, Instruction{
									Name: ts.Name.Name[:len(ts.Name.Name)-len(typ)],
									Type: ts.Name.Name,
									Args: args,
								})
							}
						}
					}
				}
			}
		}
	}

	return instructions
}

func tts(fset *token.FileSet, typ ast.Expr) string {
	var buf = new(bytes.Buffer)
	printer.Fprint(buf, fset, typ)
	return buf.String()
}
