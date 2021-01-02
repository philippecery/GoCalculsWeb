<!DOCTYPE html>
<html lang="{{ .lang }}">
	{{ template "header.html" . }}
	{{ template "teacher.navbar.html" . }}
    {{ $currentGrade := .Student.Grade }}
    <div class="col-sm-12 text-center"><h2><span class="glyphicon glyphicon-education"></span>&nbsp;<span>{{ .Student.FirstName }} {{ .Student.LastName }}</span></h2></div>
    <div class="col-sm-12 text-center"><h3><span>{{ .i18n_currentGrade }}: {{ $currentGrade.Name }}</span></h3></div>
	<div class="col-sm-12">
		<table id="grades" class="table table-striped">
			<thead>
				<tr>
					<th rowspan="2">{{ .i18n_gradeName }}</th>
					<th rowspan="2">{{ .i18n_gradeDescription }}</th>
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
			<tbody id="studentGradesData">
				{{ $i18n_assignGrade := .i18n_assignGrade }}
				{{ $userID := .Student.UserID }}
				{{ range .Grades }}
                {{ if ne .GradeID $currentGrade.GradeID }}
				<tr>
					<td>{{ .Name }}</td>
					<td>{{ .Description }}</td>
					<td>{{ .MentalMath.NbAdditions }}</td>
					<td>{{ .MentalMath.NbSubstractions }}</td>
					<td>{{ .MentalMath.NbMultiplications }}</td>
					<td>{{ .MentalMath.NbDivisions }}</td>
					<td>{{ .MentalMath.Time }}</td>
					<td>{{ .ColumnForm.NbAdditions }}</td>
					<td>{{ .ColumnForm.NbSubstractions }}</td>
					<td>{{ .ColumnForm.NbMultiplications }}</td>
					<td>{{ .ColumnForm.NbDivisions }}</td>
					<td>{{ .ColumnForm.Time }}</td>
					<td class="text-center"><a href="/teacher/student/assign?userid={{ $userID }}&gradeid={{ .GradeID }}&rnd={{ .ActionToken }}" data-toggle="tooltip" data-placement="top" title="{{ $i18n_assignGrade }}"><span class="glyphicon glyphicon-ok-circle"></span></a></td>
				</tr>
                {{ end }}
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