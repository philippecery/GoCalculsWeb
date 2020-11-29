<!DOCTYPE html>
<html lang="fr">
	{{ template "header.html" . }}
    <div class="col-sm-12 text-center"><h2><span class="glyphicon glyphicon-user"></span>&nbsp;<span>Login</span></h2></div>
	{{ template "error.html" . }}
    <div class="col-sm-12">
        <form id="loginForm" method="POST" action="/login" class="form-horizontal">
            <input type="hidden" name="token" value="{{ .Token }}"/>
            <div class="form-group">
                <label for="userId" class="col-sm-2 control-label">User ID</label>
                <div class="col-sm-3">
                    <input type="text" class="form-control" id="userId" name="userId"/>
                </div>
            </div>
            <div class="form-group">
                <label for="password" class="col-sm-2 control-label">Password</label>
                <div class="col-sm-3">
                    <input type="password" class="form-control" id="password" name="password"/>
                </div>
            </div>
            <div class="form-group">
                <div class="col-sm-offset-2 col-sm-10">
                    <button id="login" type="submit" class="btn btn-lg btn-default">Login</button>
                </div>
            </div>
        </form>
    </div>
	{{ template "footer.html" . }}
</html>