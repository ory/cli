module github.com/ory/cli

go 1.16

// Required because github.com/ory/kratos rewrites github.com/ory/kratos-client-go to
// github.com/ory/kratos/internal/httpclient
replace github.com/ory/kratos-client-go => github.com/ory/kratos-client-go v0.6.3-alpha.1

replace github.com/ory/x => github.com/ory/x v0.0.254

replace github.com/gobuffalo/pop/v5 => github.com/gobuffalo/pop/v5 v5.3.4-0.20210608105745-bb07a373cc0e

replace github.com/mattn/go-sqlite3 => github.com/mattn/go-sqlite3 v1.14.7-0.20210414154423-1157a4212dcb

replace github.com/ory/kratos => github.com/ory/kratos v0.6.3-alpha.1.0.20210608145203-b5c1658e01ca

require (
	github.com/DataDog/datadog-go v4.7.0+incompatible // indirect
	github.com/HdrHistogram/hdrhistogram-go v1.1.0 // indirect
	github.com/Masterminds/semver/v3 v3.1.1
	github.com/avast/retry-go v3.0.0+incompatible
	github.com/containerd/containerd v1.5.2 // indirect
	github.com/deckarep/golang-set v1.7.1
	github.com/elnormous/contenttype v0.0.0-20210110050721-79150725153f
	github.com/evanphx/json-patch v4.9.0+incompatible
	github.com/getkin/kin-openapi v0.48.0
	github.com/ghodss/yaml v1.0.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gobuffalo/fizz v1.13.1-0.20201104174146-3416f0e6618f
	github.com/gobuffalo/pop/v5 v5.3.3
	github.com/gofrs/uuid/v3 v3.1.2
	github.com/gomarkdown/markdown v0.0.0-20201113031856-722100d81a8e
	github.com/google/uuid v1.2.0
	github.com/hashicorp/go-retryablehttp v0.6.8
	github.com/jackc/pgx/v4 v4.11.0
	github.com/julienschmidt/httprouter v1.3.0
	github.com/markbates/pkger v0.17.1
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/ory/gochimp3 v0.0.0-20200417124117-ccd242db3655
	github.com/ory/graceful v0.1.1
	github.com/ory/herodot v0.9.6
	github.com/ory/jsonschema/v3 v3.0.3
	github.com/ory/kratos v0.6.3-alpha.1.0.20210608145540-005c8b8e28f3
	github.com/ory/kratos-client-go v0.6.3-alpha.1
	github.com/ory/x v0.0.250
	github.com/phayes/freeport v0.0.0-20180830031419-95f893ade6f2
	github.com/philhofer/fwd v1.1.1 // indirect
	github.com/pkg/errors v0.9.1
	github.com/smallstep/truststore v0.9.6
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/square/go-jose/v3 v3.0.0-20200630053402-0a67ce9b0693
	github.com/stretchr/testify v1.7.0
	github.com/tidwall/gjson v1.7.1
	github.com/tidwall/sjson v1.1.5
	github.com/txn2/txeh v1.3.0
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	github.com/urfave/negroni v1.0.0
	github.com/ziutek/mymysql v1.5.4 // indirect
	gopkg.in/gorp.v1 v1.7.2 // indirect
	gopkg.in/yaml.v2 v2.4.0
)
