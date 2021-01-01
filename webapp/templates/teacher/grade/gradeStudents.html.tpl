<!DOCTYPE html>
<html lang="{{ .lang }}">
	{{ template "header.html" . }}
	{{ template "teacher.navbar.html" . }}
    <div class="col-sm-12 text-center"><h2><span class="glyphicon glyphicon-tasks"></span>&nbsp;<span>{{ .Grade.Name }}</span></h2></div>
    <div class="col-sm-12 text-center"><h3><span class="glyphicon glyphicon-user"></span>&nbsp;<span>{{ .i18n_students }}</span></h3></div>
	<div class="col-sm-12">
		<table id="students" class="table table-striped">
			<thead>
				<tr>
					<th>{{ .i18n_firstName }}</th>
					<th>{{ .i18n_lastName }}</th>
					<th></th>
				</tr>
			</thead>
			<tbody id="studentsData">
				{{ $i18n_unassignGrade := .i18n_unassignGrade }}
				{{ $gradeID := .Grade.GradeID }}
				{{ range .Students }}
				{{ if eq .GradeID $gradeID }}
				<tr>
					<td>{{ .FirstName }}</td>
					<td>{{ .LastName }}</td>
					<td class="text-center"><a href="/teacher/grade/unassign?gradeid={{ $gradeID }}&userid={{ .UserID }}&rnd={{ .ActionToken }}" data-toggle="tooltip" data-placement="top" title="{{ $i18n_unassignGrade }}"><span class="glyphicon glyphicon-remove"></span></a></td>
				</tr>
				{{ end }}
				{{ end }}
			</tbody>
		</table>
	</div>
    <div class="col-sm-12 text-center"><h3><span class="glyphicon glyphicon-user"></span>&nbsp;<span>{{ .i18n_otherStudents }}</span></h3></div>
	<div class="col-sm-12">
        <form id="gradeStudentsForm" method="POST" action="/teacher/grade/students" class="form-horizontal">
            <input type="hidden" name="token" value="{{ .Token }}"/>
            <input type="hidden" name="gradeID" value="{{ $gradeID }}"/>
			<table id="students" class="table table-striped">
				<thead>
					<tr>
						<th rowspan="2">{{ .i18n_firstName }}</th>
						<th rowspan="2">{{ .i18n_lastName }}</th>
						<th rowspan="2">{{ .i18n_gradeName }}</th>
						<th colspan="5">{{ .i18n_mentalmath }}</th>
						<th colspan="5">{{ .i18n_columnform }}</th>
						<th rowspan="2"></th>
					</tr>
					<tr>
						<th>+</th>
						<th>-</th>
						<th>*</th>
						<th>/</th>
						<th><span class="glyphicon glyphicon-time"></span></th>
						<th>+</th>
						<th>-</th>
						<th>*</th>
						<th>/</th>
						<th><span class="glyphicon glyphicon-time"></span></th>
					</tr>
				</thead>
				<tbody id="studentsData">
					{{ $i18n_changeGrade := .i18n_changeGrade }}
					{{ $i18n_nograde := .i18n_nograde }}
					{{ range .Students }}
					{{ if ne .GradeID $gradeID }}
					<tr>
						<td>{{ .FirstName }}</td>
						<td>{{ .LastName }}</td>
						{{ if .Grade }}
						<td>{{ .Grade.Name }}</td>
						<td>{{ .Grade.MentalMath.NbAdditions }}</td>
						<td>{{ .Grade.MentalMath.NbSubstractions }}</td>
						<td>{{ .Grade.MentalMath.NbMultiplications }}</td>
						<td>{{ .Grade.MentalMath.NbDivisions }}</td>
						<td>{{ .Grade.MentalMath.Time }}</td>
						<td>{{ .Grade.ColumnForm.NbAdditions }}</td>
						<td>{{ .Grade.ColumnForm.NbSubstractions }}</td>
						<td>{{ .Grade.ColumnForm.NbMultiplications }}</td>
						<td>{{ .Grade.ColumnForm.NbDivisions }}</td>
						<td>{{ .Grade.ColumnForm.Time }}</td>
						{{ else }}
						<td colspan="11">{{ $i18n_nograde }}</td>
						{{ end }}
						<td class="text-center"><input type="checkbox" name="selectedStudents" value="{{ .UserID }}"/></td>
					</tr>
					{{ end }}
					{{ end }}
				</tbody>
			</table>
            <div class="form-group">
                <div class="col-sm-offset-2 col-sm-10">
                    <div class="btn-group" role="group">
                        <button id="save" type="submit" class="btn btn-lg btn-primary">{{ .i18n_save }}</button>
                        <a class="btn btn-lg btn-default" href="/teacher/grade/list" role="button"><span>{{ .i18n_cancel }}</span></a>
                    </div>
                </div>
            </div>
		</form>
	</div>
	{{ template "footer.html" . }}
	<script nonce="{{ .nonce }}" type="text/javascript" src="/js/jquery-2.2.2.min.js"></script>
	<script nonce="{{ .nonce }}" type="text/javascript" src="/js/bootstrap.min.js"></script>
	<script nonce="{{ .nonce }}">
	$(document).ready(function(){
		$('[data-toggle="tooltip"]').tooltip();
	});
	</script>
</html>