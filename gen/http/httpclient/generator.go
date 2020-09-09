package httpclient

import (
	"fmt"
	"sort"
	"strings"

	"github.com/RussellLuo/kok/gen/util/generator"
	"github.com/RussellLuo/kok/gen/util/misc"
	"github.com/RussellLuo/kok/gen/util/openapi"
	"github.com/RussellLuo/kok/gen/util/reflector"
)

var (
	template = `// Code generated by kok; DO NOT EDIT.
// github.com/RussellLuo/kok

package {{.Result.PkgName}}

import (
	"encoding/json"
	"net/http"
	"strconv"
	httpcodec "github.com/RussellLuo/kok/pkg/codec/httpv2"

	{{- range .Result.Imports}}
	"{{.}}"
	{{- end}}
)

type HTTPClient struct {
	codecs     httpcodec.Codecs
	httpClient *http.Client
	scheme     string
	host       string
	pathPrefix string
}

func NewHTTPClient(codecs httpcodec.Codecs, httpClient *http.Client, baseURL string) (*HTTPClient, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	return &HTTPClient{
		codecs:     codecs,
		httpClient: httpClient,
		scheme:     u.Scheme,
		host:       u.Host,
		pathPrefix: strings.TrimSuffix(u.Path, "/"),
	}, nil
}

{{- range .DocMethods}}

{{$op := getOperation .Name}}
{{$pathParams := pathParams $op.Request.Params}}
{{$queryParams := queryParams $op.Request.Params}}
{{$headerParams := headerParams $op.Request.Params}}
{{$nonCtxParams := nonCtxParams $op.Request.Params}}
{{$bodyParams := bodyParams $nonCtxParams}}
{{$nonErrReturns := nonErrReturns .Returns}}

func (c *HTTPClient) {{.Name}}({{joinParams .Params "$Name $Type" ", "}}) ({{joinParams .Returns "$Name $Type" ", "}}) {
	codec := c.codecs.EncodeDecoder("{{.Name}}")

	{{if $pathParams -}}
	{{- $fmtPatternParams := patternToFmt $op.Pattern $pathParams}}
	path := fmt.Sprintf("{{$fmtPatternParams.Pattern}}",
		{{- range $fmtPatternParams.SortedParams}}
		{{.}},
		{{- end}}
	)
	{{- else -}}
	path := "{{$op.Pattern}}"
	{{- end }}
	u := &url.URL{
		Scheme: c.scheme,
		Host:   c.host,
		Path:   c.pathPrefix+path,
	}

	{{if $queryParams -}}
	q := u.Query()
	{{- range $queryParams}}
	q.Set("{{.Alias}}", codec.EncodeRequestParam("{{.Name}}", {{.Name}}))
	{{- end}}
	u.RawQuery = q.Encode()
	{{- end}}

	{{if $bodyParams -}}
	reqBody := struct {
		{{- range $bodyParams}}
		{{title .Name}} {{.Type}} {{addTag .Name .Type}}
		{{- end}}
	}{
		{{- range $bodyParams}}
		{{title .Name}}: {{.Name}},
		{{- end}}
	}
	reqBodyReader, headers, err := codec.EncodeRequestBody(&reqBody)
	if err != nil {
		return {{returnErr .Returns}}
	}

	req, err := http.NewRequest("{{$op.Method}}", u.String(), reqBodyReader)
	if err != nil {
		return {{returnErr .Returns}}
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}
	{{- range $headerParams}}
	req.Header.Set("{{.Alias}}", codec.EncodeRequestParam("{{.Name}}", {{.Name}}))
	{{end}}

	{{- else -}}

	req, err := http.NewRequest("{{$op.Method}}", u.String(), nil)
	if err != nil {
		return {{returnErr .Returns}}
	}
	{{- range $headerParams}}
	req.Header.Set("{{.Alias}}", {{.Name}})
	{{end}}
	{{- end}}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return {{returnErr .Returns}}
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusOK && resp.StatusCode <= http.StatusNoContent {
		{{- if $nonErrReturns}}
			var respBody struct {
				{{- range $nonErrReturns}}
				{{title .Name}} {{.Type}} {{addTag .Name .Type}}
				{{- end}}
			}
			err := codec.DecodeSuccessResponse(resp.Body, &respBody)
			if err != nil {
				return {{returnErr .Returns}}
			}
			return {{joinParams $nonErrReturns "respBody.>Name" ", "}}, nil
		{{- else}}
			return nil
		{{- end}}
	} else {
		var respErr error
		err := codec.DecodeFailureResponse(resp.Body, &respErr)
		if err == nil {
			err = respErr
		}
		return {{returnErr .Returns}}
	}
}
{{- end}}
`
)

type RequestField struct {
	Name  string
	Value string
}

type Server struct {
	Service     interface{}
	NewEndpoint interface{}
	Request     interface{}
	Response    interface{}
}

type Options struct {
	SchemaPtr         bool
	SchemaTag         string
	TagKeyToSnakeCase bool
	Formatted         bool
}

type Generator struct {
	opts *Options
}

func New(opts *Options) *Generator {
	return &Generator{opts: opts}
}

