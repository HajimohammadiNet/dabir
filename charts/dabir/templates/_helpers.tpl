{{/*
Expand the name of the chart.
*/}}
{{- define "dabir.name" -}}
{{- default .Chart.Name .Values.global.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "dabir.fullname" -}}
{{- if .Values.global.fullnameOverride -}}
{{- .Values.global.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.global.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Chart label.
*/}}
{{- define "dabir.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Common labels.
*/}}
{{- define "dabir.labels" -}}
helm.sh/chart: {{ include "dabir.chart" . }}
app.kubernetes.io/name: {{ include "dabir.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}

{{/*
API selector labels.
*/}}
{{- define "dabir.api.selectorLabels" -}}
app.kubernetes.io/name: {{ include "dabir.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/component: api
{{- end -}}

{{/*
Web selector labels.
*/}}
{{- define "dabir.web.selectorLabels" -}}
app.kubernetes.io/name: {{ include "dabir.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/component: web
{{- end -}}

{{/*
Service account name.
*/}}
{{- define "dabir.serviceAccountName" -}}
{{- if .Values.serviceAccount.create -}}
{{- default (include "dabir.fullname" .) .Values.serviceAccount.name -}}
{{- else -}}
{{- default "default" .Values.serviceAccount.name -}}
{{- end -}}
{{- end -}}

{{/*
Secret name.
*/}}
{{- define "dabir.secretName" -}}
{{- if .Values.api.secrets.existingSecret -}}
{{- .Values.api.secrets.existingSecret -}}
{{- else -}}
{{- printf "%s-secret" (include "dabir.fullname" .) -}}
{{- end -}}
{{- end -}}

{{/*
API name.
*/}}
{{- define "dabir.api.name" -}}
{{- printf "%s-api" (include "dabir.fullname" .) -}}
{{- end -}}

{{/*
Web name.
*/}}
{{- define "dabir.web.name" -}}
{{- printf "%s-web" (include "dabir.fullname" .) -}}
{{- end -}}