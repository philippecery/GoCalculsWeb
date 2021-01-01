<!DOCTYPE html>
<html lang="{{ .lang }}">
	{{ template "header.html" . }}
	{{ template "teacher.navbar.html" . }}
    <div class="col-sm-12 text-center"><h2><span class="glyphicon glyphicon-user"></span>&nbsp;<span>{{ .i18n_students }}</span></h2></div>
	<div class="col-sm-12">
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
					<td class="text-center"><a href="/teacher/student/edit?userid={{ .UserID }}" data-toggle="tooltip" data-placement="top" title="{{ $i18n_changeGrade }}"><span class="glyphicon glyphicon-pencil"></span></a></td>
				</tr>
				{{ end }}
			</tbody>
		</table>
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