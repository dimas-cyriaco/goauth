## Usage

```bash
npm install # or pnpm install or yarn install
```

### Learn more on the [Solid Website](https://solidjs.com) and come chat with us on our [Discord](https://discord.com/invite/solidjs)

## Available Scripts

In the project directory, you can run:

### `npm run dev`

Runs the app in the development mode.<br>
Open [http://localhost:5173](http://localhost:5173) to view it in the browser.

### `npm run build`

Builds the app for production to the `dist` folder.<br>
It correctly bundles Solid in production mode and optimizes the build for the best performance.

The build is minified and the filenames include the hashes.<br>
Your app is ready to be deployed!

## Deployment

Learn more about deploying your application with the [documentations](https://vite.dev/guide/static-deploy.html)

The call for `client.user.Registration(params)` on the `onSubmit` of this component throw exceptions on this format:

```json
{
    "code": "invalid_argument",
    "message": "Email already taken",
    "details": null
}
```

The `message` field is a string, but can contain multiple error. For example: `validation failed: Key: 'RegistrationParams.Email' Error:Field validation for 'Email' failed on the 'email' tag\nKey: 'RegistrationParams.PasswordConfirmation' Error:Field validation for 'PasswordConfirmation' failed on the 'eqcsfield' tag`. This message says the email does is not valid due to it's format, and the password confirmation do not match the password.

I want to map this error to an local state on this format:

```json
{
    [<field>]: string[], // An array of all error on the field
}
```

This are the mapping I want to do

- "validation failed: Key: 'RegistrationParams.Email' Error:Field validation for 'Email' failed on the 'email' tag": "Invalid email format"
- "Key: 'RegistrationParams.PasswordConfirmation' Error:Field validation for 'PasswordConfirmation' failed on the 'eqcsfield' tag": Password onfirmation do not match password.
- "Email already taken": "Email already taken"
- "validation failed: Key: 'RegistrationParams.Password' Error:Field validation for 'Password' failed on the 'min' tag": "Password too short. Should be at least 6 characters"
- "validation failed: Key: 'RegistrationParams.Password' Error:Field validation for 'Password' failed on the 'max' tag": "Password too long. Should be at most 72 characters"
