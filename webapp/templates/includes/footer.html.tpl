{{ define "footer.html" }}
    </div>
  </body>
	<script nonce="{{ .nonce }}">
  document.addEventListener('DOMContentLoaded', function () {
    {{ range $key, $value := .langs }}
    document.getElementById('lang_{{ $key }}').addEventListener('click', function() { language('{{ $key }}') });
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