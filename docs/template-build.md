# Template builds

Use `Client.TemplateBuild()` to create builds, list builds, read logs, and cancel builds.
Set an organization and project on the client, or pass them to each scoped call.
Organization builds need an organization but not a project.
Build IDs select unscoped methods such as `Get`, `Cancel`, `GetLogs`, and `StreamLogs`.

```go
client, err := tektona.NewClient(tektona.WithOrg("acme"), tektona.WithProject("web"))
if err != nil {
    return err
}
build, err := client.TemplateBuild().CreateForProject(ctx, tektona.CreateTemplateBuildParams{
    Template: "go-dev",
    Build: api.TemplateBuildSpec{Image: "ubuntu:24.04"},
})
if err != nil {
    return err
}
_, err = client.TemplateBuild().GetLogs(ctx, build.Id, nil)
```

`CreateForOrg` and `CreateForProject` accept `IdempotencyKey` for safe retries.
List and log methods accept generated filters and cursor fields through root-package parameter aliases.
An HTTP error returns `*tektona.APIError`; use `errors.As` to inspect its status, headers, and problem body.

`StreamLogs` returns raw server-sent event (SSE) frames as `io.ReadCloser`.
Close the stream when you stop reading it.
Cancel the call context to stop a blocked read.
The stream does not decode SSE frames or reconnect automatically.

```go
stream, err := client.TemplateBuild().StreamLogs(ctx, build.Id, nil)
if err != nil {
    return err
}
defer stream.Close()
_, err = io.Copy(os.Stdout, stream)
```
