{{/* basic helpers — you can improve these */}}
{{- define "fullname" -}}
{{- printf "%s-%s" .Release.Name .Chart.Name -}}
{{- end -}}