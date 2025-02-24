POST {{base_url}}/login

{
  "Email": "{{user_email}}",
  "Password": "{{user_password}}"
}

HTTP 200

# Capture the CSRF token value from html body.
[Captures]
session_token: cookie "session_token"
csrf_token: cookie "csrf_token"

GET {{base_url}}/users/1
X-CSRF-TOKEN: {{csrf_token}}

HTTP 200

