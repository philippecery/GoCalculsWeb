<!DOCTYPE html>
<html lang="{{ .lang }}">
	{{ template "header.html" . }}
	{{ template "teacher.navbar.html" . }}
    <div class="col-sm-12 text-center"><h2><span class="glyphicon glyphicon-user"></span>&nbsp;<span>{{ .i18n_grades }}</span></h2></div>
	<div class="col-sm-12">
		<table id="students" class="table table-striped">
			<thead>
				<tr>
					<th rowspan="2">{{ .i18n_gradeName }}</th>
					<th rowspan="2">{{ .i18n_gradeDescription }}</th>
					<th colspan="5">{{ .i18n_mentalmath }}</th>
					<th colspan="5">{{ .i18n_columnform }}</th>
					<th rowspan="2"></th>
					<th rowspan="2"></th>
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
				{{ $i18n_editGrade := .i18n_editGrade }}
				{{ $i18n_copyGrade := .i18n_copyGrade }}
				{{ $i18n_deleteGrade := .i18n_deleteGrade }}
				{{ range .Grades }}
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
					<td class="text-center"><a href="/teacher/grade/edit?gradeid={{ .GradeID }}" data-toggle="tooltip" data-placement="top" title="{{ $i18n_editGrade }}"><span class="glyphicon glyphicon-pencil"></span></a></td>
					<td class="text-center"><a href="/teacher/grade/copy?gradeid={{ .GradeID }}" data-toggle="tooltip" data-placement="top" title="{{ $i18n_copyGrade }}"><span class="glyphicon glyphicon-copy"></span></a></td>
					<td class="text-center"><a href="/teacher/grade/delete?gradeid={{ .GradeID }}&rnd={{ .ActionToken }}" data-toggle="tooltip" data-placement="top" title="{{ $i18n_deleteGrade }}"><span class="glyphicon glyphicon-trash"></span></a></td>
				</tr>
				{{ end }}
			</tbody>
		</table>
	</div>
	<div class="text-center col-sm-12">
		<a class="btn btn-lg btn-success btn-block" href="/teacher/grade/new" role="button">
			<span class="glyphicon glyphicon-plus"></span>&nbsp;<span>{{ .i18n_addGrade }}</span>
		</a>
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