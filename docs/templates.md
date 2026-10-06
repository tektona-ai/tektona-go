# Templates

Use `Client.Template()` to access organization, project, and system templates.
Set `TEKTONA_API_KEY` before you create a client, or pass `WithAPIKey`.
Set the default scope with `WithOrg` and `WithProject`, or use `TEKTONA_ORG` and `TEKTONA_PROJECT`.
Pass an organization or project in a call to override its default.

```go
client, err := tektona.NewClient(tektona.WithOrg("acme"), tektona.WithProject("web"))
if err != nil {
    return err
}

created, err := client.Template().CreateForProject(ctx, tektona.CreateTemplateForProjectParams{
    Name: "go-dev",
})
if err != nil {
    return err
}

template, err := client.Template().GetForProject(ctx, created.Metadata.Name, nil)
if err != nil {
    return err
}
fmt.Println(template.Reference)
```

`ListForOrg` and `ListForProject` accept optional `Q`, `Cursor`, and `Limit` parameters.
Pass `ListTemplatesBody.Pagination.NextCursor` as `Cursor` to read another page.
`CreateForOrg` and `CreateForProject` accept a name, display name, and description.
`UpdateForOrg` and `UpdateForProject` use the path name in `metadata.name`.
`DeleteForOrg` and `DeleteForProject` return `DeleteTemplateResponse` with the number of deleted builds.
`ActivateForOrg`, `ArchiveForOrg`, `ActivateForProject`, and `ArchiveForProject` return the updated template.
`GetSystem` reads a system template by name and does not use the client scope.

If a required name or scope is empty, the call returns a local error before it sends a request.
If the API rejects a request, inspect `*tektona.APIError` with `errors.As`.
