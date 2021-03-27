{{/*
Expand the name of the chart.
*/}}
{{- define "gotway.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "gotway.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "gotway.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "gotway.commonLabels" -}}
helm.sh/chart: {{ include "gotway.chart" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "gotway.selectorLabels" -}}
app.kubernetes.io/name: {{ include "gotway.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Labels
*/}}
{{- define "gotway.labels" -}}
{{ include "gotway.commonLabels" . }}
{{ include "gotway.selectorLabels" . }}
{{- end }}

{{/*
Full name Catalog
*/}}
{{- define "gotway.fullnameCatalog" -}}
{{- printf "%s-%s" (include "gotway.fullname" .) "catalog" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Selector labels Catalog
*/}}
{{- define "gotway.selectorLabelsCatalog" -}}
app.kubernetes.io/name: {{ include "gotway.fullnameCatalog" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Labels Catalog
*/}}
{{- define "gotway.labelsCatalog" -}}
{{ include "gotway.commonLabels" . }}
{{ include "gotway.selectorLabelsCatalog" . }}
{{- end }}

{{/*
Full name Route
*/}}
{{- define "gotway.fullnameRoute" -}}
{{- printf "%s-%s" (include "gotway.fullname" .) "route" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Selector labels Route
*/}}
{{- define "gotway.selectorLabelsRoute" -}}
app.kubernetes.io/name: {{ include "gotway.fullnameRoute" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Labels Route
*/}}
{{- define "gotway.labelsRoute" -}}
{{ include "gotway.commonLabels" . }}
{{ include "gotway.selectorLabelsRoute" . }}
{{- end }}

{{/*
Full name Stock
*/}}
{{- define "gotway.fullnameStock" -}}
{{- printf "%s-%s" (include "gotway.fullname" .) "stock" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Selector labels Stock
*/}}
{{- define "gotway.selectorLabelsStock" -}}
app.kubernetes.io/name: {{ include "gotway.fullnameStock" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Labels Stock
*/}}
{{- define "gotway.labelsStock" -}}
{{ include "gotway.commonLabels" . }}
{{ include "gotway.selectorLabelsStock" . }}
{{- end }}
