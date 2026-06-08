module github.com/ory/cli

go 1.26

replace github.com/gorilla/sessions => github.com/ory/sessions v1.2.2-0.20220110165800-b09c17334dc2

require (
	dario.cat/mergo v1.0.2
	github.com/Masterminds/semver/v3 v3.5.0
	github.com/evanphx/json-patch v5.9.0+incompatible
	github.com/getkin/kin-openapi v0.131.0
	github.com/ghodss/yaml v1.0.0
	github.com/go-jose/go-jose/v3 v3.0.5
	github.com/gofrs/uuid v4.4.0+incompatible
	github.com/gomarkdown/markdown v0.0.0-20240730141124-034f12af3bf6
	github.com/hashicorp/go-retryablehttp v0.7.8
	github.com/ory/client-go v1.22.41
	github.com/ory/gochimp3 v0.0.0-20200417124117-ccd242db3655
	github.com/ory/graceful v0.2.0
	github.com/ory/herodot v0.10.9-0.20260330111132-da75ef0fbc22
	github.com/ory/hydra-client-go/v2 v2.4.0-alpha.1.0.20251107123905-f3d35665821b
	github.com/ory/hydra/v2 v2.3.1-0.20260521102426-2e62d4d951ef
	github.com/ory/jsonschema/v3 v3.0.9-0.20250317235931-280c5fc7bf0e
	github.com/ory/keto v0.14.1-0.20260521082736-bcf74421eaf2
	github.com/ory/kratos v1.3.1-0.20260521101516-8b2f33bd5c3d
	github.com/ory/x v0.0.730-0.20260521082735-5aba9539f9e9
	github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c
	github.com/pkg/errors v0.9.1
	github.com/playwright-community/playwright-go v0.4702.0
	github.com/rs/cors v1.11.1
	github.com/sabhiram/go-gitignore v0.0.0-20210923224102-525f6e181f06
	github.com/spf13/cobra v1.10.2
	github.com/spf13/pflag v1.0.10
	github.com/stretchr/testify v1.11.1
	github.com/tidwall/gjson v1.19.0
	github.com/tidwall/sjson v1.2.5
	github.com/urfave/negroni v1.0.0
	golang.org/x/oauth2 v0.36.0
	golang.org/x/text v0.37.0
	gopkg.in/yaml.v2 v2.4.0
)

