{{ define "teacher.navbar.html" }}
<ul class="nav nav-tabs">
  <li role="presentation"><a href="/teacher/student/list">{{ .i18n_students }}</a></li>
  <li role="presentation"><a href="/teacher/grade/list">{{ .i18n_grades }}</a></li>
</ul>
{{ end }}