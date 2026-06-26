# zed-bedrock-env

Provisions a scoped IAM role for AWS Bedrock and exports temporary credentials into your shell so [Zed](https://zed.dev) can call Bedrock models directly.

## How it works

1. **`infra/`** — Terraform creates an IAM role with a minimal inline policy that allows `bedrock:InvokeModel` and `bedrock:InvokeModelWithResponseStream` for a specific model.
2. **`app/`** — A small Go CLI (`zed-creds`) assumes that role via STS (resolving the account ID automatically via `GetCallerIdentity`) and prints `export` statements for the resulting temporary credentials.
3. Eval the output in your shell before launching Zed.

## Prerequisites

- AWS credentials configured locally (e.g. `~/.aws/credentials`, SSO, or environment variables)
- [Terraform](https://developer.hashicorp.com/terraform) >= 1.0
- [Go](https://go.dev) >= 1.26 (or [Task](https://taskfile.dev) for the build shortcuts)

## 1 — Create the IAM role

Use a Terraform workspace per AWS account to keep state isolated:

```bash
cd infra
terraform workspace new my-project   # once per account; use "select" afterwards
terraform init
terraform apply
```

> **Note:** The role must exist before running `zed-creds`.  
> If you want to use a different Bedrock model, update the `Resource` ARNs in `infra/main.tf` and the `role` / `role_name` values in `app/config.yaml` and `app/config/config.go` (`DefaultRoleName`) to match the new role name.

## 2 — Configure the app

Copy and adjust the example config:

```bash
cp app/config-example.yaml app/config.yaml
```

Minimal `config.yaml`:

```yaml
role: bedrock-direct-call-sonnet-4-6   # must match the role created by Terraform

aws:
  region: eu-central-1
  session_name: zed-cli
  duration_seconds: 3600
```

The account ID is resolved automatically via `sts:GetCallerIdentity` — no need to hard-code it.  
You can still set `role_arn` explicitly if you need a cross-account or non-standard ARN.

## 3 — Build and install

```bash
task build        # builds to bin/zed-creds
task install      # builds and installs to /usr/local/bin/zed-creds
```

Or directly with Go:

```bash
cd app && go build -o ../bin/zed-creds main.go
```

## 4 — Export credentials

```bash
eval $(zed-creds)
```

This outputs `ZED_ACCESS_KEY_ID`, `ZED_SECRET_ACCESS_KEY`, `ZED_SESSION_TOKEN`, and `ZED_AWS_REGION` in your current shell session. Launch Zed from the same shell.

Credentials are valid for `duration_seconds` (default 1 hour). Re-run to refresh.
