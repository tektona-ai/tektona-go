# Git credentials, registries, and repositories

Use `GitCredential()`, `Registry()`, and `Repository()` to manage resources in a project.
Set the default organization and project on the client, or pass `Org` and `Project` with each call.
A nonempty call value takes precedence over the client default.
The client reads `TEKTONA_ORG` and `TEKTONA_PROJECT` when you do not set explicit defaults.
Calls that need a project return a local error if either scope value is empty.

```go
client, err := tektona.NewClient(
    tektona.WithAPIKey(apiKey),
    tektona.WithOrg("acme"),
    tektona.WithProject("web"),
)
if err != nil {
    return err
}

repos, err := client.Repository().List(ctx, nil)
if err != nil {
    return err
}
_ = repos.Items
```

`GitCredential().List` accepts a `ListGitCredentialsParams` value with an optional credential scope.
`GitCredential().Create` accepts the token, forge, credential scope, and repository IDs in `CreateGitCredentialParams`.
`GitCredential().Update` and `Delete` require a credential ID and a credential scope.
The token stays in the request body; error responses can include server details.

`Registry().List` accepts `Cursor` and `Limit` for pagination.
`Registry().Create` and `Update` accept `DryRun` with the registry fields in one parameter value.
The registry response does not return the password or token.

`Repository().List` accepts a `Default` filter.
`Repository().Create` and `Update` accept the name, URL, default branch, and default-set flag.
`Repository().Update`, `Delete`, and `ListBranches` require a UUID string as the repository ID.
`Repository().ListBranches` accepts `Refresh`, `Q`, `Cursor`, and `Limit`.

Every method returns a generated `api` response type or an error.
An HTTP error returns `*tektona.APIError`, which exposes its status, headers, body, and parsed problem.
Delete methods return nil only for HTTP 204.
Use `client.Raw()` when you need a generated method without client scope defaults.
