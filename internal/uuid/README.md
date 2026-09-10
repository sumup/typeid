# Vendored Go UUID implementation

`uuid.go` is copied from Go **1.27.1**:

- https://github.com/golang/go/blob/go1.27.1/src/uuid/uuid.go

The BSD 3-Clause license is reproduced in `LICENSE`. The implementation is
unchanged.
Keep local compatibility code outside this directory so future updates can be
compared directly with upstream.

This copy allows typeid to keep supporting Go 1.25 without a UUID module
dependency. It uses only APIs available in Go 1.25.
