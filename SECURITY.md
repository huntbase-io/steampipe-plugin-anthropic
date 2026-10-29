# Security policy

## Reporting a vulnerability

Please report vulnerabilities privately through GitHub's
[private vulnerability reporting](https://github.com/huntbase-io/steampipe-plugin-anthropic/security/advisories/new)
("Report a vulnerability" on the repository's **Security** tab). Do not open a
public issue or pull request for a suspected vulnerability.

Include the plugin version (or commit), how you run it (Steampipe, the SQLite
extension or the Postgres FDW), and steps to reproduce. We will acknowledge
your report within 5 working days, keep you updated while we look into it,
and credit you in the advisory unless you prefer otherwise.

Security fixes are made on the latest release only.

## Scope

The plugin reads the Anthropic API with the keys you configure and returns the
results as tables. In scope are issues in the plugin itself, for example:

- an API key (`api_key`, `admin_api_key`, `compliance_api_key`,
  `analytics_api_key` or their environment variables) leaking into query
  results, logs, error messages or anywhere other than the request to the
  Anthropic API;
- the plugin sending a key or request to a host other than the Anthropic API;
- a qualifier or other query input changing an API request in a way the table
  does not document;
- vulnerabilities in the published plugin binaries or the SQLite and Postgres
  extensions.

Vulnerabilities in the Anthropic API or claude.ai should be reported to
Anthropic through its responsible disclosure programme. Issues in Steampipe,
the Steampipe plugin SDK, steampipe-sqlite or steampipe-postgres-fdw should be
reported to [Turbot](https://github.com/turbot).

## Using the plugin safely

**Use the least-privileged key that covers the tables you query.** The admin,
compliance and analytics keys each unlock organisation-wide data. Configure
only the keys you need.

**Treat query results as sensitive.** The `anthropic_compliance_*` tables can
return chat content, project files and user details for a whole organisation.
Steampipe caches results; restrict who can query the connection and where its
output goes.

**Keep keys out of shared files.** Prefer the environment variables to
writing keys into `anthropic.spc`, and keep that file out of version control.
