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
				{{ $i18n_editStudent := .i18n_editStudent }}
				{{ range .Students }}
				<tr>
					<td>{{ .FirstName }}</td>
					<td>{{ .LastName }}</td>
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
					<td class="text-center"><a href="/teacher/student/edit?userid={{ .UserID }}" data-toggle="tooltip" data-placement="top" title="{{ $i18n_editStudent }}"><span class="glyphicon glyphicon-pencil"></span></a></td>
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