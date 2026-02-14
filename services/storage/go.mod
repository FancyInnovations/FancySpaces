module github.com/fancyinnovations/fancyspaces/storage

go 1.25.5

require (
	github.com/OliverSchlueter/goutils v0.0.28
	github.com/justinas/alice v1.2.0
	github.com/mattn/go-sqlite3 v1.14.33
)

require (
	github.com/klauspost/compress v1.18.3 // indirect
	github.com/nats-io/nats.go v1.48.0 // indirect
	github.com/nats-io/nkeys v0.4.12 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	golang.org/x/crypto v0.47.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
)

replace github.com/fancyinnovations/fancyspaces/integrations/storage-go-sdk => ../integrations/storage-go-sdk