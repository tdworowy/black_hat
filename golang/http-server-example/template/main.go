package main

import (
	"html/template"
	"os"
)

var template_html = `
	<html>
	<body>
		Hello {{.}}
	</body>
	</html>
`

func main() {

	t, err := template.New("Hello").Parse(template_html)
	if err != nil {
		panic(err)
	}
	t.Execute(os.Stdout, "<script>alert('world')</script>")
}
