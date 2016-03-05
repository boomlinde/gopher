package main

var tpltext = `<!doctype html>
<html>
<head>
<meta charset="utf-8">
<title>{{.Title}}</title>
</head>
<body>
<pre>
{{range .Lines}} {{if .Link}}({{.Type}}) <a href="{{.Link}}">{{.Text}}</a>{{else}}      {{.Text}}{{end}}
{{end}}</pre>
</body>
</html>`