func (g *Generator) Generate(result *reflector.Result, spec *openapi.Specification) ([]byte, error) {
	operationMap := make(map[string]*openapi.Operation)
	for _, op := range spec.Operations {
		operationMap[op.Name] = op
	}

	var docMethods []*reflector.Method
	for _, m := range result.Interface.Methods {
		if _, ok := operationMap[m.Name]; ok {
			docMethods = append(docMethods, m)
		}
	}

	data := struct {
		Result     *reflector.Result
		DocMethods []*reflector.Method
		Opts       *Options
	}{
		Result:     result,
		DocMethods: docMethods,
		Opts:       g.opts,
	}

	type FmtPatternParams struct {
		Pattern      string
		SortedParams []string
	}

	return generator.Generate(template, data, generator.Options{
		Funcs: map[string]interface{}{
			"joinParams": func(params []*reflector.Param, format, sep string) string {
				var results []string

				for _, p := range params {
					r := strings.NewReplacer(
						"$Name", p.Name, "$Type", p.Type,
						">Name", strings.Title(p.Name),
					)
					results = append(results, r.Replace(format))
				}
				return strings.Join(results, sep)
			},
			"nonErrReturns": func(params []*reflector.Param) (out []*reflector.Param) {
				for _, p := range params {
					if p.Type != "error" {
						out = append(out, p)
					}
				}
				return
			},
			"getOperation": func(name string) *openapi.Operation {
				return operationMap[name]
			},
			"title": strings.Title,
			"addTag": func(name, typ string) string {
				if g.opts.SchemaTag == "" {
					return ""
				}

				if typ == "error" {
					name = "-"
				} else if g.opts.TagKeyToSnakeCase {
					name = misc.ToSnakeCase(name)
				}

				return fmt.Sprintf("`%s:\"%s\"`", g.opts.SchemaTag, name)
			},
			"addAmpersand": func(name string) string {
				if g.opts.SchemaPtr {
					return "&" + name
				}
				return name
			},
			"patternToFmt": func(pattern string, params []*openapi.Param) FmtPatternParams {
				type nameType struct {
					Name string
					Type string
				}
				type nameIndex struct {
					NameType nameType
					Index    int
				}
				var nameIndices []nameIndex
				for _, p := range params {
					varname := "{" + p.Alias + "}"
					i := strings.Index(pattern, varname)
					if i == -1 {
						panic(fmt.Errorf("no param named %s in pattern %s", p.Alias, pattern))
					}
					nameIndices = append(nameIndices, nameIndex{Index: i, NameType: nameType{Name: p.Name, Type: p.Type}})

					pattern = strings.ReplaceAll(pattern, varname, "%s")
				}

				// Sort names by index.
				sort.Slice(nameIndices, func(i, j int) bool {
					return nameIndices[i].Index < nameIndices[j].Index
				})

				var sortedNames []string
				for _, ni := range nameIndices {
					name := ni.NameType.Name
					sortedNames = append(sortedNames, fmt.Sprintf("codec.EncodeRequestParam(%q, %s)", name, name))
				}
				return FmtPatternParams{
					Pattern:      pattern,
					SortedParams: sortedNames,
				}
			},
			"bodyParams": func(in []*openapi.Param) (out []*openapi.Param) {
				for _, p := range in {
					if p.In == openapi.InBody {
						out = append(out, p)
					}
				}
				return
			},
			"pathParams": func(in []*openapi.Param) (out []*openapi.Param) {
				for _, p := range in {
					if p.In == openapi.InPath {
						out = append(out, p)
					}
				}
				return
			},
			"queryParams": func(in []*openapi.Param) (out []*openapi.Param) {
				for _, p := range in {
					if p.In == openapi.InQuery {
						out = append(out, p)
					}
				}
				return
			},
			"headerParams": func(in []*openapi.Param) (out []*openapi.Param) {
				for _, p := range in {
					if p.In == openapi.InHeader {
						out = append(out, p)
					}
				}
				return
			},
			"nonCtxParams": func(params []*openapi.Param) (out []*openapi.Param) {
				for _, p := range params {
					if p.Type != "context.Context" {
						out = append(out, p)
					}
				}
				return
			},
			"returnErr": func(params []*reflector.Param) string {
				emptyValue := func(typ string) string {
					switch typ {
					case "int", "int8", "int16", "int32", "int64",
						"uint", "uint8", "uint16", "uint32", "uint64":
						return "0"
					case "string":
						return `""`
					case "bool":
						return "false"
					default:
						if strings.HasPrefix(typ, "map") || //map
							strings.HasPrefix(typ, "chan") || // channel
							strings.HasPrefix(typ, "[") || // slice or array
							strings.HasPrefix(typ, "*") { // pointer
							return "nil"
						} else {
							// interface or struct
							return typ + "{}"
						}
					}
				}

				var returns []string
				for i := 0; i < len(params)-1; i++ {
					returns = append(returns, emptyValue(params[i].Type))
				}

				returns = append(returns, "err")

				return strings.Join(returns, ", ")
			},
		},
		Formatted: g.opts.Formatted,
	})
}
