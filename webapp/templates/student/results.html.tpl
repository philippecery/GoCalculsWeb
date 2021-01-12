<!DOCTYPE html>
<html lang="{{ .lang }}">
	{{ template "header.html" . }}
	<div class="col-sm-12">
		<div class="col-sm-4"><h2><span class="text-capitalize">{{ .User.FirstName }} {{ .User.LastName }}</span></h2></div>
		<div class="col-sm-4 text-center"><h2>{{ .i18n_results }}</span></h2></div>
		<div class="col-sm-4 text-right">
			<h2>
				<button type="button" class="btn btn-default bestFilter"><span class="glyphicon glyphicon-education"></span></button>
				<div class="btn-group" role="group">
					<button type="button" class="btn btn-default homeworkTypeFilter active" homeworkType="-1">{{ .i18n_allTypes }}</button>
					<button type="button" class="btn btn-default homeworkTypeFilter" homeworkType="1"><span class="glyphicon glyphicon-hourglass"></span></button>
					<button type="button" class="btn btn-default homeworkTypeFilter" homeworkType="2"><span class="glyphicon glyphicon-pencil"></span></button>
				</div>
				<div class="btn-group" role="group">
					<button type="button" class="btn btn-default statusFilter active" status="-1">{{ .i18n_allStatuses }}</button>
					<button type="button" class="btn btn-success statusFilter" status="3"><span class="glyphicon glyphicon-thumbs-up"></span></button>
					<button type="button" class="btn btn-danger statusFilter" status="2"><span class="glyphicon glyphicon-time"></span></button>
					<button type="button" class="btn btn-danger statusFilter" status="1"><span class="glyphicon glyphicon-thumbs-down"></span></button>
					<button type="button" class="btn btn-danger statusFilter" status="0"><span class="glyphicon glyphicon-remove"></span></button>
				</div>
			</h2>
		</div>
	</div>
	<div class="col-sm-12">
		<div id="resultspage" class="hidden">
			<table id="results" class="table table-striped">
				<thead>
					<tr>
						<th>{{ .i18n_startDate }}</th>
						<th><span class="glyphicon glyphicon-hourglass"></span>/<span class="glyphicon glyphicon-pencil"></span></th>
						<th>+</th>
						<th>-</th> 
						<th>*</th>
						<th>/</th>
						<th>{{ .i18n_duration }}</th>
						<th></th>
						<th></th>
					</tr>
				</thead>
				<tbody id="resultsData">
				</tbody>
			</table>
			<nav id="pagination">
				<ul class="pager">
					<li class="previous"><a href="#"><span aria-hidden="true">&larr;</span> {{ .i18n_previous }}</a></li>
					<li class="next"><a href="#">{{ .i18n_next }} <span aria-hidden="true">&rarr;</span></a></li>
				</ul>
			</nav>
		</div>
		<div id="noresults" class="col-sm-offset-3 col-sm-6 text-center hidden">
			<h3>{{ .i18n_noResults }}</h3>
		</div>
		<div class="col-sm-offset-3 col-sm-6 text-center" id="wait">
			<div class="col-sm-12">
				<h3>{{ .i18n_loading }}</h3>
			</div>
			<div class="col-sm-12">
				<div class="progress">
					<div class="progress-bar progress-bar-default progress-bar-striped active" style="width: 100%;"></div>
				</div>
			</div>
		</div>
	</div>
	<div class="col-sm-offset-3 col-sm-6 text-center">
		<button id="exit" type="button" class="btn btn-lg btn-default btn-block"><span class="glyphicon glyphicon-home"></span>&nbsp;Quitter</button>
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

	{{ template "footer.html" . }}
    <script nonce="{{ .nonce }}" type="text/javascript" src="/js/jquery-2.2.2.min.js"></script>
    <script nonce="{{ .nonce }}" type="text/javascript" src="/js/bootstrap.min.js"></script>
    <script nonce="{{ .nonce }}" type="text/javascript" charset="utf-8">
		$(document).ready(function(){
			var token = '{{ .Token }}';
			var namespace = '/websocket?token='+token;
			var socket = new WebSocket('wss://' + document.domain + ':' + location.port + namespace);
			var nbTotal = 0;
			var homeworkType = -1;
			var status = -1;
			var page = 1;

			function processResultsResponse(msg) {
				var tbody = $('tbody#resultsData');
				tbody.html('');
				var sessions = msg.sessions;
				for (var item in sessions) {
					var result = sessions[item];
					var $tr = $('<tr/>').append(
							$('<td/>').text(result.startTime),
							$('<td/>').html('<span class="glyphicon glyphicon-'+result.type+'"></span>'),
							$('<td/>').text(result.nbAdditions),
							$('<td/>').text(result.nbSubstractions),
							$('<td/>').text(result.nbMultiplications),
							$('<td/>').text(result.nbDivisions),
							$('<td/>').text(result.duration),
							$('<td/>').html('<span class="glyphicon glyphicon-'+result.status+'"></span>'),
							$('<td/>').html('<button type="button" class="btn btn-default details" sessionID="' + result.sessionID + '"><span class="glyphicon glyphicon-th-list"></span></button>'));
					tbody.append($tr);
				}
				nbTotal = msg.nbTotal;
				if (page == 1) {
					$('li.previous').addClass('disabled');
				} else {
					$('li.previous').removeClass('disabled');
				}
				if (page >= (nbTotal/10)) {
					$('li.next').addClass('disabled');
				} else {
					$('li.next').removeClass('disabled');
				}
				$('div#wait').addClass('hidden');
				if (sessions.length > 0) {
					$('div#resultspage').removeClass('hidden');
				} else {
					$('div#noresults').removeClass('hidden');
				}
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
			$('.bestFilter').click(function(event) {
				homeworkType = -1;
				status = -1;
				$('.homeworkTypeFilter').removeClass('active');
				$('.statusFilter').removeClass('active');
				$(this).addClass('active');
				page = 1;
				$('div#resultspage').addClass('hidden');
				$('div#noresults').addClass('hidden');
				$('div#wait').removeClass('hidden');
				socket.send(JSON.stringify({ request: "records", page: page, token: token }));
			});
			$('.statusFilter').click(function(event) {
				status = parseInt($(this).attr('status'));
				$('.statusFilter').removeClass('active');
				$(this).addClass('active');
				page = 1;
				reload();
			});
			$('.homeworkTypeFilter').click(function(event) {
				homeworkType = parseInt($(this).attr('homeworkType'));
				$('.homeworkTypeFilter').removeClass('active');
				$(this).addClass('active');
				page = 1;
				reload();
			});
			$('.previous').click(function(event) {
				if (page > 1) {
					page--;
				}
				reload();
			});
			$('.next').click(function(event) {
				if (page < (nbTotal/10)) {
					page++;
				}
				reload();
			});
			function reload() {
				$('.bestFilter').removeClass('active');
				$('div#resultspage').addClass('hidden');
				$('div#noresults').addClass('hidden');
				$('div#wait').removeClass('hidden');
				socket.send(JSON.stringify({ request: "results", type: homeworkType, status: status, page: page, token: token }));
			}
			$('tbody#resultsData').on('click', 'button.details', function(event) {
				details($(this).attr('sessionID'));
			});
			function details(sessionID) {
				$('div#resultsOperations').text('');
				socket.send(JSON.stringify({ request: "details", sessionID: sessionID, token: token }));
				$('div#results').modal('show');
			}
			$('button#exit').click(function(event) {
				socket.close(1000);
				top.location = "/student/dashboard";
			});
			function log(msg) {
				if(!$('div#logs').hasClass('hidden')) {
					var now = new Date();
					$('div#logs').append('<p>' + now.getHours() + ':' + now.getMinutes() + ':' + now.getSeconds() + ': ' + msg + '</p>');
				}
			}
			socket.onopen = function(event) {
				socket.send(JSON.stringify({ request: "results", type: homeworkType, status: status, page: page, token: token }));
			}

			socket.onmessage = function(event) {
				var msg = JSON.parse(event.data);
				switch(msg.response) {
					case "results":
						processResultsResponse(msg);
						break;
					case "summary":
						processSummaryResponse(msg);
						break;
				}
			}
		});
	</script>
</html>