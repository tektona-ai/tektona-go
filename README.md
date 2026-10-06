# Tektona Go SDK

The client exposes sandbox, location, metadata, organization, and project methods.
Use `client.Raw()` for public HTTP operations without resource methods.

```go
client, err := tektona.NewClient(
    tektona.WithOrg("acme"),
    tektona.WithProject("web"),
)
if err != nil { return err }

sandbox, err := client.Sandbox().Create(ctx, tektona.CreateSandboxParams{
    Template: "go-dev:stable",
})
if err != nil { return err }
fmt.Println(sandbox.Id)
```

Set `TEKTONA_API_KEY` before you construct the client, or use `WithAPIKey`.
Set optional `TEKTONA_API_URL`, `TEKTONA_ORG`, and `TEKTONA_PROJECT` defaults as needed.
An option overrides its environment value. A method parameter overrides the scope default.
The default API URL is `https://api.tektona.ai`.

Pass `&tektona.ListSandboxesParams{Project: &empty}` to list every project in an organization.
Pass a sandbox ID to `Get` or `Delete` without scope, or pass a name with org and project scope.
Use `client.Project().List(ctx, nil)` to list projects across organizations.
Use `client.Project().ListForOrg(ctx, "", nil)` to list projects in the default organization.
Pass a project name to `Get`, `Update`, `Delete`, and lifecycle default methods. An empty name returns an error.
Set `Org` in project method parameters to override the default organization.
Location, metadata, and cross-organization list methods do not use scope defaults.
API failures return `*tektona.APIError`. Use `errors.As` to inspect the status and problem.

The `api` package contains generated methods and types. The public OpenAPI snapshot and its source commit are in `api/`.
Run `go generate ./api` after you update the snapshot.
