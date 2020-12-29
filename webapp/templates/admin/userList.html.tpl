<!DOCTYPE html>
<html lang="{{ .lang }}">
	{{ template "header.html" . }}
    <div class="col-sm-12 text-center"><h2><span class="glyphicon glyphicon-user"></span>&nbsp;<span>{{ .i18n_registeredUsers }}</span></h2></div>
	<div class="col-sm-12">
		<table id="users" class="table table-striped">
			<thead>
				<tr>
					<th>{{ .i18n_userid }}</th>
					<th>{{ .i18n_firstName }}</th>
					<th>{{ .i18n_lastName }}</th> 
					<th>{{ .i18n_emailAddress }}</th>
					<th>{{ .i18n_role }}</th>
					<th>{{ .i18n_lastConnection }}</th>
					<th></th>
					<th></th>
				</tr>
			</thead>
			<tbody id="usersData">
				{{ $currentUserID := .User.UserID }}
				{{ $i18n_disableAccount := .i18n_disableAccount }}
				{{ $i18n_enableAccount := .i18n_enableAccount }}
				{{ $i18n_delete := .i18n_deleteUser }}
				{{ range .RegisteredUsers }}
				<tr>
					<td>{{ .UserID }}</td>
					<td>{{ .FirstName }}</td>
					<td>{{ .LastName }}</td>
					<td>{{ .EmailAddress }}</td>
					<td>{{ .Role }}</td>
					<td>{{ .LastConnection }}</td>
					<td class="text-center">{{ if eq $currentUserID .UserID }}<span class="glyphicon glyphicon-ban-circle"></span>{{ else }}<a href="/admin/user/status?userid={{ .UserID }}&rnd={{ .ActionToken }}" data-toggle="tooltip" data-placement="top" title="{{ if .Enabled }}{{ $i18n_disableAccount }}{{ else }}{{ $i18n_enableAccount }}{{ end }}"><span class="glyphicon glyphicon-{{ if .Enabled }}ok{{ else }}remove{{ end }}-circle"></span></a>{{ end }}</td>
					<!-- TODO: Add confirmation before deletion-->
					<td class="text-center">{{ if eq $currentUserID .UserID }}<span class="glyphicon glyphicon-ban-circle"></span>{{ else }}<a href="/admin/user/delete?userid={{ .UserID }}&rnd={{ .ActionToken }}" data-toggle="tooltip" data-placement="top" title="{{ $i18n_delete }}"><span class="glyphicon glyphicon-trash"></span></a>{{ end }}</td>
				</tr>
				{{ end }}
			</tbody>
		</table>
	</div>
    <div class="col-sm-12 text-center"><h2><span class="glyphicon glyphicon-user"></span>&nbsp;<span>{{ .i18n_unregisteredUsers }}</span></h2></div>
	<div class="col-sm-12">
		<table id="userTokens" class="table table-striped">
			<thead>
				<tr>
					<th>{{ .i18n_userid }}</th>
					<th>{{ .i18n_role }}</th>
					<th>{{ .i18n_token }}</th>
					<th>{{ .i18n_expires }}</th>
					<th></th>
					<th></th>
				</tr>
			</thead>
			<tbody id="userTokensData">
				{{ $i18n_copyRegistrationLink := .i18n_copyRegistrationLink }}
				{{ range .UnregisteredUsers }}
				<tr>
					<td>{{ .UserID }}</td>
					<td>{{ .Role }}</td>
					<td>{{ .Token }}</td>
					<td>{{ .Expires }}</td>
					<td class="text-center"><a href="{{ .Link }}" id="link_{{ .UserID }}" data-toggle="tooltip" data-placement="top" title="{{ $i18n_copyRegistrationLink }}"><span class="glyphicon glyphicon-copy"></span></a></td>
					<!-- TODO: Add confirmation before deletion-->
					<td class="text-center"><a href="/admin/user/delete?userid={{ .UserID }}&rnd={{ .ActionToken }}" data-toggle="tooltip" data-placement="top" title="{{ $i18n_delete }}"><span class="glyphicon glyphicon-trash"></span></a></td>
				</tr>
				{{ end }}
			</tbody>
		</table>
	</div>
	<div class="text-center col-sm-12">
		<a class="btn btn-lg btn-success btn-block" href="/admin/user/new" role="button">
			<span class="glyphicon glyphicon-plus"></span>&nbsp;<span>{{ .i18n_addUser }}</span>
		</a>
	</div>
	{{ template "footer.html" . }}
	<script nonce="{{ .nonce }}" type="text/javascript" src="/js/jquery-2.2.2.min.js"></script>
	<script nonce="{{ .nonce }}" type="text/javascript" src="/js/bootstrap.min.js"></script>
	<script nonce="{{ .nonce }}">
	$(document).ready(function(){
		$('[data-toggle="tooltip"]').tooltip();
		{{ range .UnregisteredUsers }}
		document.getElementById('link_{{ .UserID }}').addEventListener('click', function(evt) { copyURI(evt) });
		{{ end }}
	});
	function copyURI(evt) {
		evt.preventDefault();
		navigator.clipboard.writeText(window.location.origin + evt.currentTarget.getAttribute('href')).then(() => {
		}, () => {
		});
	}
	</script>
</html>