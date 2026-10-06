# Egress network policies

Use `client.EgressNetworkPolicy()` to read system policies and manage organization or project policies.
System methods do not use organization or project defaults.

```go
client, err := tektona.NewClient(
    tektona.WithAPIKey(key),
    tektona.WithOrg("acme"),
    tektona.WithProject("web"),
)
if err != nil {
    return err
}

domains := []string{"github.com"}
policy, err := client.EgressNetworkPolicy().CreateForProject(ctx,
    tektona.CreateEgressNetworkPolicyForProjectParams{
        Name: "source", AllowedDomains: &domains,
    })
if err != nil {
    return err
}
fmt.Println(policy.Name)
```

Import `github.com/tektona-ai/tektona-go` as `tektona` for this example.
Import `fmt` and define `ctx` and `key` before the call.

`ListSystem` returns `*api.ListSystemBody`; `GetSystem` returns `*api.Policy`.
Organization and project lists return `*api.ListPoliciesBody`.
Their create, get, and update methods return `*api.PolicyResponse`.
Delete methods return only an error and require HTTP 204.

Pass `Org` and `Project` in a method's parameters to override client defaults.
An empty override inherits the matching default.
The client reads `TEKTONA_ORG` and `TEKTONA_PROJECT` when no explicit client option exists.
Use `WithOrg("")` or `WithProject("")` to clear an environment default.
Organization methods require an organization; project methods require both scopes.
Get, update, and delete methods require a policy name.
The client rejects missing scopes and names before it sends a request.

Pass `api` model fields through create and update parameter structs.
The server validates policy rules and names.
Inspect HTTP errors with `errors.As(err, &apiErr)` where `apiErr` has type `*tektona.APIError`.
Its `StatusCode`, `Header`, `Body`, and `Problem` fields contain the response details.
