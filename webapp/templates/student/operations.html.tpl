<!DOCTYPE html>
<html lang="fr">
    {{ template "header.html" . }}
    <div class="col-sm-12">
        <div class="col-sm-4"><h2><span class="text-capitalize">{{ .User.FirstName }} {{ .User.LastName }}</span></h2></div>
        <div class="col-sm-4 text-center"><h2><span id="operationName" class="text-capitalize"></span></h2></div>
        <div class="col-sm-4 text-right"><h2><span id="timer"><span id="minutes"></span>:<span id="seconds"></span></span></h2></div>
    </div>
    <div class="col-sm-12">
        <div class="progress percentAll">
            <div id="percentAllGood" class="progress-bar progress-bar-success progress-bar-striped percentGood"></div>
        </div>
    </div>
    <div class="col-sm-10 hidden">
        <div class="progress">
            <div id="percentGood" class="progress-bar progress-bar-info progress-bar-striped percentGood"></div>
        </div>
    </div>
    <div id="nbResults" class="col-sm-2 text-right hidden">
        <span class="glyphicon glyphicon-ok text-success"></span> <span id="nbGood">0</span>
        <span class="glyphicon glyphicon-remove text-danger"></span> <span id="nbWrong">0</span>
    </div>

    <div class="col-sm-12">
        <form id="operation" method="POST" action="#" class="form-inline">
            <div class="form-group form-group-lg col-sm-offset-2 col-sm-8 text-right">
                <div id="answer-input-group" class="input-group input-group-lg">
                    <span id="operation" class="input-group-addon"></span>
                    <input type="text" class="form-control" id="answer" readonly="readonly"/>
                    <span class="input-group-addon"><span id="isGoodAnswer"></span></span>
                </div>
                <div id="answer2-input-group" class="input-group input-group-lg col-sm-7 hidden">
                    <span class="input-group-addon">{{ .i18n_remainder }}</span>
                    <input type="text" class="form-control" id="answer2" readonly="readonly" disabled="disabled"/>
                    <span class="input-group-addon"><span id="isGoodAnswer"></span></span>				
                </div>
            </div>
            <button id="toggleResult" type="button" class="btn btn-lg btn-default hidden">
                <span id="toggleResult" class="glyphicon glyphicon-eye-open"></span>
            </button>
        </form>
    </div>

    <div class="col-sm-12">
        <hr/>
    </div>

    {{ template "operations.keyboard.html" .Keyboards }}

    <div class="col-sm-offset-3 col-sm-6 text-center">
        <button id="submit" type="button" class="btn btn-lg btn-default btn-block"><span class="glyphicon glyphicon-sunglasses"></span>&nbsp;{{ .i18n_check }}</button>
        <div id="next" class="btn-group btn-group-justified hidden" role="group">
            <div class="btn-group" role="group">
                <button id="next" type="button" class="btn btn-lg btn-success btn-block"><span class="glyphicon glyphicon-step-forward"></span>&nbsp;{{ .i18n_continue }}</button>
            </div>
            <div class="btn-group" role="group">
                <button id="restart" type="button" class="btn btn-lg btn-danger btn-block"><span class="glyphicon glyphicon-fast-backward"></span>&nbsp;{{ .i18n_restart }}</button>
            </div>
        </div>
        <div id="end" class="btn-group btn-group-justified hidden" role="group">
            <div class="btn-group" role="group">
                <button id="summary" type="button" class="btn btn-lg btn-primary btn-block"><span class="glyphicon glyphicon-list"></span>&nbsp;{{ .i18n_results }}</button>
            </div>
            <div class="btn-group" role="group">
                <button id="retry" type="button" class="btn btn-lg btn-danger btn-block"><span class="glyphicon glyphicon-fast-backward"></span>&nbsp;{{ .i18n_retry }}</button>
            </div>
        </div>
    </div>
    <div class="col-sm-12">
        <hr/>
    </div>
    <div class="col-sm-offset-3 col-sm-6 text-center">
        <button id="exit" type="button" class="btn btn-lg btn-default btn-block"><span class="glyphicon glyphicon-home"></span>&nbsp;{{ .i18n_quit }}</button>
    </div>

    <div class="modal fade" id="results" tabindex="-1" role="dialog">
        <div class="modal-dialog" role="document">
            <div class="modal-content">
                <div class="modal-header">
                    <h4 class="modal-title">{{ .i18n_results }}</h4>
                </div>
                <div class="modal-body" id="operationlist">
                    <div id="resultsOperations"></div>
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-default" data-dismiss="modal">{{ .i18n_close }}</button>
                </div>
            </div>
        </div>
    </div>

    <div id="logs" class="col-sm-12 hidden">
    </div>

    </div>
    </body>
    <script nonce="{{ .nonce }}" type="text/javascript" src="/js/jquery-2.2.2.min.js"></script>
    <script nonce="{{ .nonce }}" type="text/javascript" src="/js/bootstrap.min.js"></script>
    <script nonce="{{ .nonce }}" type="text/javascript" charset="utf-8">
        $(document).ready(function(){
            var token = '{{ .Token }}';
            var namespace = '/student/websocket?token='+token;
            var socket = new WebSocket('wss://' + document.domain + ':' + location.port + namespace);
            var type = {{ .TypeID }};

            function disableLink(id) {
                $('#'+id).addClass('disabled');
                $('#'+id).click(function(event) { window.alert(`{{ .i18n_errDisabled }}`) });
            }
            disableLink('profile');
            {{ range $key, $value := .langs }}
            disableLink('lang_{{ $key }}');
            {{ end }}

            $('div#keyboard button[id^=keynum]').click(function(event) {
                log(this.innerText);
                $('input#answer').val($('input#answer').val() + this.innerText);
            });
            $('div#keyboard button#keydel').click(function(event) {
                log('DEL');
                var value = $('input#answer').val();
                $('input#answer').val(value.substring(0, value.length-1));
            });
            $('div#keyboard2 button[id^=keynum]').click(function(event) {
                log(this.innerText);
                $('input#answer2').val($('input#answer2').val() + this.innerText);
            });
            $('div#keyboard2 button#keydel').click(function(event) {
                log('DEL');
                var value = $('input#answer2').val();
                $('input#answer2').val(value.substring(0, value.length-1));
            });

            function getTimeRemaining(endtime) {
                var t = Date.parse(endtime) - Date.parse(new Date());
                var seconds = Math.floor((t / 1000) % 60);
                var minutes = Math.floor((t / 1000 / 60) % 60);
                return {
                'total': t,
                'minutes': minutes,
                'seconds': seconds
                };
            }
            var ended = false;
            var timeout = false;
            var duration = {{ .Homework.Time }} * 60 * 1000;
            var deadline = new Date(Date.parse(new Date()) + duration);
            function updateClock() {
                var t = getTimeRemaining(deadline);
                if(t.minutes == 1 & t.seconds == 0) {
                    $('span#timer').addClass('text-danger');
                }
                $('span#minutes').text(('0' + t.minutes).slice(-2));
                $('span#seconds').text(('0' + t.seconds).slice(-2));
                if (t.total <= 0 || ended) {
                    clearInterval(timeinterval);
                    timeout = (t.total <= 0);
                    $('input#answer').val('');
                    $('div#keyboard').addClass('hidden');
                    $('div#keyboard2').addClass('hidden');
                    $('input#answer2').val('');
                    $('button#submit').addClass('hidden');
                    $('div#next').addClass('hidden');
                    $('div#end').removeClass('hidden');
                    end();
                }
            }
            updateClock();
            var timeinterval = setInterval(updateClock, 1000);

            $('button#next').click(function(event) {
                socket.send(JSON.stringify({ request: "operation", token: token }));
            });

            function processOperationResponse(msg) {
                $('div#keyboard').removeClass('hidden');
                $('div#answer-input-group').removeClass('has-success has-warning has-error');
                $('span#isGoodAnswer').removeClass().addClass('glyphicon glyphicon-question-sign');
                $('button#toggleResult').addClass('hidden');
                $('span#toggleResult').removeClass(). addClass('glyphicon glyphicon-eye-open');
                $('span#operationName').text(msg.operationName);
                $('span#operation').text(msg.operand1 + ' ' + msg.operator + ' ' + msg.operand2 + ' = ').html();
                $('input#answer').val('');
                $('input#answer2').val('');
                $('div#answer2-input-group').removeClass('has-success has-warning has-error');
                if(msg.operator == '/') {
                    $('div#keyboard2').removeClass('hidden');
                    $('input#answer2').prop('disabled', false);
                    $('div#answer2-input-group').removeClass('hidden');
                } else {
                    $('div#keyboard2').addClass('hidden');
                    $('input#answer2').prop('disabled', true);
                    $('div#answer2-input-group').addClass('hidden');
                }
                $('button#submit').prop('disabled', false).removeClass('hidden');
                $('div#next').addClass('hidden');
                $('div#end').addClass('hidden');
            }
            $('button#submit').click(function(event) {
                log('button "submit" clicked');
                if(!$('input#answer').val() || (!$('input#answer2').prop('disabled') && !$('input#answer2').val())) {
                    if(!$('input#answer').val()) {
                        $('div#answer-input-group').removeClass('has-success has-warning has-error').addClass('has-warning');
                    } else {
                        $('div#keyboard').addClass('hidden');
                    }
                    if(!$('input#answer2').prop('disabled') && !$('input#answer2').val()) {
                        $('div#answer2-input-group').removeClass('has-success has-warning has-error').addClass('has-warning');
                    } else {
                        $('div#keyboard2').addClass('hidden');
                    }
                } else {
                    $('div#keyboard').addClass('hidden');
                    $('div#keyboard2').addClass('hidden');
                    $('button#submit').prop('disabled', true);
                    socket.send(JSON.stringify({
                        request: "answer",
                        answer: $('input#answer').val(),
                        answer2: $('input#answer2').val(),
                        token: token
                    }));
                }
            });

            function processAnswerResponse(msg) {
                log('answerResponse');
                log('msg.percentUpdate = ' + msg.percentUpdate);
                if(msg.good == true) {
                    $('span#nbGood').text(msg.nbUpdate);
                    $('div#percentGood').css('width', msg.percentUpdate + '%');
                    $('div#percentAllGood').css('width', msg.percentAll + '%');
                } else {
                    $('div#answer-input-group').removeClass('has-success has-warning has-error').addClass('has-error');
                    $('div#answer2-input-group').removeClass('has-success has-warning has-error').addClass('has-error');
                    $('span#isGoodAnswer').removeClass().addClass('glyphicon glyphicon-remove');
                    $('span#nbWrong').text(msg.nbUpdate);
                    $('button#toggleResult').removeClass('hidden');
                    if($('div#percentAllGood').hasClass('progress-bar-success')) {
                        $('div#percentAllGood').removeClass('progress-bar-success').addClass('progress-bar-danger');
                    }
                }
                //log('msg.nbRemaining = ' + msg.nbRemaining);
                log('msg.nbTotalRemaining = ' + msg.nbTotalRemaining);
                $('button#submit').addClass('hidden');
                if(msg.nbTotalRemaining > 0) {
                    if(msg.good == true) {
                    	socket.send(JSON.stringify({ request: "operation", token: token }));
                    } else {
                        if(type == 1) {
                            ended = true;
                            updateClock();
                        } else {
                            $('div#next').removeClass('hidden');
                        }
                    }
                } else {
                    ended = true;
                    updateClock();
                }
            }

            $('button#toggleResult').click(function(event) {
                socket.send(JSON.stringify({ request: "toggle", show: $('span#toggleResult').hasClass('glyphicon-eye-open'), token: token }));
            });

            function processToggleResponse(msg) {
                if(msg.showResult) {
                    $('div#answer-input-group').removeClass('has-success has-warning has-error').addClass('has-success');
                    $('div#answer2-input-group').removeClass('has-success has-warning has-error').addClass('has-success');
                    $('span#toggleResult').removeClass(). addClass('glyphicon glyphicon-eye-close');
                    $('input#answer').val(msg.result);
                    $('input#answer2').val(msg.result2);
                } else {
                    $('div#answer-input-group').removeClass('has-success has-warning has-error').addClass('has-error');
                    $('div#answer2-input-group').removeClass('has-success has-warning has-error').addClass('has-error');
                    $('span#toggleResult').removeClass(). addClass('glyphicon glyphicon-eye-open');
                    $('input#answer').val(msg.answer);
                    $('input#answer2').val(msg.answer2);
                }
            }
            $('button#summary').click(function(event) {
                $('#results').modal('show');
            });
            function end() {
                if (timeout) {
                    $('div#resultsOperations').append(
                        $('<p class="lead text-center text-danger"/>').append(
                            $('<span class="glyphicon glyphicon-time"/>').text(''),
                            $('<span/>').text('{{ .i18n_timeout }}'),
                            $('<span class="glyphicon glyphicon-time"/>').text('')
                        )
                    )
                }
                socket.send(JSON.stringify({ request: "end", timeout: timeout, token: token }));
                $('#results').modal('show');
            }

            function processSummaryResponse(msg) {
                log(msg.operationName + ' msg.nbTotal = ' + msg.nbTotal);
                if(msg.nbTotal > 0) {
                    var textClass, comment;
                    if(msg.nbGood < msg.nbTotal) {
                        if(msg.nbWrong > 0) {
                            textClass = 'text-danger';
                        } else {
                            textClass = 'text-warning';
                        }
                    } else {
                        if(msg.nbWrong > 0) {
                            textClass = 'text-warning';
                        } else {
                            textClass = 'text-success';
                        }
                    }
                    $('div#resultsOperations').append(
                        $('<p class="text-center ' + textClass + '"/>').html(msg.operationsTodo + ', ' + msg.operationsGood + ', ' + msg.operationsWrong)
                    );
                    if(msg.success) {
                        $('span#resultsOperations').append('<p class="lead text-center text-success"><span class="glyphicon glyphicon glyphicon-thumbs-up"></span><span>&nbsp;{{ .i18n_success }}&nbsp;</span><span class="glyphicon glyphicon glyphicon-thumbs-up"></span></p>');
                    } else {
                        $('span#resultsOperations').append('<p class="lead text-center text-danger"><span class="glyphicon glyphicon-thumbs-down"></span><span>&nbsp;{{ .i18n_failure }}&nbsp;</span><span class="glyphicon glyphicon-thumbs-down"></span></p>');
                    }
                }
            }

            function pluriel(mot, nombre) {
                if(nombre>1) {
                    return mot+'s';
                }
                return mot;
            }

            $('button#restart').click(function(event) {
                restart();
            });
            $('button#retry').click(function(event) {
                restart();
            });
            function restart() {
                socket.close(1000);
                top.location = "/student/operations?type={{ .TypeID }}";
            }
            $('button#exit').click(function(event) {
                socket.send(JSON.stringify({ request: "end", token: token }));
                socket.close(1000);
                top.location = "/student/dashboard";
            });

            function log(msg) {
                if(!$('div#logs').hasClass('hidden')) {
                    var now = new Date();
                    $('div#logs').append('<p>' + now.getHours() + ':' + now.getMinutes() + ':' + now.getSeconds() + ': ' + msg + '</p>');
                }
            }
            if(type==2) {
                $('nbResults').removeClass('hidden');
            }

            socket.onopen = function(event) {
            	socket.send(JSON.stringify({ request: "operation", token: token }));
            }

            socket.onmessage = function(event) {
                var msg = JSON.parse(event.data);
                switch(msg.response) {
                    case "operation":
                        processOperationResponse(msg);
                        break;
                    case "answer":
                        processAnswerResponse(msg);
                        break;
                    case "toggle":
                        processToggleResponse(msg);
                        break;
                    case "summary":
                        processSummaryResponse(msg);
                        break;
                }
            }

            $(window).unload(function() {
                socket.send(JSON.stringify({ request: "end", token: token }));
                socket.close(1000);
            });
        });
    </script>
</html>