<div align="center">

# TypeID

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

### Maintainers

- [Johannes Gräger](mailto:johannes.graeger@sumup.com)
- [Matouš Dzivjak](mailto:matous.dzivjak@sumup.com)

---

Based on the go implementation of typeid found at: https://github.com/jetify-com/typeid-go by [Jetify](https://www.jetify.com/).
Modifications made available under the same license as the original.

[^UUIDv7]: https://datatracker.ietf.org/doc/html/rfc9562
[^UUIDv4]: https://datatracker.ietf.org/doc/html/rfc4122
