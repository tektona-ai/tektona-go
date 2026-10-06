# Template versions

Use `client.TemplateVersion()` to list, get, delete, prune, activate, or archive template versions.
Organization and project routes use the client scope defaults.
Set `Org` or `Project` in a parameter struct to override a default for one call.
System routes ignore these defaults.

```go
limit := int32(10)
versions, err := client.TemplateVersion().ListForProject(ctx, "go-dev", &tektona.ListTemplateVersionsForProjectParams{
    Limit: &limit,
})
// Check err before use.
version, err := client.TemplateVersion().GetForOrg(ctx, "go-dev", "01ABC", nil)
// Check err before use.
dryRun := true
report, err := client.TemplateVersion().PruneForOrg(ctx, "go-dev", tektona.PruneTemplateVersionsForOrgParams{
    DryRun: &dryRun,
})
```

`ListForOrg`, `ListForProject`, and `ListSystem` accept state, tag, text, cursor, and limit filters.
Pass `DryRun` to preview a prune without deletion.
Delete succeeds only on HTTP 204.
Other calls return generated `api` response bodies.
HTTP failures return `*tektona.APIError`; use `errors.As` to inspect the status and problem.
