<!DOCTYPE html>
<html lang="{{ .lang }}">
	{{ template "header.html" . }}
	<div class="text-center col-sm-12">
		<div><h2>{{ .User.FirstName }}</h2></div>
		<div class="well well-lg center block">
			<a class="btn btn-lg btn-success btn-block" href="/student/operations?type=1" role="button">
				<span class="glyphicon glyphicon-hourglass"></span>&nbsp;<span>{{ .i18n_mentalmath }}</span>
			</a>
			<a class="btn btn-lg btn-primary btn-block" href="/student/operations?type=2" role="button">
				<span class="glyphicon glyphicon-pencil"></span>&nbsp;<span>{{ .i18n_columnform }}</span>
			</a>
			<a class="btn btn-lg btn-default btn-block" href="/student/results" role="button">
				<span class="glyphicon glyphicon-education"></span>&nbsp;<span>{{ .i18n_results }}</span>
			</a>
		</div>
	</div>
	{{ template "footer.html" . }}
</html>