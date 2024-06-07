{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "browserkube.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "browserkube.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "browserkube.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Generate labels
*/}}
{{- define "labels" }}
heritage: {{ $.Release.Service | quote }}
release: {{ $.Release.Name | quote }}
chart: {{ include "browserkube.chart" . }}
app: {{ $.Chart.Name | quote }}
{{- end -}}

Create the name of the service account to use
*/}}
{{- define "browserkube.serviceAccountName" -}}
{{- if .Values.rbac.serviceAccount.create -}}
    {{ default "browserkube" .Values.rbac.serviceAccount.serviceAccountName }}
{{- else -}}
    {{ default "default" .Values.rbac.serviceAccount.serviceAccountName }}
{{- end -}}
{{- end -}}
