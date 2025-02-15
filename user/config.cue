SendEmails: bool | *true

if #Meta.Environment.Type == "test" {
    SendEmails: false
}

SendEmailsFrom: [
    if #Meta.Environment.Type == "production" { "noreply@example.com" },
    if #Meta.Environment.Name == "staging"    { "staging@example.com" },

    "dev-system@example.dev",
][0]

SMTPHost: "smtp.gmail.com"
SMTPPort: 587
SMTPUsername: "dimascyriaco@gmail.com"

BaseURL: [
    if #Meta.Environment.Type == "production" { "https://goauth.com" },
    if #Meta.Environment.Name == "staging"    { "https://staging.goauth.com" },

    "http://localhost:4000",
][0]
