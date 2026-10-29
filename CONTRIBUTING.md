# Contributing to steampipe-plugin-anthropic

Thanks for your interest. Bug reports, new tables for Anthropic API endpoints,
query examples and docs fixes are all welcome.

## Before you start

- For a new table or anything bigger than a small fix, open an issue first so
  we can agree on the approach.
- Security issues: do not open a public issue; see [SECURITY.md](SECURITY.md).
- This project follows the [Code of Conduct](CODE_OF_CONDUCT.md).
- Never commit API keys, organisation IDs, email addresses or real query
  results. Use made-up values in docs and examples.

## Development setup

You need Go (see `go.mod` for the version) and
[Steampipe](https://steampipe.io/downloads).

```bash
git clone https://github.com/huntbase-io/steampipe-plugin-anthropic.git
cd steampipe-plugin-anthropic
go build ./...
make install                                # builds into ~/.steampipe/plugins/local/anthropic/
cp config/anthropic.spc ~/.steampipe/config/anthropic.spc
steampipe service restart                   # pick up a new binary
steampipe query "select id, display_name from anthropic_model"
```

CI runs the same checks on every pull request:

```bash
go build ./...
go vet ./...
gofmt -l .        # must print nothing
```

`make sqlite` and `make postgres` build the SQLite extension and the Postgres
FDW; see the [README](README.md#sqlite-extension).

## Pull requests

- Keep each PR focused on one change.
- Make sure `go build`, `go vet` and `gofmt` are clean before pushing.
- Test new or changed tables against the live API with `make install`, and
  say in the PR which key type you used and which queries you ran.
- Add or update the table's page in `docs/tables/` and, for user-visible
  changes, an entry in [CHANGELOG.md](CHANGELOG.md).

## Licence

By contributing you agree that your contributions are licensed under the
project's [Apache License 2.0](LICENSE).
