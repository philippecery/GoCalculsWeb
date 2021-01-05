<!DOCTYPE html>
<html lang="{{ .lang }}">
	{{ template "header.html" . }}
	{{ template "teacher.navbar.html" . }}
    <div class="col-sm-12 text-center"><h2><span class="glyphicon glyphicon-tasks"></span>&nbsp;<span>{{ .i18n_gradeFormTitle }}</span></h2></div>
	{{ template "error.html" . }}
    <div class="col-sm-12">
        <form id="gradeForm" method="POST" action="/teacher/grade/save" class="form-horizontal">
            <input type="hidden" name="token" value="{{ .Token }}"/>
            <input type="hidden" name="operation" value="{{ .Operation }}"/>
            <input type="hidden" name="gradeID" value="{{ .Grade.GradeID }}"/>
            <div class="form-group">
                <label for="name" class="col-sm-2 control-label">{{ .i18n_gradeName }}</label>
                <div class="col-sm-10">
                    <input type="text" class="form-control" id="name" name="name" maxlength="32" value="{{ .Grade.Name }}" autofocus required/>
                </div>
            </div>
            <div class="form-group">
                <label for="description" class="col-sm-2 control-label">{{ .i18n_gradeDescription }}</label>
                <div class="col-sm-10">
                    <textarea class="form-control" id="description" name="description" rows="2">{{ .Grade.Description }}</textarea>
                </div>
            </div>
            <label class="col-sm-offset-2 control-label">{{ .i18n_mentalmath }}</label>
            <div class="form-group">
				<label for="mm_nbAdditions" class="col-sm-offset-3 control-label">{{ .i18n_nbAdditions }}:&nbsp;</label><span id="mm_nbAdditionsDisplay">{{ .Grade.MentalMath.NbAdditions }}</span>
                <div class="col-sm-offset-3 col-sm-9">
                    <input type="range" class="form-control" id="mm_nbAdditions" name="mm_nbAdditions" min="0" max="100" step="10" value="{{ .Grade.MentalMath.NbAdditions }}"/>
                </div>
            </div>
            <div class="form-group">
				<label for="mm_nbSubstractions" class="col-sm-offset-3 control-label">{{ .i18n_nbSubstractions }}:&nbsp;</label><span id="mm_nbSubstractionsDisplay">{{ .Grade.MentalMath.NbSubstractions }}</span>
                <div class="col-sm-offset-3 col-sm-9">
                    <input type="range" class="form-control" id="mm_nbSubstractions" name="mm_nbSubstractions" min="0" max="100" step="10" value="{{ .Grade.MentalMath.NbAdditions }}"/>
                </div>
            </div>
            <div class="form-group">
				<label for="mm_nbMultiplications" class="col-sm-offset-3 control-label">{{ .i18n_nbMultiplications }}:&nbsp;</label><span id="mm_nbMultiplicationsDisplay">{{ .Grade.MentalMath.NbMultiplications }}</span>
                <div class="col-sm-offset-3 col-sm-9">
                    <input type="range" class="form-control" id="mm_nbMultiplications" name="mm_nbMultiplications" min="0" max="100" step="10" value="{{ .Grade.MentalMath.NbMultiplications }}"/>
                </div>
            </div>
            <div class="form-group">
				<label for="mm_nbDivisions" class="col-sm-offset-3 control-label">{{ .i18n_nbDivisions }}:&nbsp;</label><span id="mm_nbDivisionsDisplay">{{ .Grade.MentalMath.NbDivisions }}</span>
                <div class="col-sm-offset-3 col-sm-9">
                    <input type="range" class="form-control" id="mm_nbDivisions" name="mm_nbDivisions" min="0" max="100" step="10" value="{{ .Grade.MentalMath.NbDivisions }}"/>
                </div>
            </div>
            <div class="form-group">
				<label for="mm_time" class="col-sm-offset-3 control-label">{{ .i18n_timeInMinutes }}:&nbsp;</label><span id="mm_timeDisplay">{{ .Grade.MentalMath.Time }}</span>
                <div class="col-sm-offset-3 col-sm-9">
                    <input type="range" class="form-control" id="mm_time" name="mm_time" min="1" max="10" step="1" value="{{ .Grade.MentalMath.Time }}"/>
                </div>
            </div>
            <label class="col-sm-offset-2 control-label">{{ .i18n_columnform }}</label>
            <div class="form-group">
				<label for="cf_nbAdditions" class="col-sm-offset-3 control-label">{{ .i18n_nbAdditions }}:&nbsp;</label><span id="cf_nbAdditionsDisplay">{{ .Grade.ColumnForm.NbAdditions }}</span>
                <div class="col-sm-offset-3 col-sm-9">
                    <input type="range" class="form-control" id="cf_nbAdditions" name="cf_nbAdditions" min="0" max="10" step="1" value="{{ .Grade.ColumnForm.NbAdditions }}"/>
                </div>
            </div>
            <div class="form-group">
				<label for="cf_nbSubstractions" class="col-sm-offset-3 control-label">{{ .i18n_nbSubstractions }}:&nbsp;</label><span id="cf_nbSubstractionsDisplay">{{ .Grade.ColumnForm.NbSubstractions }}</span>
                <div class="col-sm-offset-3 col-sm-9">
                    <input type="range" class="form-control" id="cf_nbSubstractions" name="cf_nbSubstractions" min="0" max="10" step="1" value="{{ .Grade.ColumnForm.NbSubstractions }}"/>
                </div>
            </div>
            <div class="form-group">
				<label for="cf_nbMultiplications" class="col-sm-offset-3 control-label">{{ .i18n_nbMultiplications }}:&nbsp;</label><span id="cf_nbMultiplicationsDisplay">{{ .Grade.ColumnForm.NbMultiplications }}</span>
                <div class="col-sm-offset-3 col-sm-9">
                    <input type="range" class="form-control" id="cf_nbMultiplications" name="cf_nbMultiplications" min="0" max="10" step="1" value="{{ .Grade.ColumnForm.NbMultiplications }}"/>
                </div>
            </div>
            <div class="form-group">
				<label for="cf_nbDivisions" class="col-sm-offset-3 control-label">{{ .i18n_nbDivisions }}:&nbsp;</label><span id="cf_nbDivisionsDisplay">{{ .Grade.ColumnForm.NbDivisions }}</span>
                <div class="col-sm-offset-3 col-sm-9">
                    <input type="range" class="form-control" id="cf_nbDivisions" name="cf_nbDivisions" min="0" max="10" step="1" value="{{ .Grade.ColumnForm.NbDivisions }}"/>
                </div>
            </div>
            <div class="form-group">
				<label for="cf_time" class="col-sm-offset-3 control-label">{{ .i18n_timeInMinutes }}:&nbsp;</label><span id="cf_timeDisplay">{{ .Grade.ColumnForm.Time }}</span>
                <div class="col-sm-offset-3 col-sm-9">
                    <input type="range" class="form-control" id="cf_time" name="cf_time" min="5" max="60" step="5" value="{{ .Grade.ColumnForm.Time }}"/>
                </div>
            </div>
            <div class="form-group">
                <div class="col-sm-offset-2 col-sm-10">
                    <div class="btn-group" role="group">
                        <button id="save" type="submit" class="btn btn-lg btn-primary">{{ .i18n_save }}</button>
                        <a class="btn btn-lg btn-default" href="/teacher/grade/list" role="button"><span>{{ .i18n_cancel }}</span></a>
                    </div>
                </div>
            </div>
        </form>
    </div>
	{{ template "footer.html" . }}
	<script nonce="{{ .nonce }}" type="text/javascript" src="/js/jquery-2.2.2.min.js"></script>
    <script nonce="{{ .nonce }}" type="text/javascript" charset="utf-8">
        $(document).ready(function(){
			$('input#mm_nbAdditions').on('input', function(event) {
				$('span#mm_nbAdditionsDisplay').text(this.value);
			});
			$('input#mm_nbSubstractions').on('input', function(event) {
				$('span#mm_nbSubstractionsDisplay').text(this.value);
			});
			$('input#mm_nbMultiplications').on('input', function(event) {
				$('span#mm_nbMultiplicationsDisplay').text(this.value);
			});
			$('input#mm_nbDivisions').on('input', function(event) {
				$('span#mm_nbDivisionsDisplay').text(this.value);
			});
			$('input#mm_time').on('input', function(event) {
				$('span#mm_timeDisplay').text(this.value);
			});
			$('input#cf_nbAdditions').on('input', function(event) {
				$('span#cf_nbAdditionsDisplay').text(this.value);
			});
			$('input#cf_nbSubstractions').on('input', function(event) {
				$('span#cf_nbSubstractionsDisplay').text(this.value);
			});
			$('input#cf_nbMultiplications').on('input', function(event) {
				$('span#cf_nbMultiplicationsDisplay').text(this.value);
			});
			$('input#cf_nbDivisions').on('input', function(event) {
				$('span#cf_nbDivisionsDisplay').text(this.value);
			});
			$('input#cf_time').on('input', function(event) {
				$('span#cf_timeDisplay').text(this.value);
			});
        });
	</script>
</html>