require (
	cloud.google.com/go/compute/metadata v0.9.0 // indirect
	code.dny.dev/ssrf v0.2.0 // indirect
	filippo.io/edwards25519 v1.2.0 // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/sprig/v3 v3.3.0 // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/ProtonMail/go-crypto v0.0.0-20230717121422-5aa5874ade95 // indirect
	github.com/ProtonMail/go-mime v0.0.0-20230322103455-7d82a3887f2f // indirect
	github.com/ProtonMail/gopenpgp/v2 v2.7.5 // indirect
	github.com/ThalesGroup/crypto11 v1.4.1 // indirect
	github.com/XSAM/otelsql v0.42.0 // indirect
	github.com/arbovm/levenshtein v0.0.0-20160628152529-48b4e1c0c4d0 // indirect
	github.com/asaskevich/govalidator v0.0.0-20230301143203-a9d515a09cc2 // indirect
	github.com/avast/retry-go/v4 v4.7.0 // indirect
	github.com/aymerick/douceur v0.2.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/boombuler/barcode v1.1.0 // indirect
	github.com/bradleyjkemp/cupaloy/v2 v2.8.0 // indirect
	github.com/bwmarrin/discordgo v0.29.0 // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/cenkalti/backoff/v3 v3.2.2 // indirect
	github.com/cenkalti/backoff/v5 v5.0.3 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cloudflare/circl v1.6.3 // indirect
	github.com/cockroachdb/cockroach-go/v2 v2.4.3 // indirect
	github.com/containerd/errdefs v1.0.0 // indirect
	github.com/containerd/errdefs/pkg v0.3.0 // indirect
	github.com/coreos/go-oidc/v3 v3.18.0 // indirect
	github.com/cristalhq/jwt/v4 v4.0.2 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/deckarep/golang-set/v2 v2.6.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.4.1 // indirect
	github.com/dghubble/oauth1 v0.7.3 // indirect
	github.com/dgraph-io/ristretto/v2 v2.4.0 // indirect
	github.com/distribution/reference v0.6.0 // indirect
	github.com/docker/docker v28.5.2+incompatible // indirect
	github.com/docker/go-connections v0.7.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/evanphx/json-patch/v5 v5.9.11 // indirect
	github.com/fatih/color v1.19.0 // indirect
	github.com/fatih/structs v1.1.0 // indirect
	github.com/felixge/fgprof v0.9.5 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/fsnotify/fsnotify v1.10.1 // indirect
	github.com/fxamacker/cbor/v2 v2.9.2 // indirect
	github.com/go-crypt/crypt v0.4.14 // indirect
	github.com/go-crypt/x v0.4.15 // indirect
	github.com/go-faker/faker/v4 v4.7.0 // indirect
	github.com/go-jose/go-jose/v4 v4.1.4 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-openapi/errors v0.22.7 // indirect
	github.com/go-openapi/jsonpointer v0.23.1 // indirect
	github.com/go-openapi/strfmt v0.26.2 // indirect
	github.com/go-openapi/swag/jsonname v0.26.0 // indirect
	github.com/go-sql-driver/mysql v1.10.0 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/go-viper/mapstructure/v2 v2.5.0 // indirect
	github.com/go-webauthn/webauthn v0.17.3 // indirect
	github.com/go-webauthn/x v0.2.5 // indirect
	github.com/gobuffalo/envy v1.10.2 // indirect
	github.com/gobuffalo/fizz v1.14.4 // indirect
	github.com/gobuffalo/flect v1.0.3 // indirect
	github.com/gobuffalo/github_flavored_markdown v1.1.4 // indirect
	github.com/gobuffalo/helpers v0.6.10 // indirect
	github.com/gobuffalo/nulls v0.4.2 // indirect
	github.com/gobuffalo/plush/v4 v4.1.22 // indirect
	github.com/gobuffalo/plush/v5 v5.0.11 // indirect
	github.com/gobuffalo/tags/v3 v3.1.4 // indirect
	github.com/gobuffalo/validate/v3 v3.3.3 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/goccy/go-json v0.10.6 // indirect
	github.com/goccy/go-yaml v1.19.2 // indirect
	github.com/gofrs/flock v0.13.0 // indirect
	github.com/gogo/googleapis v1.4.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.2 // indirect
	github.com/golang-jwt/jwt/v5 v5.3.1 // indirect
	github.com/golang/gddo v0.0.0-20210115222349-20d68f94ee1f // indirect
	github.com/google/go-github/v38 v38.1.0 // indirect
	github.com/google/go-jsonnet v0.22.0 // indirect
	github.com/google/go-querystring v1.2.0 // indirect
	github.com/google/go-tpm v0.9.8 // indirect
	github.com/google/pprof v0.0.0-20260507013755-92041b743c96 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gopherjs/gopherjs v1.20.1 // indirect
	github.com/gorilla/context v1.1.2 // indirect
	github.com/gorilla/css v1.0.1 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/gorilla/pat v1.0.2 // indirect
	github.com/gorilla/securecookie v1.1.2 // indirect
	github.com/gorilla/sessions v1.4.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.3.3 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.29.0 // indirect
	github.com/gtank/cryptopasta v0.0.0-20170601214702-1f550f6f2f69 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/huandu/xstrings v1.5.0 // indirect
	github.com/ian-kent/envconf v0.0.0-20141026121121-c19809918c02 // indirect
	github.com/ian-kent/go-log v0.0.0-20160113211217-5731446c36ab // indirect
	github.com/ian-kent/goose v0.0.0-20141221090059-c3541ea826ad // indirect
	github.com/ian-kent/linkio v0.0.0-20170807205755-97566b872887 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/inhies/go-bytesize v0.0.0-20220417184213-4913239db9cf // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.9.2 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jaegertracing/jaeger-idl v0.7.1 // indirect
	github.com/jmoiron/sqlx v1.4.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51 // indirect
	github.com/knadh/koanf/maps v0.1.2 // indirect
	github.com/knadh/koanf/parsers/json v1.0.0 // indirect
	github.com/knadh/koanf/parsers/toml v0.1.0 // indirect
	github.com/knadh/koanf/parsers/yaml v1.1.0 // indirect
	github.com/knadh/koanf/providers/posflag v1.0.1 // indirect
	github.com/knadh/koanf/v2 v2.3.4 // indirect
	github.com/landlock-lsm/go-landlock v0.8.1 // indirect
	github.com/lestrrat-go/backoff/v2 v2.0.8 // indirect
	github.com/lestrrat-go/blackmagic v1.0.4 // indirect
	github.com/lestrrat-go/httpcc v1.0.1 // indirect
	github.com/lestrrat-go/iter v1.0.2 // indirect
	github.com/lestrrat-go/jwx v1.2.31 // indirect
	github.com/lestrrat-go/option v1.0.1 // indirect
	github.com/lib/pq v1.12.3 // indirect
	github.com/magiconair/properties v1.8.9 // indirect
	github.com/mailhog/MailHog-Server v1.0.1 // indirect
	github.com/mailhog/data v1.0.1 // indirect
	github.com/mailhog/smtp v1.0.1 // indirect
	github.com/mailhog/storage v1.0.1 // indirect
	github.com/mailru/easyjson v0.9.2 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.22 // indirect
	github.com/mattn/goveralls v0.0.12 // indirect
	github.com/microcosm-cc/bluemonday v1.0.27 // indirect
	github.com/miekg/pkcs11 v1.1.1 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/moby/docker-image-spec v1.3.1 // indirect
	github.com/moby/moby/api v1.54.2 // indirect
	github.com/moby/moby/client v0.4.1 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/montanaflynn/stats v0.9.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/ncruces/go-strftime v1.0.0 // indirect
	github.com/nyaruka/phonenumbers v1.7.4 // indirect
	github.com/oasdiff/yaml v0.0.0-20250309154309-f31be36b4037 // indirect
	github.com/oasdiff/yaml3 v0.0.0-20250309153720-d2182401db90 // indirect
	github.com/oklog/ulid/v2 v2.1.1 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.1 // indirect
	github.com/openzipkin/zipkin-go v0.4.3 // indirect
	github.com/ory/analytics-go/v5 v5.0.1 // indirect
	github.com/ory/dockertest/v4 v4.0.0 // indirect
	github.com/ory/go-acc v0.2.9-0.20230103102148-6b1c9a70dbbe // indirect
	github.com/ory/go-convenience v0.1.1-0.20251022160015-e2a1f648d0b1 // indirect
	github.com/ory/keto/proto v0.13.0-alpha.0.0.20260521082736-bcf74421eaf2 // indirect
	github.com/ory/kratos-client-go v1.3.9-0.20251107123727-a6ddbd382e38 // indirect
	github.com/ory/mail/v3 v3.0.1-0.20260416102637-e762925f059d // indirect
	github.com/ory/nosurf v1.2.7 // indirect
	github.com/ory/pop/v6 v6.4.2-0.20260507161217-89126558d369 // indirect
	github.com/pborman/uuid v1.2.1 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/perimeterx/marshmallow v1.1.5 // indirect
	github.com/peterhellberg/link v1.2.0 // indirect
	github.com/philhofer/fwd v1.2.0 // indirect
	github.com/pkg/profile v1.7.0 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/pquerna/otp v1.5.0 // indirect
	github.com/prometheus/client_golang v1.23.2 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.67.5 // indirect
	github.com/prometheus/procfs v0.20.1 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/sagikazarmark/locafero v0.4.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/samber/lo v1.53.0 // indirect
	github.com/sawadashota/encrypta v0.0.5 // indirect
	github.com/seatgeek/logrus-gelf-formatter v0.0.0-20210414080842-5b05eb8ff761 // indirect
	github.com/segmentio/backo-go v1.1.0 // indirect
	github.com/sergi/go-diff v1.4.0 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	github.com/sirupsen/logrus v1.9.4 // indirect
	github.com/slack-go/slack v0.23.1 // indirect
	github.com/smarty/assertions v1.16.0 // indirect
	github.com/soheilhy/cmux v0.1.5 // indirect
	github.com/sourcegraph/annotate v0.0.0-20160123013949-f4cad6c6324d // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/sourcegraph/syntaxhighlight v0.0.0-20170531221838-bd320f5d308e // indirect
	github.com/spf13/afero v1.11.0 // indirect
	github.com/spf13/cast v1.10.0 // indirect
	github.com/spf13/viper v1.18.2 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/t-k/fluent-logger-golang v1.0.0 // indirect
	github.com/thales-e-security/pool v0.0.2 // indirect
	github.com/tidwall/match v1.2.0 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/tinylib/msgp v1.6.4 // indirect
	github.com/tomnomnom/linkheader v0.0.0-20180905144013-02ca5825eb80 // indirect
	github.com/toqueteos/webbrowser v1.2.0 // indirect
	github.com/twmb/murmur3 v1.1.8 // indirect
	github.com/ugorji/go v1.1.4 // indirect
	github.com/wI2L/jsondiff v0.7.1 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	github.com/xtgo/uuid v0.0.0-20140804021211-a0b114877d4c // indirect
	github.com/zmb3/spotify/v2 v2.4.3 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.68.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace v0.68.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.68.0 // indirect
	go.opentelemetry.io/contrib/propagators/b3 v1.43.0 // indirect
	go.opentelemetry.io/contrib/propagators/jaeger v1.43.0 // indirect
	go.opentelemetry.io/contrib/samplers/jaegerremote v0.37.0 // indirect
	go.opentelemetry.io/otel v1.43.0 // indirect
	go.opentelemetry.io/otel/exporters/jaeger v1.17.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.43.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.43.0 // indirect
	go.opentelemetry.io/otel/exporters/zipkin v1.43.0 // indirect
	go.opentelemetry.io/otel/metric v1.43.0 // indirect
	go.opentelemetry.io/otel/sdk v1.43.0 // indirect
	go.opentelemetry.io/otel/trace v1.43.0 // indirect
	go.opentelemetry.io/proto/otlp v1.10.0 // indirect
	go.uber.org/automaxprocs v1.6.0 // indirect
	go.uber.org/mock v0.6.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.yaml.in/yaml/v2 v2.4.4 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/crypto v0.51.0 // indirect
	golang.org/x/exp v0.0.0-20260508232706-74f9aab9d74a // indirect
	golang.org/x/mod v0.36.0 // indirect
	golang.org/x/net v0.54.0 // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/sys v0.44.0 // indirect
	golang.org/x/telemetry v0.0.0-20260508192327-42602be52be6 // indirect
	golang.org/x/tools v0.45.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260519071638-aa98bba5eb94 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260519071638-aa98bba5eb94 // indirect
	google.golang.org/grpc v1.81.1 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	kernel.org/pub/linux/libs/security/libcap/psx v1.2.78 // indirect
	modernc.org/libc v1.72.3 // indirect
	modernc.org/mathutil v1.7.1 // indirect
	modernc.org/memory v1.11.0 // indirect
	modernc.org/sqlite v1.50.1 // indirect
	sigs.k8s.io/yaml v1.6.0 // indirect
)

tool golang.org/x/tools/cmd/goimports
