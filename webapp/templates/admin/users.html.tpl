<!DOCTYPE html>
<html lang="fr">
	{{ template "header.html" . }}
	{{ template "logout.html" . }}
    <div class="col-sm-12 text-center"><h2><span class="glyphicon glyphicon-user"></span>&nbsp;<span>Registered Users</span></h2></div>
	<div class="col-sm-12">
		<table id="users" class="table table-striped">
			<thead>
				<tr>
					<th>ID</th>
					<th>First Name</th>
					<th>Last Name</th> 
					<th>Email Address</th>
					<th>Role</th>
					<th>Last Connection</th>
					<th></th>
					<th></th>
				</tr>
			</thead>
			<tbody id="usersData">
				{{ $userID := .UserID }}
				{{ range .RegisteredUsers }}
				<tr>
					<td>{{ .UserID }}</td>
					<td>{{ .FirstName }}</td>
					<td>{{ .LastName }}</td>
					<td>{{ .EmailAddress }}</td>
					<td>{{ .Role }}</td>
					<td>{{ .LastConnection }}</td>
					<td class="text-center">{{ if eq $userID .UserID }}<span class="glyphicon glyphicon-ban-circle"></span>{{ else }}<a href="/admin/status?userid={{ .UserID }}&rnd={{ .ActionToken }}" data-toggle="tooltip" data-placement="top" title="{{ if .Enabled }}Disable{{ else }}Enable{{ end }} Account"><span class="glyphicon glyphicon-{{ if .Enabled }}ok{{ else }}remove{{ end }}-circle"></span></a>{{ end }}</td>
					<!-- TODO: Add confirmation before deletion-->
					<td class="text-center">{{ if eq $userID .UserID }}<span class="glyphicon glyphicon-ban-circle"></span>{{ else }}<a href="/admin/delete?userid={{ .UserID }}&rnd={{ .ActionToken }}" data-toggle="tooltip" data-placement="top" title="Delete User"><span class="glyphicon glyphicon-trash"></span></a>{{ end }}</td>
				</tr>
				{{ end }}
			</tbody>
		</table>
	</div>
    <div class="col-sm-12 text-center"><h2><span class="glyphicon glyphicon-user"></span>&nbsp;<span>Unregistered Users</span></h2></div>
	<div class="col-sm-12">
		<table id="userTokens" class="table table-striped">
			<thead>
				<tr>
					<th>ID</th>
					<th>Role</th>
					<th>Token</th>
					<th>Expires</th>
					<th></th>
					<th></th>
				</tr>
			</thead>
			<tbody id="userTokensData">
				{{ range .UnregisteredUsers }}
				<tr>
					<td>{{ .UserID }}</td>
					<td>{{ .Role }}</td>
					<td>{{ .Token }}</td>
					<td>{{ .Expires }}</td>
					<td class="text-center"><a href="{{ .Link }}" onclick="copyURI(event)" data-toggle="tooltip" data-placement="top" title="Copy Registration Link"><span class="glyphicon glyphicon-copy"></span></a></td>
					<!-- TODO: Add confirmation before deletion-->
					<td class="text-center"><a href="/admin/delete?userid={{ .UserID }}&rnd={{ .ActionToken }}" data-toggle="tooltip" data-placement="top" title="Delete User"><span class="glyphicon glyphicon-trash"></span></a></td>
				</tr>
				{{ end }}
			</tbody>
		</table>
	</div>
	<div class="text-center col-sm-12">
		<a class="btn btn-lg btn-success btn-block" href="/admin/newUser" role="button">
			<span class="glyphicon glyphicon-plus"></span>&nbsp;<span>Add User</span>
		</a>
	</div>
	{{ template "footer.html" . }}
	<script type="text/javascript" src="/js/jquery-2.2.2.min.js"></script>
	<script type="text/javascript" src="/js/bootstrap.min.js"></script>
	<script>
	$(document).ready(function(){
		$('[data-toggle="tooltip"]').tooltip(); 
	});
	function copyURI(evt) {
		evt.preventDefault();
		navigator.clipboard.writeText(window.location.origin + evt.currentTarget.getAttribute('href')).then(() => {
		}, () => {
		});
	}
	</script>
</html>