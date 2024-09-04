package main

import (
	"os"
	"text/template"
)

// template今もよく使うやつ
// else句とwith句を合わせたelse with句が追加された
// こちらも参考: https://future-architect.github.io/articles/20240722a/
type Inventory struct {
	Material string
	Count    uint
}

type Profile struct {
	Name string
	Age  int
}
type Profiles []Profile

func pipeline(name string) string {
	m := map[string]string{
		"Alice": "Alice",
		"Bob":   "Bob",
		"Carol": "Carol",
	}
	return m[name]
}

// ifに近いね。ifのelse句がない場合に使う
const TemplData = `
	{{with pipeline .Name}}
		{{ . }}これがわかりやすいかも
	{{else}}
		noting
	{{end}}
	{{with pipeline .Name}} T1 {{else}}{{with pipeline .Name}} T0 {{end}}{{end}}
	じゃなくて
	{{with pipeline .Name}} T1 {{else with pipeline .Name}} T0 {{end}}
	が使えるようになった

	年齢は{{.Age}}です
`

func main() {
	// sweaters := Inventory{"wool", 17}
	sweaters := Profiles{
		Profile{"Alice", 20},
		Profile{"aaa", 21},
	}
	// tmpl, err := template.New("test").Parse(TemplData)
	tmpl, err := template.New("test").Funcs(template.FuncMap{"pipeline": pipeline}).Parse(TemplData)
	if err != nil {
		panic(err)
	}
	for _, p := range sweaters {
		err = tmpl.Execute(os.Stdout, p)
	}
	if err != nil {
		panic(err)
	}

}
