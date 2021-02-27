{{ define "modalSettings.html" }}
<div class="modal fade" id="modalSettings" tabindex="-1" role="dialog" aria-labelledby="Settings">
    <div class="modal-dialog" role="document">
        <div class="modal-content">
            <div class="modal-header">
                <ul class="nav nav-pills">
                    <li role="presentation" class="active"><a href="#" id="showProfile">{{ .i18n_profile }}</a></li>
                    <li role="presentation"><a href="#" id="showChangePassword">{{ .i18n_changePassword }}</a></li>
                </ul>
            </div>
            <div class="modal-body">
                <div id="message" class="col-sm-12 text-center hidden"></div>
                <div id="tabProfile">
                    <div class="row">
                        <div class="col-sm-12 text-center"><span class="glyphicon glyphicon-info-sign"></span>&nbsp;<span>{{ .i18n_lastConnection }}: </span><span id="lastConnection"></span></div>
                    </div>
                    <div class="row">
                        <div class="col-sm-12 text-center">&nbsp;</div>
                    </div>
                    <div class="row">
                        <div class="col-sm-12">
                            <form id="profileForm" method="POST" action="/user/profile" class="form-horizontal">
                                <input type="hidden" name="token" value="{{ .Token }}"/>
                                <div class="form-group">
                                    <label for="userId" class="col-sm-5 control-label">{{ .i18n_userid }}</label>
                                    <div class="col-sm-5">
                                        <input type="text" autocomplete="username" class="form-control" id="userId" name="userId" readonly/>
                                    </div>
                                </div>
                                <div class="form-group">
                                    <label for="firstName" class="col-sm-5 control-label">{{ .i18n_firstName }}</label>
                                    <div class="col-sm-5">
                                        <input type="text" autocomplete="given-name" class="form-control" id="firstName" name="firstName"/>
                                    </div>
                                </div>
                                <div class="form-group">
                                    <label for="lastName" class="col-sm-5 control-label">{{ .i18n_lastName }}</label>
                                    <div class="col-sm-5">
                                        <input type="text" autocomplete="family-name" class="form-control" id="lastName" name="lastName"/>
                                    </div>
                                </div>
                                <div class="form-group">
                                    <label for="emailAddress" class="col-sm-5 control-label">{{ .i18n_emailAddress }}</label>
                                    <div class="col-sm-5">
                                        <input type="email" autocomplete="email" pattern="^(?=^.{6,254}$)[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$" class="form-control" id="emailAddress" name="emailAddress"/>
                                    </div>
                                </div>
                                <div class="form-group">
                                    <div class="col-sm-offset-5 col-sm-5">
                                        <button id="saveProfile" type="submit" class="btn btn-lg btn-primary">{{ .i18n_save }}</button>
                                    </div>
                                </div>
                            </form>
                        </div>
                    </div>
                </div>
                <div id="tabChangePassword" class="hidden">
                    <div class="row">
                        <div class="col-sm-12">
                            <form id="changePasswordForm" method="POST" action="/user/changePassword" class="form-horizontal">
                                <input type="hidden" name="token" value="{{ .Token }}"/>
                                <div class="form-group">
                                    <label for="password" class="col-sm-5 control-label">{{ .i18n_currentPassword }}</label>
                                    <div class="col-sm-5">
                                        <input type="password" autocomplete="current-password" class="form-control" id="password" name="password" required/>
                                    </div>
                                </div>
                                <div class="form-group">
                                    <label for="newPassword" class="col-sm-5 control-label">{{ .i18n_newPassword }}</label>
                                    <div class="col-sm-5">
                                        <input type="password" autocomplete="new-password" pattern="^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.*[^a-zA-Z0-9\x01-\x1f]).{8,}$" title="Must be minimum 8 characters and contain at least one digit, one uppercase letter, one lowercase letter and one special characters" class="form-control" id="newPassword" name="newPassword" required/>
                                    </div>
                                </div>
                                <div class="form-group">
                                    <label for="newPasswordConfirm" class="col-sm-5 control-label">{{ .i18n_newPasswordConfirm }}</label>
                                    <div class="col-sm-5">
                                        <input type="password" autocomplete="new-password" class="form-control" id="newPasswordConfirm" name="newPasswordConfirm" required/>
                                    </div>
                                </div>
                                <div class="form-group">
                                    <div class="col-sm-offset-5 col-sm-5">
                                        <button id="changePassword" type="submit" class="btn btn-lg btn-primary">{{ .i18n_save }}</button>
                                    </div>
                                </div>
                            </form>
                        </div>
                    </div>
                </div>
            </div>
            <div class="modal-footer">
                <button type="button" class="btn btn-default" data-dismiss="modal">{{ .i18n_close }}</button>
            </div>
        </div>
    </div>
</div>
{{ end }}

{{ define "modalSettings.js" }}
function getProfile() {
    $("div#message").addClass('hidden');
    $.getJSON("/user/profile")
        .done(function(data){
            $("span#lastConnection").text(data.LastConnection);
            $("input#userId").val(data.UserProfile.UserID);
            $("input#firstName").val(data.UserProfile.FirstName);
            $("input#lastName").val(data.UserProfile.LastName);
            $("input#emailAddress").val(data.UserProfile.EmailAddress);
        });
}

$("form").submit(function(event) {
    event.preventDefault();
    $.post($(this).attr("action"), $(this).serialize(), function(data) {
            var text, icon;
            switch(data.Result) {
                case 0:
                    text = 'success';
                    icon = 'ok-sign';
                    break;
                case 1:
                    text = 'warning';
                    icon = 'warning-sign';
                    break;
                case 2:
                    text = 'danger';
                    icon = 'ban-circle';
                    break;
            }
            $("div#message").removeClass('hidden').html('<p class="lead text-center text-' + text + '"><span class="glyphicon glyphicon-' + icon + '"></span>&nbsp;<span>' + data.Message + '</span></p>');
            $("#changePasswordForm").trigger("reset")
        }, "json");
});

$('a#openSettings').click(function(event) {
    getProfile();
    $('#modalSettings').modal({ backdrop: 'static', show: true });
});
$('a#showProfile').click(function(event) {
    getProfile();
    $('#tabProfile').removeClass('hidden');
    $('a#showProfile').parent().addClass('active');
    $('#tabChangePassword').addClass('hidden');
    $('a#showChangePassword').parent().removeClass('active')
});
$('a#showChangePassword').click(function(event) {
    $("div#message").addClass('hidden');
    $('#tabProfile').addClass('hidden');
    $('a#showProfile').parent().removeClass('active')
    $('#tabChangePassword').removeClass('hidden');
    $('a#showChangePassword').parent().addClass('active');
});
{{ end }}