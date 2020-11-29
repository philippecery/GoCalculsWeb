<!DOCTYPE html>
<html lang="fr">
	{{ template "header.html" . }}
    <div class="col-sm-12 text-center"><h2><span class="glyphicon glyphicon-user"></span>&nbsp;<span>Registration</span></h2></div>
	{{ template "error.html" . }}
    <div class="col-sm-12">
        <form id="registrationForm" method="POST" action="/register" class="form-horizontal">
            <input type="hidden" name="token" value="{{ .Token }}"/>
            <div class="form-group">
                <label for="userId" class="col-sm-2 control-label">User ID</label>
                <div class="col-sm-3">
                    <input type="text" class="form-control" id="userId" name="userId" readonly value="{{ .UserID }}"/>
                </div>
            </div>
            <div class="form-group">
                <label for="emailAddress" class="col-sm-2 control-label">Email Address</label>
                <div class="col-sm-3">
                    <input type="text" class="form-control" id="emailAddress" name="emailAddress"/>
                </div>
            </div>
            <div class="form-group">
                <label for="password" class="col-sm-2 control-label">Password</label>
                <div class="col-sm-3">
                    <input type="password" class="form-control" id="password" name="password"/>
                </div>
            </div>
            <div class="form-group">
                <label for="passwordConfirm" class="col-sm-2 control-label">Password</label>
                <div class="col-sm-3">
                    <input type="password" class="form-control" id="passwordConfirm" name="passwordConfirm"/>
                </div>
            </div>
            <div class="form-group">
                <label for="firstName" class="col-sm-2 control-label">First Name</label>
                <div class="col-sm-3">
                    <input type="text" class="form-control" id="firstName" name="firstName"/>
                </div>
            </div>
            <div class="form-group">
                <label for="lastName" class="col-sm-2 control-label">Last Name</label>
                <div class="col-sm-3">
                    <input type="text" class="form-control" id="lastName" name="lastName"/>
                </div>
            </div>
            <div class="form-group">
                <div class="col-sm-offset-2 col-sm-10">
                    <button id="register" type="submit" class="btn btn-lg btn-default">Register</button>
                </div>
            </div>
        </form>
    </div>
	{{ template "footer.html" . }}
</html>