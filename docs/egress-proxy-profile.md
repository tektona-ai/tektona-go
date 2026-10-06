# Egress proxy profiles

Use `client.EgressProxyProfile()` for organization and project profiles. The accessor uses the same HTTP client as `client.Raw()`.

Set the API key and the default scope when you create the client:

```go
client, err := tektona.NewClient(
    tektona.WithAPIKey(key),
    tektona.WithOrg("acme"),
    tektona.WithProject("web"),
)
if err != nil {
    return err
}

profiles, err := client.EgressProxyProfile().ListForProject(ctx, "", nil)
if err != nil {
    return err
}
_ = profiles.Items
```

An empty scope override uses the client default. A call override takes precedence over the client default.
The client reads `TEKTONA_ORG` and `TEKTONA_PROJECT` when no explicit defaults exist.
Organization methods do not use the project default.

The nine methods are:

| Scope | Methods |
| --- | --- |
| Organization | `ListForOrg`, `CreateForOrg` |
| Project | `ListForProject`, `CreateForProject`, `AddRuleForProject`, `DeleteForProject`, `SetDefaultForProject`, `DeleteRuleForProject`, `UpdateRuleForProject` |

Pass a profile name to `AddRuleForProject`. Pass a profile ID to delete a profile or change its default flag.
Pass both the profile ID and rule ID to change or delete a rule.
`CreateForOrg` sends `scope: "org"`; `CreateForProject` sends `scope: "project"`.

The list methods return `*api.ProfileListBody`. The create methods return `*api.ProfileCreateBody`.
`AddRuleForProject` returns `*api.RuleCreatedBody`. The four methods with no response body return only an error.
Use the `Cursor` and `Limit` fields in the list parameters to request a page.

For an HTTP error, use `errors.As` to get `*tektona.APIError`. It contains the status, headers, body, and parsed problem.
The client does not retry a request. The caller controls timeouts through the context or HTTP client.
