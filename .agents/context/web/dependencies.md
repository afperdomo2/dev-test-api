# Frontend dependency safety rules

Before adding any new npm package, you MUST perform a security and quality validation. If any check fails, do not install it, warn the user, and suggest alternatives.

## Checks

Run these checks **before** executing `pnpm add`. Skip only if the user explicitly overrides.

### 1. Source reputation & authenticity

- Prefer well-known packages: `axios`, `vue`, `vuetify`, `@tanstack/vue-query`, `@vueuse/core`, `pinia`, `vue-router`, `typescript`, `prettier`, `eslint`, official `@types/*`.
- Check the npm page (https://www.npmjs.com/package/<name>):
  - Verified publisher or well-known organization.
  - Weekly downloads > 10k (unless the ecosystem is niche and the alternative is clearly legitimate).
  - GitHub repository linked and matches the official source.
- **Typosquatting**: check the name for visual similarity to popular packages (e.g., `vuetfiy`, `axio`, `vuue`).

### 2. Community & maintenance health

- Last publish → less than 6 months ago (prefer < 3 months).
- GitHub repository: not archived, has recent commits, real issues/PRs.
- At least one stable release (`v1.0.0+`).
- Package is not deprecated on npm (`npm deprecate`).

### 3. Security record

- Run `pnpm audit` after install. If it reports critical/high vulnerabilities, investigate and warn.
- Check npm advisories page: https://www.npmjs.com/advisories?search=<package>
- No suspicious patterns:
  - Obfuscated code, minified blobs, encoded payloads.
  - Postinstall scripts that open network connections or read environment variables without documentation.
  - Hidden data collection or telemetry without explicit opt-in.

### 4. License

- MUST have a clear OSI-compliant license (MIT, Apache 2.0, BSD 3-clause, etc.).
- No license? Warn and ask user before proceeding.
- Avoid GPL/AGPL unless confirmed compatible.

### 5. Quality & API stability

- Prefer libraries with TypeScript types built-in or available via `@types/`.
- Prefer libraries that follow Tree Shaking-friendly patterns (ESM exports).
- Avoid packages that pull excessive transitive dependencies for a simple feature.
- Check bundle size impact: use https://bundlephobia.com/package/<name>.

### 6. Supply chain

- Check `package.json` dependencies of the target: does it pull in an excessive tree?
- Verify the package is hosted on npm's official registry (not a custom registry).
- Avoid packages with `prepare`, `preinstall`, or `postinstall` scripts unless clearly justified (these run arbitrary code during `pnpm install`).

## Procedure

```
1. Identify the best candidate package for the task
2. Run the checks above on the candidate
3. If it fails any check → search for a better alternative
4. If no safe alternative exists → inform the user with the risk details
5. Only then run pnpm add <package>
6. After adding: run make fe-check AND make fe-lint to verify no breakage
```

## Red flags (automatic rejection)

- < 100 weekly downloads AND < 5 contributors AND < 1 year of maintenance.
- Last publish > 1 year ago.
- Known critical/high CVE without a fix.
- Contains obfuscated or intentionally unreadable code.
- Zero documentation (no README, no API docs).
- Impersonates a well-known package (typosquatting).
- GitHub repo is archived, deleted, or returns 404.
- Package has postinstall scripts that do anything beyond building native modules.
