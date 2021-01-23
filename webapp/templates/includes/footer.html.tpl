{{ define "footer.html" }}
    </div>
  </body>
  <script nonce="{{ .nonce }}" type="text/javascript" src="/js/jquery-2.2.2.min.js"></script>
  <script nonce="{{ .nonce }}" type="text/javascript" src="/js/bootstrap.min.js"></script>
	<script nonce="{{ .nonce }}">
    $(document).ready(function(){
  		$('[data-toggle="tooltip"]').tooltip();

      {{ range $key, $value := .langs }}
      $('#lang_{{ $key }}').click(function(event) { language('{{ $key }}') });
      {{ end }}
    });
    function language(lang) {
      var date = new Date();
      date.setTime(date.getTime() + (10 * 365 * 24 * 60 * 60 * 1000));
      document.cookie = "lang="+lang+"; expires="+date.toGMTString()+"; path=/";
      window.location.reload();
    }
	</script>
{{ end }}