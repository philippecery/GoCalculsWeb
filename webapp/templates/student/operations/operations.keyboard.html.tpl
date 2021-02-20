{{ define "operations.keyboard.html" }}
{{ $keys := .Keys }}
{{ range .IDs }}
{{ $id := . }}
{{ $class := "primary" }}
{{ if eq $id "keyboard2" }}
    {{ $class = "success" }}
{{ end }}
<div id="{{ $id }}" class="col-sm-12 btn-toolbar hidden" role="toolbar">
	<div class="col-sm-11">
		<div class="btn-group btn-group-justified" role="group">
            {{ range $key := $keys }}
			<div class="btn-group" role="group">
				<button type="button" id="keynum_{{ $key }}" class="btn btn-{{ $class }} btn-lg">{{ if eq $key "dot" }}{{ "." }}{{ else }}{{ $key }}{{ end }}</button>
			</div>
            {{ end }}
		</div>
	</div>
	<div class="col-sm-1">
		<div class="btn-group">
			<button type="button" id="keydel" class="btn btn-danger btn-lg"><span class="glyphicon glyphicon-arrow-left"></span></button>
		</div>
	</div>
</div>
{{ end }}
{{ end }}