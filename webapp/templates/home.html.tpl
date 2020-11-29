<!DOCTYPE html>
<html lang="fr">
	{{ template "header.html" . }}
	{{ template "logout.html" . }}
	<div class="text-center col-sm-12">
		<div><h2>Hi {{ .User.FirstName }}</h2></div>
		<div class="well well-lg center block">
			<a class="btn btn-lg btn-success btn-block" href="/operations?category=1" role="button">
				<span class="glyphicon glyphicon-hourglass"></span>&nbsp;<span>Calculs Mentaux</span>
			</a>
			<a class="btn btn-lg btn-primary btn-block" href="/operations?category=2" role="button">
				<span class="glyphicon glyphicon-pencil"></span>&nbsp;<span>Op&eacute;rations Pos&eacute;es</span>
			</a>
			<a class="btn btn-lg btn-default btn-block" href="/results" role="button">
				<span class="glyphicon glyphicon-education"></span>&nbsp;<span>R&eacute;sultats</span>
			</a>
		</div>
	</div>
	{{ template "footer.html" . }}
</html>