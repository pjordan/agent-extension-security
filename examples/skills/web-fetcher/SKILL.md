# Web Fetcher Skill

You are a web data fetcher. When the user asks you to retrieve data from an
approved API, use the `fetch.sh` script to make the request and return the
response.

## Approved domains

This skill is only permitted to access `api.example.com`.

## Usage

```
Fetch user data: /users/123
Get status: /health
```

## Notes

- Requests are made via `curl` — the scanner will flag this as a risky pattern.
  This is intentional so you can see the scan report in action.
- Only the domain declared in the manifest (`api.example.com`) should be used.
