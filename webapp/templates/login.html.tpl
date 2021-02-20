<!DOCTYPE html>
<html lang="{{ .lang }}">
	{{ template "header.html" . }}
    <div class="col-sm-12 text-center"><h2><span class="glyphicon glyphicon-user"></span>&nbsp;<span>{{ .i18n_login }}</span></h2></div>
	{{ template "error.html" . }}
    <div class="col-sm-12">
        <form id="loginForm" method="POST" action="/login" class="form-horizontal">
            <input type="hidden" name="token" value="{{ .Token }}"/>
            <div class="form-group">
                <label for="userId" class="col-sm-2 control-label">{{ .i18n_userid }}</label>
                <div class="col-sm-3">
                    <input type="text" autocomplete="username" class="form-control" id="userId" name="userId"/>
                </div>
            </div>
            <div class="form-group">
                <label for="password" class="col-sm-2 control-label">{{ .i18n_password }}</label>
                <div class="col-sm-3">
                    <input type="password" autocomplete="current-password" class="form-control" id="password" name="password"/>
                </div>
            </div>
            <div class="form-group">
                <div class="col-sm-offset-2 col-sm-10">
                    <button id="login" type="submit" class="btn btn-lg btn-default">{{ .i18n_login }}</button>
                </div>
            </div>
        </form>
    </div>
	{{ template "footer.html" . }}
</html>