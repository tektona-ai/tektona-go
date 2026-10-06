# Sandbox actions

Use `client.Sandbox()` to call the ten sandbox actions below.
Pass a sandbox ID or name as the second argument.
For a name, set both `org` and `project` on the client or the call parameters.
An ID does not require either scope value.

| Method | Parameters | Response |
| --- | --- | --- |
| `StartDesktop` | `*StartDesktopParams` | `*api.DesktopActionOutputBody` |
| `StopDesktop` | `*StopDesktopParams` | `*api.DesktopActionOutputBody` |
| `Pause` | `PauseSandboxParams` | `*api.PauseSandboxResponse` |
| `Reboot` | `*RebootSandboxParams` | `*api.RebootSandboxResponseBody` |
| `Reset` | `*ResetSandboxParams` | `*api.RebootSandboxResponseBody` |
| `Resize` | `ResizeSandboxParams` | `*api.ResizeResponse` |
| `Resume` | `*ResumeSandboxParams` | `*api.ResumeSandboxResponseBody` |
| `Fork` | `ForkSandboxParams` | `*api.ForkResponse` |
| `Rename` | `RenameSandboxParams` | `*api.RenameSandboxResponse` |
| `Transfer` | `TransferSandboxParams` | `*api.SandboxTransferResponse` |

Each method takes `context.Context` first.
The five pointer parameter types hold optional scope only; pass `nil` to use client scope.
The value parameter types combine body fields with optional scope pointers.
For `Rename`, a nil `Name` removes the name.
For `Transfer`, set `To` to the new owner email or user ID.
The API checks other field values.

```go
result, err := client.Sandbox().Resize(ctx, "api-dev", tektona.ResizeSandboxParams{
    Cpu: &cpu,
})
if err != nil {
    var apiErr *tektona.APIError
    if errors.As(err, &apiErr) {
        fmt.Println(apiErr.StatusCode)
    }
    return err
}
fmt.Println(result.AppliesImmediately)
```

These methods accept only HTTP 200 as success.
On another status, inspect `*tektona.APIError` for the status, headers, body, and parsed problem.
