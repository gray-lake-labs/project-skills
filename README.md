# scrum-master — moved

This repository is archived. Scrum Master — `tpm`, the process machines it runs, and the specification corpus that governs both — now lives in the declarative-agents monorepo as an application:

**[gray-lake-labs/declarative-agents › applications/scrum-master](https://github.com/gray-lake-labs/declarative-agents/tree/main/applications/scrum-master)**

It moved with its full history (`git log` in the monorepo follows every file through the move), and its specification, issues, and releases continue there.

## Install

Releases of the monorepo carry `tpm` and the `agent` runtime together, under one `SHA256SUMS`:

```bash
curl -fsSL https://raw.githubusercontent.com/gray-lake-labs/declarative-agents/main/applications/scrum-master/install.sh | sh
```

The releases this repository published remain downloadable here, so a pinned install of an older `tpm` keeps working.

## Why it moved

The decision is recorded in [#189](https://github.com/gray-lake-labs/scrum-master/issues/189) and the move in [#190](https://github.com/gray-lake-labs/scrum-master/issues/190): tpm's skills became declarative-agents profiles, and the runtime, the format library, and the corpus now move together at one tag instead of pinning one another across repositories.
