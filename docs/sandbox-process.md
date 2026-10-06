# Sandbox processes

Use `client.SandboxProcess()` to manage processes in a sandbox.

Pass the sandbox ID or name to each method. Pass a process ID or name to methods that target one process.

For a sandbox name, set `WithOrg` and `WithProject`, or pass `Org` and `Project` in the call parameters.
A sandbox ID needs no scope. Each call can override the client defaults.

```go
client, err := tektona.NewClient(
    tektona.WithAPIKey(key),
    tektona.WithOrg("acme"),
    tektona.WithProject("web"),
)
if err != nil { return err }

process, err := client.SandboxProcess().Start(ctx, "my-sandbox", tektona.StartSandboxProcessOptions{
    Command: "go test ./...",
})
if err != nil { return err }

logs, err := client.SandboxProcess().GetLogs(ctx, "my-sandbox", process.Id, nil)
```

`List`, `Stop`, `Get`, and `GetLogs` accept parameter aliases from the generated API.
Other methods take one options struct for scope, query, and body fields.

| Method | Result |
| --- | --- |
| `List` | Process list |
| `Start` | Process after HTTP 201 |
| `Stop` | Process after HTTP 202 |
| `Get`, `Rename` | Process |
| `GetLogs` | Buffered log frames |
| `SetAutostart`, `Resize`, `Signal`, `WriteStdin` | No body after HTTP 204 |
| `CreateStreamAccess` | Token and WebSocket URL |

`CreateStreamAccess` returns a short-lived access token and a WebSocket URL. It does not open a stream. Treat the token and URL as credentials.

HTTP errors return `*tektona.APIError`. Inspect the status and headers with `errors.As`. Transport errors and invalid responses return regular Go errors.
