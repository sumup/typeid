// TypeIDs are a draft standard for *type-safe, globally unique identifiers* based on the upcoming [UUIDv7 standard].
// Their properties, particularly k-sortability, make them suitable primary identifiers for classic database systems like PostgreSQL.
// However, k-sortability may not always be desirable. For instance, you might require an identifier with high randomness entropy for security reasons.
// Additionally, in distributed database systems like CockroachDB, having a k-sortable primary key can lead to hotspots and performance issues.
//
// While this package draws inspiration from the original typeid Go package ([go.jetpack.io/typeid]), it provides multiple ID types:
//
//   - [typeid.Sortable] is based on UUIDv7 and is k-sortable. Its implementation adheres to the draft standard.
//     The suffix part is encoded in **lowercase** crockford base32.
//   - [typeid.Random] is also based on UUIDv4 and is completely random. Unlike `typeid.Sortable`,
//     the suffix part is encoded in **uppercase** crockford base32.
//
// Please refer to the respective type documentation for more details.
//
// # Database Support
//
// ID types in this package can be used with [database/sql] and [github.com/jackc/pgx].
//
// When using the standard library sql, IDs will be stored as their string representation and can be scanned and valued accordingly.
// When using pgx, both TEXT and UUID columns can be used directly. However, note that the type information is lost when using UUID columns, unless you take additional steps
// at the database layer. Be mindful of your identifier semantics, especially in complex JOIN queries.
//
// # Usage
//
// To create a new ID type, define a prefix type that implements the [Prefix] interface. Then, define a TypeAlias for your ID type to [Random] or [Sortable] with your
// prefix type as generic argument.
//
// Example:
//
//	import "github.com/sumup/x/typeid"
//
//	type UserPrefix struct{}
//
//	func (UserPrefix) Prefix() string {
//	    return "user"
//	}
//
//	type UserID = typeid.Sortable[UserPrefix]
//
//	userID, err := typeid.New[UserID]()
//	if err != nil {
//	    fmt.Println("create user id:", err)
//	}
//	fmt.Println(userID) // --> user_01hf98sp99fs2b4qf2jm11hse4
//
// [UUIDv7 standard]: https://www.ietf.org/archive/id/draft-peabody-dispatch-new-uuid-format-01.html#name-versions
//
// [UUIDv4 standard]: https://datatracker.ietf.org/doc/html/rfc4122
package main

//go:generate godoc-readme-gen -f -title "github.com/sumup/typeid"
