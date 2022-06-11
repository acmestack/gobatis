{{define "namespace"}}test{{end}}

{{define "selectTestTable"}}
SELECT id,username,password,createtime FROM test_table
{{where .Id "AND" "id = " (arg .Id) "" | where .Username "AND" "username = " (arg .Username) | where .Password "AND" "password = " (arg .Password) | where .Createtime "AND" "createtime = " (arg .Createtime)}}
{{end}}

{{define "selectTestTableCount"}}
SELECT COUNT(*) FROM test_table
{{where .Id "AND" "id = " (arg .Id) "" | where .Username "AND" "username = " (arg .Username) | where .Password "AND" "password = " (arg .Password) | where .Createtime "AND" "createtime = " (arg .Createtime)}}
{{end}}

{{define "insertTestTable"}}
INSERT INTO test_table(id,username,password,createtime)
VALUES(
{{arg .Id}}, {{arg .Username}}, {{arg .Password}}, {{arg .Createtime}})
{{end}}

{{define "insertBatchTestTable"}}
{{$size := len . | add -1}}
INSERT INTO test_table(id,username,password,createtime)
VALUES {{range $i, $v := .}}
({{arg $v.Id}}, {{arg $v.Username}}, {{arg $v.Password}}, {{arg $v.Createtime}}){{if lt $i $size}},{{end}}
{{end}}
{{end}}

{{define "updateTestTable"}}
UPDATE test_table
{{set .Id "id = " (arg .Id) "" | set .Username "username = " (arg .Username) | set .Password "password = " (arg .Password) | set .Createtime "createtime = " (arg .Createtime)}}
{{where .Id "AND" "id = " (arg .Id) ""}}
{{end}}

{{define "deleteTestTable"}}
DELETE FROM test_table
{{where .Id "AND" "id = " (arg .Id) "" | where .Username "AND" "username = " (arg .Username) | where .Password "AND" "password = " (arg .Password) | where .Createtime "AND" "createtime = " (arg .Createtime)}}
{{end}}