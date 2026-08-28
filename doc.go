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
// When using pgx, TEXT and UUID columns can be used directly. With UUID columns the type prefix is not stored in the database
// unless you take additional steps at the database layer.
//
// To keep both the type prefix and UUID in PostgreSQL, define a composite type and register it on your pgx
// connections (for example in AfterConnect). This package only implements the composite field accessors;
// type loading and OID registration stay in the integrating application:
//
//	CREATE TYPE typeid AS (
//	    "type" varchar(63),
//	    "uuid" UUID
//	);
//
//	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
//	    t, err := conn.LoadType(ctx, "typeid")
//	    if err != nil {
//	        return err
//	    }
//	    conn.TypeMap().RegisterType(t)
//	    return nil
//	}
//
// [Sortable] and [Random] implement pgx's CompositeIndexGetter and CompositeIndexScanner interfaces, so pgx's
// CompositeCodec can encode and scan composite typeid columns.
//
// # Usage
//
// To create a new ID type, define a prefix type that implements the [Prefix] interface. Then, define a TypeAlias for your ID type to [Random] or [Sortable] with your
// prefix type as generic argument.
//
// Example:
//
//	import "github.com/sumup/typeid"
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
package typeid

//go:generate godoc-readme-gen -f -title "github.com/sumup/typeid"
