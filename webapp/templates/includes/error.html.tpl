{{ define "error.html" }}
    {{ if .ErrorMessage }}
    <div class="col-sm-12 text-center"><span class="glyphicon glyphicon-warning-sign"></span>&nbsp;<span>{{ .ErrorMessage }}</span></div>
    {{ end }}
{{ end }}