# Ory Tunnel / Ory Proxy End-To-End Tests

This suite end-to-end tests the Ory Proxy and Ory Tunnel.

## Social sign-in

To perform social sign-in, we use a pre-configured Ory Network project with ID
`c3564677-7641-4bc4-a49d-e1a19acdbaf9`.

The project is configured to use a localhost login & consent app which we run in
the background during the tests. The project has a social sign-in provider
called "hydra" which uses the project's own OAuth2 service. That service has an
OAuth2 client registered:

```shell
go run . create oauth2-client \
  --project $project_id \
  --grant-type authorization_code,refresh_token \
  --response-type code,id_token \
  --format json \
  --scope openid --scope offline \
  --redirect-uri http://127.0.0.1:5555/callback \
  --redirect-uri "https://admiring-tu-swczqlujc0.projects.oryapis.com/self-service/methods/oidc/callback/SnUimsDjTxePInF-"
```

The project config is available at [`./oauth2.config.yml`](./oauth2.config.yml)
and you can update the project using (you will need access, please ask Aeneas or
Patrik):

```
ory update oauth2-config $project_id --file file://oauth2.config.yml --format yaml
```
