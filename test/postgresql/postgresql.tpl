{{define "selectTestTable"}}
{{$COLUMNS := "id, username, password"}}
SELECT {{$COLUMNS}} FROM test_table WHERE 1=1
        {{if (ne .Username "")}} AND username = '{{.Username}}' {{end}}
        {{if (ne .Password "")}} AND password = '{{.Password}}' {{end}}
{{end}}

{{define "insertTestTable"}}
{{$COLUMNS := "id, username, password"}}
INSERT INTO test_table ({{$COLUMNS}})
  VALUES(
  {{.Id}},
  '{{.Username}}',
  '{{.Password}}'
  )
{{end}}

{{define "updateTestTable"}}
UPDATE test_table SET id = id
  {{if (ne .Username "")}} , username = '{{.Username}}' {{end}}
  {{if (ne .Password "")}} , password = '{{.Password}}' {{end}}
{{if (ne .Id 0)}} WHERE id = {{.Id}} {{end}}
{{end}}

{{define "deleteTestTable"}}
DELETE FROM test_table WHERE 1=1
    {{if (ne .Id 0)}} AND id = {{.Id}} {{end}}
    {{if (ne .Username "")}} AND username = {{.Username}} {{end}}
    {{if (ne .Password "")}} AND password = {{.Password}} {{end}}
{{end}}
