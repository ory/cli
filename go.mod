module github.com/ory/cli

go 1.16

// Required because github.com/ory/kratos rewrites github.com/ory/kratos-client-go to
// github.com/ory/kratos/internal/httpclient
replace github.com/ory/kratos-client-go => github.com/ory/kratos-client-go v0.5.4-alpha.1.0.20210210170256-960b093d8bf9

replace github.com/ory/kratos/corp => github.com/ory/kratos/corp v0.0.0-20210118092700-c2358be1e867

require (
	github.com/Masterminds/semver/v3 v3.1.1
	github.com/avast/retry-go v3.0.0+incompatible
	github.com/caddyserver/caddy/v2 v2.3.0
	github.com/caddyserver/certmagic v0.12.1-0.20201215190346-201f83a06067
	github.com/deckarep/golang-set v1.7.1
	github.com/elnormous/contenttype v0.0.0-20210110050721-79150725153f
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gobuffalo/fizz v1.13.1-0.20201104174146-3416f0e6618f
	github.com/gobuffalo/pop/v5 v5.3.2-0.20210128124218-e397a61c1704
	github.com/gofrs/uuid/v3 v3.1.2
	github.com/gomarkdown/markdown v0.0.0-20201113031856-722100d81a8e
	github.com/google/uuid v1.1.5
	github.com/jackc/pgx/v4 v4.10.1
	github.com/julienschmidt/httprouter v1.3.0
	github.com/markbates/pkger v0.17.1
	github.com/ory/gochimp3 v0.0.0-20200417124117-ccd242db3655
	github.com/ory/graceful v0.1.1
	github.com/ory/herodot v0.9.0
	github.com/ory/jsonschema/v3 v3.0.2
	github.com/ory/kratos v0.5.5-alpha.1.0.20210225142555-c4c6ed9619fa
	github.com/ory/kratos-client-go v0.5.4-alpha.1
	github.com/ory/x v0.0.198
	github.com/phayes/freeport v0.0.0-20180830031419-95f893ade6f2
	github.com/pkg/errors v0.9.1
	github.com/smallstep/cli v0.15.2
	github.com/smallstep/truststore v0.9.6
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/square/go-jose/v3 v3.0.0-20200630053402-0a67ce9b0693
	github.com/stretchr/testify v1.7.0
	github.com/tidwall/gjson v1.6.7
	github.com/tidwall/sjson v1.1.4
	github.com/urfave/negroni v1.0.0
	gopkg.in/yaml.v2 v2.4.0
)
