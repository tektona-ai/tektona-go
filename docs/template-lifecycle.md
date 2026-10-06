# Template lifecycle

Use `client.TemplateLifecycle()` to read, update, and preview template lifecycle settings.
Organization methods use the `org` argument or the client default when `org` is empty.
Project methods require a project name and use the client organization default unless the parameters set `Org`.
An empty project name or a missing organization returns a local error before any request.

```go
client, err := tektona.NewClient(tektona.WithOrg("acme"))
if err != nil {
    return err
}

archiveDays := int32(30)
deleteDays := int32(0)
proposal := tektona.PreviewTemplateLifecycleForProjectParams{
    ArchiveAfterDaysUnused: &archiveDays,
    DeleteAfterDaysArchived: &deleteDays,
}
preview, err := client.TemplateLifecycle().PreviewForProject(ctx, "web", proposal)
if err != nil {
    return err
}
fmt.Println(preview.Archive.Versions, preview.Delete.Versions)
```

`GetForOrg` and `GetForProject` return `*api.TemplateLifecycleSettingsResponse`.
`UpdateForOrg` and `UpdateForProject` return the same response type.
`PreviewForOrg` and `PreviewForProject` return `*api.TemplateLifecyclePreviewResponse` without a saved change.
Both settings fields accept pointers: use `nil` to inherit the setting and `&zero` to send zero.
The server checks valid values. The client checks only the required scope.
HTTP failures return `*tektona.APIError`; use `errors.As` to inspect the status and response.
