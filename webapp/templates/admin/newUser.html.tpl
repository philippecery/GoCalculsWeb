<!DOCTYPE html>
<html lang="fr">
	{{ template "header.html" . }}
	{{ template "logout.html" . }}
    <div class="col-sm-12 text-center"><h2><span class="glyphicon glyphicon-user"></span>&nbsp;<span>New User</span></h2></div>
	{{ template "error.html" . }}
    <div class="col-sm-12">
        <form id="registrationForm" method="POST" action="/admin/newUser" class="form-horizontal">
            <input type="hidden" name="token" value="{{ .Token }}"/>
            <div class="form-group">
                <label for="userId" class="col-sm-2 control-label">User ID</label>
                <div class="col-sm-3">
                    <input type="text" class="form-control" id="userId" name="userId"/>
                </div>
            </div>
            <div class="form-group">
                <label for="role" class="col-sm-2 control-label">Role</label>
                <div class="col-sm-3">
                    <select class="form-control" id="role" name="role">
                        <option value="0">Select...</option>
                        <option value="1">Admin</option>
                        <option value="2">Teacher</option>
                        <option value="3">Student</option>
                    </select>
                </div>
            </div>
            <div class="form-group">
                <div class="col-sm-offset-2 col-sm-10">
                    <div class="btn-group btn-group-justified" role="group">
                        <div class="btn-group" role="group">
                            <button id="register" type="submit" class="btn btn-lg btn-success btn-block">Register</button>
                        </div>
                        <div class="btn-group" role="group">
                            <a class="btn btn-lg btn-default btn-block" href="/admin/users" role="button"><span>Cancel</span></a>
                        </div>
                    </div>
                </div>
            </div>
        </form>
    </div>
	{{ template "footer.html" . }}
</html>