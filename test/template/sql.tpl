{{define "selectTestTable"}}
{{$COLUMNS := "`id`, `username`, `password`"}}
SELECT {{$COLUMNS}} FROM `TEST_TABLE` WHERE 1=1
        {{if (ne .UserName "")}} AND `username` = {{.UserName}} {{end}}
        {{if (ne .Password "")}} AND `password` = {{.Password}} {{end}}
{{end}}

{{define "insertTestTable"}}
{{$COLUMNS := "`id`, `username`, `password`"}}
INSERT INTO `TEST_TABLE` ({{$COLUMNS}})
  VALUES(
  {{.UserName}},
  {{.Password}}
  )
{{end}}

{{define "updateTestTable"}}
Update `TEST_TABLE` set
  {{if (ne .UserName "")}} `username` = {{.UserName}} {{end}}
  {{if (ne .Password "")}} `password` = {{.Password}} {{end}}
{{if (ne .Id 0)}} WHERE `id` = {{.Id}} {{end}}
{{end}}

{{define "deleteTestTable"}}
DELETE FROM `TEST_TABLE` WHERE 1=1
    {{if (ne .Id 0)}} AND `id` = {{.Id}} {{end}}
    {{if (ne .UserName "")}} AND `username` = {{.UserName}} {{end}}
    {{if (ne .Password "")}} AND `password` = {{.Password}} {{end}}
{{end}}

{{template "selectTestTable"}}
{{template "insertTestTable"}}
{{template "updateTestTable"}}
{{template "deleteTestTable"}}