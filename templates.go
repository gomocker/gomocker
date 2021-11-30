package main

var output = `// Code generated by gomocker. DO NOT EDIT.
// versions:
//     gomocker {{.Version}}

package {{.Package}}
{{if .Imports}}
import ({{ range $short, $long := .Imports }}
{{$short}} {{$long}} {{end}}
) {{end}}

// Impls
type ( {{ range $tmpl := .Tmpls }}
{{$tmpl.StructName}}Impl struct { behavior {{$tmpl.MockName}}Behavior } {{end}}
)

var ErrMock = errors.New("")

// Check
var ({{ range $tmpl := .Tmpls }}
_ {{$tmpl.OriginalInterface}} = &{{$tmpl.StructName}}Impl {} {{end}}
)

{{ range $tmpl := .Tmpls }}
// New{{$tmpl.MockName}} creates mocked implementation of {{$tmpl.OriginalInterface}} interface
func New{{$tmpl.MockName}}Mock(behavior {{$tmpl.MockName}}Behavior) {{$tmpl.OriginalInterface}} {
	return &{{$tmpl.StructName}}Impl {
		behavior: behavior,
	}
}
{{end}}

{{ range $tmpl := .Tmpls}}
type {{$tmpl.MockName}}Behavior struct { {{range $field := $tmpl.StructFields}}
    {{$field}} {{end}}
}
{{end}}
{{ range $tmpl := .Tmpls}}
{{ range $method := $tmpl.Methods }}
func (aeiouy *{{ $tmpl.StructName }}Impl) {{$method.WithTypes}} { {{if $method.DoesReturn}}
         if aeiouy.behavior.{{ $method.Name }} != nil {
              return aeiouy.behavior.{{$method.WithoutTypes}}
         }
         return {{ else }}
         if aeiouy.behavior.{{ $method.Name }} != nil {
              aeiouy.behavior.{{$method.WithoutTypes}}
         } {{ end }}
}
{{ end }}
{{ end }}`
