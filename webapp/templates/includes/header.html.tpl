{{ define "header.html" }}
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>{{ .i18n_title }}</title>
    <link href="/css/bootstrap.min.css" rel="stylesheet">
    <link href="/css/range.css" rel="stylesheet">
  </head>
  <body>
    <div class="container">
      <div class="text-right col-sm-12">
        {{ if .User }}
        <a class="btn btn-xs btn-danger" href="/logout" role="button">
          <span class="glyphicon glyphicon-log-out"></span>&nbsp;<span>{{ .i18n_logout }}</span>
        </a>
        {{ end }}
        <div class="btn-group" role="group">
            {{ range $key, $value := .langs }}
              <a class="btn btn-xs btn-default" id="lang_{{ $key }}" role="button"><span>{{ $value }}</span></a>
            {{ end }}
        </div>
      </div>
      <h1 class="text-center alert alert-danger"><span class="glyphicon glyphicon-education"></span>&nbsp;{{ .i18n_title }}</h1>
{{ end }}