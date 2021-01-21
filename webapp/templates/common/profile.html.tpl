<!DOCTYPE html>
<html lang="{{ .lang }}">
	{{ template "header.html" . }}
    <div class="col-sm-12 text-center"><h2><span class="glyphicon glyphicon-user"></span>&nbsp;<span>{{ .i18n_profile }}</span></h2></div>
	{{ template "error.html" . }}
    <div class="col-sm-12 text-center"><span class="glyphicon glyphicon-info-sign"></span>&nbsp;<span>{{ .i18n_lastConnection }}: {{ .LastConnection }}</span></div>
    <div class="col-sm-12 text-center">&nbsp;</div>
    <div class="col-sm-12">
        <form id="profileForm" method="POST" action="/profile" class="form-horizontal">
            <input type="hidden" name="token" value="{{ .Token }}"/>
            <div class="form-group">
                <label for="userId" class="col-sm-2 control-label">{{ .i18n_userid }}</label>
                <div class="col-sm-3">
                    <input type="text" autocomplete="username" class="form-control" id="userId" name="userId" readonly value="{{ .UserProfile.UserID }}"/>
                </div>
            </div>
            <div class="form-group">
                <label for="firstName" class="col-sm-2 control-label">{{ .i18n_firstName }}</label>
                <div class="col-sm-3">
                    <input type="text" autocomplete="given-name" class="form-control" id="firstName" name="firstName" value="{{ .UserProfile.FirstName }}"/>
                </div>
            </div>
            <div class="form-group">
                <label for="lastName" class="col-sm-2 control-label">{{ .i18n_lastName }}</label>
                <div class="col-sm-3">
                    <input type="text" autocomplete="family-name" class="form-control" id="lastName" name="lastName" value="{{ .UserProfile.LastName }}"/>
                </div>
            </div>
            <div class="form-group">
                <label for="emailAddress" class="col-sm-2 control-label">{{ .i18n_emailAddress }}</label>
                <div class="col-sm-3">
                    <input type="email" autocomplete="email" pattern="^(?=^.{6,254}$)[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$" class="form-control" id="emailAddress" name="emailAddress" value="{{ .UserProfile.EmailAddress }}"/>
                </div>
            </div>
            <div class="form-group">
                <div class="col-sm-offset-2 col-sm-10">
                    <div class="btn-group" role="group">
                        <button id="saveProfile" type="submit" class="btn btn-lg btn-primary">{{ .i18n_save }}</button>
                        <a class="btn btn-lg btn-default" href="{{ .LastVisitedPage }}" role="button"><span>{{ .i18n_cancel }}</span></a>
                    </div>
                </div>
            </div>
        </form>
    </div>
	{{ template "footer.html" . }}
</html>