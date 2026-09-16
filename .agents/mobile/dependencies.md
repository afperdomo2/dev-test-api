# Mobile dependency safety rules

Same safety checks as [`.agents/web/dependencies.md`](../web/dependencies.md). Before adding any npm package via `pnpm add` in `apps/mobile/`:

1. **Source reputation**: npm weekly downloads > 10k, verified publisher, GitHub linked
2. **Maintenance**: last publish < 6 months, not deprecated, not archived
3. **Security**: `pnpm audit`, no obfuscated code, no suspicious postinstall scripts
4. **License**: OSI-compliant (MIT, Apache 2.0, BSD)
5. **Quality**: TypeScript types, tree-shakeable, no excessive transitive deps

**Additional mobile-specific checks:**

- Prefer Expo-compatible packages (`expo-*`, or libraries known to work with Expo managed workflow)
- Check if a native module requires `expo prebuild` or ejecting — flag to the user
- Avoid packages that pull excessive native dependencies (check `react-native` peer compatibility)

## Procedure

```
1. Identify candidate package
2. Run checks above
3. If fails any check → search alternative
4. If no safe alternative → warn user with risk details
5. Run pnpm add in workspace root (or apps/mobile/)
6. After adding: run pnpm --filter mobile typecheck to verify
```

## Red flags (automatic rejection)

- Same as web: < 100 downloads, > 1 year stale, archived repo, postinstall scripts, obfuscated code
- Package assumes `node` APIs not available in React Native (e.g., `fs`, `path`, `localStorage`)
