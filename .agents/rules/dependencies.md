# Dependency safety rules

Before adding any new library, package, or dependency, you MUST perform a security and quality validation. If any check fails, do not install it, warn the user, and suggest alternatives.

## Checks

Run these checks **before** executing `go get` or adding to `go.mod`. Skip only if the user explicitly overrides.

### 1. Source reputation & authenticity

- What is the module path? Prefer well-known organizations and maintainers (e.g., `github.com/gin-gonic/gin`, `github.com/go-jose/go-jose/v3`, `golang.org/x/...`).
- Check the GitHub/GitLab repository:
  - Verified organization or individual maintainer (check badge, contributor history).
  - Not a fork with obscure modifications — verify against the official source.
  - The repository has real issues, PRs, and discussion activity (not an empty shell).
- **Typosquatting**: check if the name impersonates a popular library by visual similarity (e.g., `gihub.com/...`, `gormm.io/...`, `go-jose` vs `gopkg.in/jose`).

### 2. Community & maintenance health

- Last commit → less than 6 months ago (prefer < 3 months).
- Open issues → not overwhelmingly high relative to stars/contributors.
- Stars and contributors → a healthy project has community traction (not a single-contributor one-push repo).
- Release tags → at least one stable `v1.0.0+` or evidence of real semantic versioning.
- Is the project archived or unmaintained? Skip if yes.

### 3. Security record

- Check the repository's Security tab in GitHub. Any unaddressed advisories? Skip if critical vulnerabilities exist.
- No known CVEs or malicious reports (search: `[package name] CVE`, `[package name] malware`).
- No suspicious patterns:
  - Obfuscated code, minified blobs, encoded payloads.
  - `init()` functions that open network connections, read environment variables, or execute system commands without clear documentation.
  - Hidden data collection or telemetry without explicit opt-in.

### 4. License

- MUST have a clear OSI-compliant license (MIT, Apache 2.0, BSD 3-clause, etc.).
- No license? Warn and ask user before proceeding — unlicensed code has no grant of rights.
- Avoid GPL/AGPL unless the user confirms they are compatible with the project license.

### 5. Quality & API stability

- Prefer libraries that follow Go idioms: context-aware functions, proper error handling, well-documented exported API.
- Prefer standard library or the official extended library (`golang.org/x/...`) over third-party when both exist.
- If the library wraps something simple (e.g., an entire dependency for a single helper function), prefer implementing it inline or using a smaller alternative.
- The Go module proxy (`proxy.golang.org`) should have the module indexed — verify with `go env GOPROXY`.

### 6. Supply chain

- Check `go.mod` of the target dependency: does it pull in an excessive tree of transitive dependencies? Warn if a small feature drags in 50+ indirect packages.
- Avoid dependencies that have not been updated to use the project's Go version (check `go` directive in the dependency's `go.mod`).

## Procedure

```
1. Identify the best candidate library for the task
2. Run the checks above on the candidate
3. If it fails any check → search for a better alternative
4. If no safe alternative exists → inform the user with the risk details
5. Only then run go get <module>
6. After go get: run make check (fmt + vet) AND make build to verify compilation
```

## Red flags (automatic rejection)

- < 10 stars AND < 2 contributors AND no release tags.
- Last commit > 1 year ago.
- Open critical security advisory or known CVE without a fix.
- Contains obfuscated or intentionally unreadable code.
- Zero documentation (no README, no GoDoc comments, no examples).
- Impersonates a well-known package (typosquatting).
- Repository is archived, deleted, or returns 404.
- Uses `unsafe` or `syscall` without clear justification.
