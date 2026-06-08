module github.com/sumup/typeid/benchmark

go 1.25.0

replace github.com/sumup/typeid => ../

require (
	github.com/sumup/typeid v0.0.0-20240207125954-757b87eaff3c
	go.jetify.com/typeid v1.1.0
)

require (
	github.com/gofrs/uuid/v5 v5.4.0 // indirect
	github.com/jackc/pgx/v5 v5.9.0 // indirect
)
