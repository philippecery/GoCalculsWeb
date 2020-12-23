<!DOCTYPE html>
<html lang="{{ .lang }}">
	{{ template "header.html" . }}
    <div class="col-sm-12 text-center"><h2><span class="glyphicon glyphicon-user"></span>&nbsp;<span>{{ .i18n_newUser }}</span></h2></div>
	{{ template "error.html" . }}
    <div class="col-sm-12">
        <form id="registrationForm" method="POST" action="/admin/newUser" class="form-horizontal">
            <input type="hidden" name="token" value="{{ .Token }}"/>
            <div class="form-group">
                <label for="userId" class="col-sm-2 control-label">{{ .i18n_userid }}</label>
                <div class="col-sm-3">
                    <input type="text" autocomplete="off" pattern="^(?=^[a-z.]{3,32}$)[a-z]{2,}(\.?[a-z]{2,})*$" title="Between 3 and 32 characters. Only lower case letters and dot (.) character are allowed. Must start and end with letters." class="form-control" id="userId" name="userId"/>
                </div>
            </div>
            <div class="form-group">
                <label for="role" class="col-sm-2 control-label">{{ .i18n_role }}</label>
                <div class="col-sm-3">
                    <select class="form-control" id="role" name="role">
                        <option value="0">{{ .i18n_select }}</option>
                        <option value="1">{{ .i18n_admin }}</option>
                        <option value="2">{{ .i18n_teacher }}</option>
                        <option value="3">{{ .i18n_student }}</option>
                    </select>
                </div>
            </div>
            <div class="form-group">
                <div class="col-sm-offset-2 col-sm-10">
                    <div class="btn-group" role="group">
                        <button id="register" type="submit" class="btn btn-lg btn-success">{{ .i18n_createUser }}</button>
                        <a class="btn btn-lg btn-default" href="/admin/users" role="button"><span>{{ .i18n_cancel }}</span></a>
                    </div>
                </div>
            </div>
        </form>
    </div>
	{{ template "footer.html" . }}
</html>