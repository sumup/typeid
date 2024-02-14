module github.com/sumup/typeid/benchmark

go 1.21.6

toolchain go1.21.7

replace github.com/sumup/typeid => ../

require (
	github.com/sumup/typeid v0.0.0-20240207125954-757b87eaff3c
	go.jetpack.io/typeid v1.0.0
)

require (
	github.com/gofrs/uuid/v5 v5.0.0 // indirect
	github.com/jackc/pgx/v5 v5.5.3 // indirect
)
