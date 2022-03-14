module github.com/ory/cli

go 1.16

// Required because github.com/ory/kratos rewrites github.com/ory/kratos-client-go to
// github.com/ory/kratos/internal/httpclient
replace github.com/ory/kratos-client-go => github.com/ory/kratos-client-go v0.7.6-alpha.7.0.20211020080137-4b204b1a9f2d

replace github.com/gobuffalo/pop/v5 => github.com/gobuffalo/pop/v5 v5.3.4-0.20210608105745-bb07a373cc0e

replace github.com/mattn/go-sqlite3 => github.com/mattn/go-sqlite3 v1.14.9

replace github.com/ory/kratos => github.com/ory/kratos v0.7.6-alpha.1.0.20211023171302-21270a85f39c

replace github.com/ory/x => ../x

require (
	github.com/Masterminds/semver/v3 v3.1.1
	github.com/avast/retry-go v3.0.0+incompatible
	github.com/deckarep/golang-set v1.7.1
	github.com/evanphx/json-patch v4.11.0+incompatible
	github.com/getkin/kin-openapi v0.48.0
	github.com/ghodss/yaml v1.0.0
	github.com/go-sql-driver/mysql v1.6.0
	github.com/gobuffalo/fizz v1.14.0
	github.com/gobuffalo/pop/v5 v5.3.4
	github.com/gofrs/uuid/v3 v3.1.2
	github.com/gomarkdown/markdown v0.0.0-20201113031856-722100d81a8e
	github.com/google/uuid v1.3.0
	github.com/hashicorp/go-retryablehttp v0.7.0
	github.com/jackc/pgx/v4 v4.13.0
	github.com/markbates/pkger v0.17.1
	github.com/ory/gochimp3 v0.0.0-20200417124117-ccd242db3655
	github.com/ory/graceful v0.1.1
	github.com/ory/herodot v0.9.13
	github.com/ory/jsonschema/v3 v3.0.7
	github.com/ory/kratos-client-go v0.8.0-alpha.2
	github.com/ory/x v0.0.309
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	github.com/square/go-jose/v3 v3.0.0-20200630053402-0a67ce9b0693
	github.com/stretchr/testify v1.7.0
	github.com/tidwall/gjson v1.9.4
	github.com/tidwall/sjson v1.2.2
	github.com/urfave/negroni v1.0.0
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211
	gopkg.in/yaml.v2 v2.4.0
)
