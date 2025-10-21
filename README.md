<div align="center">

# TypeID

[![Stars](https://img.shields.io/github/stars/sumup/typeid?style=social)](https://github.com/sumup/typeid/)
[![Go Reference](https://pkg.go.dev/badge/github.com/sumup/typeid.svg)](https://pkg.go.dev/github.com/sumup/typeid)
[![CI Status](https://github.com/sumup/typeid/workflows/CI/badge.svg)](https://github.com/sumup/typeid/actions/workflows/ci.yml)
[![License](https://img.shields.io/github/license/sumup/typeid)](./LICENSE)
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-v2.1%20adopted-ff69b4.svg)](https://github.com/sumup/typeid/tree/main/CODE_OF_CONDUCT.md)

</div>

TypeIDs are a draft standard for *type-safe, globally unique identifiers* based on the [UUIDv7 standard](https://datatracker.ietf.org/doc/html/rfc9562). Their properties, particularly k-sortability, make them suitable primary identifiers for classic database systems like PostgreSQL. However, k-sortability may not always be desirable. For instance, you might require an identifier with high randomness entropy for security reasons. Additionally, in distributed database systems like CockroachDB, having a k-sortable primary key can lead to hotspots and performance issues.

While this package draws inspiration from the original typeid-go package ([github.com/jetify-com/typeid-go](https://github.com/jetify-com/typeid-go)), it provides multiple ID types:

- `typeid.Sortable` is based on UUIDv7[^UUIDv7] and is k-sortable. Its implementation adheres to the draft standard. The suffix part is encoded in **lowercase** crockford base32.
- `typeid.Random` is also based on UUIDv4[^UUIDv4] and is completely random. Unlike `typeid.Sortable`, the suffix part is encoded in **uppercase** crockford base32.

Please refer to the respective type documentation for more details.

## Install

```shell
go get github.com/sumup/typeid
```

# Usage

To create a new ID type, define a prefix type that implements the `typeid.Prefix` interface. Then, define a TypeAlias for your ID type to `typeid.Random` or `typeid.Sortable` with your prefix type as generic argument.

Example:

```go
import "github.com/sumup/typeid"

type UserPrefix struct{}

func (UserPrefix) Prefix() string {
    return "user"
}

type UserID = typeid.Sortable[UserPrefix]

userID, err := typeid.New[UserID]()
if err != nil {
    fmt.Println("create user id:", err)
}

fmt.Println(userID) // --> user_01hf98sp99fs2b4qf2jm11hse4
```

# Database Support

ID types in this package can be used with [database/sql](https://pkg.go.dev/database/sql) and [github.com/jackc/pgx](https://pkg.go.dev/github.com/jackc/pgx/v5).

When using the standard library SQL, IDs will be stored as their string representation and can be scanned and valued accordingly. When using pgx, both TEXT and UUID columns can be used directly. However, note that the type information is lost when using UUID columns, unless you take additional steps at the database layer. Be mindful of your identifier semantics, especially in complex JOIN queries.

## Using with sqlc

TypeIDs work seamlessly with [sqlc](https://sqlc.dev/) by using column overrides in your `sqlc.yaml` configuration:

```yaml
version: "2"
sql:
  - schema: "schema.sql"
    queries: "queries"
    engine: postgresql
    gen:
      go:
        package: postgres
        out: postgres
        sql_package: pgx/v5
        overrides:
          - column: users.id
            go_type:
              import: github.com/yourorg/yourproject/internal/domain
              type: "UserID"
```

With this configuration, sqlc will generate Go code that uses your TypeID type directly:

```go
package domain

import "github.com/sumup/typeid"

type UserPrefix struct{}

func (UserPrefix) Prefix() string {
    return "user"
}

type UserID = typeid.Sortable[UserPrefix]

type User struct {
    ID          UserID
    Name        string
    // ... other fields
}
```

You can then use your TypeID types directly in your queries:

```sql
-- name: CreateUser :one
INSERT INTO users (
    id,
    name
) VALUES (
    @id,
    @name
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = @id;
```

And call them from Go:

```go
userID := typeid.Must(typeid.New[domain.UserID]())
user, err := queries.CreateUser(ctx, postgres.CreateUserParams{
    ID:          userID,
    Name:        "Karl",
})
```

## Using with oapi-codegen

TypeIDs can be used with [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) to generate type-safe API clients and servers. Use the `x-go-type` and `x-go-type-import` extensions in your OpenAPI specification:

```yaml
paths:
  /users/{user_id}:
    parameters:
      - in: path
        name: user_id
        description: The ID of the uesr to retrieve.
        required: true
        schema:
          type: string
          example: user_01hf98sp99fs2b4qf2jm11hse4
          x-go-type: "domain.UserID"
          x-go-type-import:
            path: github.com/yourorg/yourproject/internal/domain

components:
  schemas:
    User:
      type: object
      properties:
        id:
          type: string
          description: Unique identifier of the user.
          example: user_01hf98sp99fs2b4qf2jm11hse4
          x-go-type: "domain.UserID"
          x-go-type-import:
            path: github.com/yourorg/yourproject/internal/domain
        name:
          type: string
```

The generated code will use your TypeID types:

```go
type User struct {
    Id   domain.UserID `json:"id"`
    Name string        `json:"name"`
}
```

TypeIDs implement `encoding.TextMarshaler` and `encoding.TextUnmarshaler`, so they work with JSON encoding/decoding in generated API code without any additional configuration.

### Maintainers

- [Johannes Gräger](mailto:johannes.graeger@sumup.com)
- [Matouš Dzivjak](mailto:matous.dzivjak@sumup.com)

---

Based on the go implementation of typeid found at: https://github.com/jetify-com/typeid-go by [Jetify](https://www.jetify.com/).
Modifications made available under the same license as the original.

[^UUIDv7]: https://datatracker.ietf.org/doc/html/rfc9562
[^UUIDv4]: https://datatracker.ietf.org/doc/html/rfc4122
