module github.com/dexidp/dex

go 1.15

require (
	github.com/beevik/etree v1.1.0
	github.com/coreos/go-oidc v2.2.1+incompatible
	github.com/dexidp/dex/api/v2 v2.0.0
	github.com/felixge/httpsnoop v1.0.1
	github.com/ghodss/yaml v1.0.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/golang/protobuf v1.3.2
	github.com/gorilla/handlers v1.4.2
	github.com/gorilla/mux v1.7.3
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/kylelemons/godebug v1.1.0
	github.com/lib/pq v1.3.0
	github.com/masami10/dex/connector/suc v0.0.0-00010101000000-000000000000
	github.com/mattn/go-sqlite3 v1.11.0
	github.com/prometheus/client_golang v1.4.0
	github.com/russellhaering/goxmldsig v0.0.0-20180430223755-7acd5e4a6ef7
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.5
	github.com/stretchr/testify v1.4.0
	github.com/testcontainers/testcontainers-go v0.0.9
	go.etcd.io/etcd v0.0.0-20191023171146-3cf2f69b5738
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	golang.org/x/net v0.0.0-20200513185701-a91f0712d120
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	google.golang.org/api v0.15.0
	google.golang.org/grpc v1.26.0
	gopkg.in/ldap.v2 v2.5.1
	gopkg.in/square/go-jose.v2 v2.4.1
	sigs.k8s.io/testing_frameworks v0.1.2
)

replace github.com/dexidp/dex/api/v2 => ./api/v2

replace github.com/masami10/dex/connector/suc => ./connector/suc
