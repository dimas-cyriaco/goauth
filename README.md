# 🔒 GOAuth

An experimental OAuth 2.0 server made for learning Go, Encore, and the
intricacies of OAuth.

## ⚡️ Requirements

The project uses [Mise](https://mise.jdx.dev/) for tools version managing and
task runner.

It's necessary to install some extra tools that Mise do not have support yet:

1. [Encore](https://encore.dev/go): The framework used for the backend.
2. [Process Compose](https://f1bonacc1.github.io/process-compose/installation/):
   A tool to orchestrate process inspired by docker compose.

After installing encore, run `encore auth login` to setup the environment and
allow you to fetch the projects secrets.

## 🏃 Running

After installing the requirements, you are ready to run the project. Run:

```sh
mise pc:up
```

This will start both, the backend and the frontend, run the backend tests in
watch mode and start the Playwright UI.

Run the following command to run the Process Compose UI:

```sh
mise pc:ui
```

## 🤫 Secrets

This step is optional, and is only needed to run the CI tasks locally.

Some secrets are managed through
[Mise](https://mise.jdx.dev/environments/secrets.html) and the encrypted file
`.env.secrets.json`. If you don't have access to the encryption key, you can
recreate this file with your own data using `sops`:

1. Edit the `.env.secrets.json` to be on the format:

```json
{
    "ENCORE_AUTH_KEY": "key"
}
```

2. Run: `sops encrypt --age age-key .env.secrets.json --inplace`

This will encrypt the file and allow Mise to set the variables automatically.

## Systems

### [WIP] Developer Area

This is where the developers create/manage their apps

### [TODO] OAuth Flows
