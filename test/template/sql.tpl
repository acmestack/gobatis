{{define "selectTestTable"}}
{{$COLUMNS := "`id`, `username`, `password`"}}
SELECT {{$COLUMNS}} FROM `TEST_TABLE`
{{where (ne .UserName "") "AND" "username" .UserName "" | where (ne .Password "pw") "AND" "password" .Password | where (ne .Status -1) "AND" "status" .Status }}
{{end}}

{{define "insertTestTable"}}
{{$COLUMNS := "`id`, `username`, `password`"}}
INSERT INTO `TEST_TABLE` ({{$COLUMNS}})
  VALUES(
  {{.UserName}},
  {{.Password}}
  )
{{end}}

{{define "insertBatchTestTable"}}
{{$COLUMNS := "id, username, password"}}
{{$size := len . | add -1}}
INSERT INTO test_table ({{$COLUMNS}})
  VALUES {{range $i, $v := .}}
  ({{$v.Id}}, '{{$v.UserName}}', '{{$v.Password}}'){{if lt $i $size}},{{end}}
  {{end}}
{{end}}

{{define "updateTestTable"}}
UPDATE `TEST_TABLE`
  {{set (ne .UserName "") "username" .UserName "" | set (ne .Password "") "password" .Password | set (ne .Status -1) "status" .Status}}
{{if (ne .Id 0)}} WHERE `id` = {{.Id}} {{end}}
{{end}}

{{define "deleteTestTable"}}
DELETE FROM `TEST_TABLE`
{{where (ne .Id 0) "AND" "id" .Id "" | where (ne .UserName "") "AND" "username" .UserName | where (ne .Password "pw") "AND" "password" .Password | where (ne .Status -1) "AND" "status" .Status }}
{{end}}

{{template "selectTestTable"}}
{{template "insertTestTable"}}
{{template "updateTestTable"}}
{{template "deleteTestTable"}}