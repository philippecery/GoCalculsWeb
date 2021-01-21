{{ define "header.html" }}
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>{{ .i18n_title }}</title>
    <link rel="icon" type="image/png" href="/img/maths.png"/>
    <link rel="stylesheet" href="/css/bootstrap.min.css">
    <link rel="stylesheet" href="/css/range.css">
    <style nonce="{{ .nonce }}">
      div.percentAll {
        height: 40px;
      }
      div.percentGood {
        width: 0%;
      }
    </style>
  </head>
  <body>
    <div class="container">
      <div class="text-left col-sm-6">
        {{ range $key, $value := .langs }}
        <a id="lang_{{ $key }}" role="button"><img src="/img/{{ $key }}.svg" alt="{{ $value }}"></a>
        {{ end }}
      </div>
      <div class="text-right col-sm-6">
        {{ if .User }}
        <a class="btn btn-xs btn-primary" id="profile" href="/user/profile" role="button">
          <span class="glyphicon glyphicon-user"></span>&nbsp;<span>{{ .i18n_viewProfile }}</span>
        </a>
        <a class="btn btn-xs btn-danger" href="/logout" role="button">
          <span class="glyphicon glyphicon-log-out"></span>&nbsp;<span>{{ .i18n_logout }}</span>
        </a>
        {{ end }}
      </div>
      <h1 class="text-center alert alert-danger"><span class="glyphicon glyphicon-education"></span>&nbsp;{{ .i18n_title }}</h1>
{{ end }}