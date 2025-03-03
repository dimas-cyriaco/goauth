import { faker } from '@faker-js/faker'
import { expect, test } from '@playwright/test'

import { SigninPage } from './pages/signin-page'
import { SignupPage } from './pages/signup-page'

test.describe('Signup Page', () => {
  test('should redirect to Home page on success', async ({ page }) => {
    // Arrange

    const email = faker.internet.email()
    const password = faker.internet.password()

    const signupPage = new SignupPage(page)
    await signupPage.createUser(email, password)

    const signinPage = new SigninPage(page)
    await signinPage.goto()

    // Act

    await signinPage.fillEmail(email)
    await signinPage.fillPassword(password)

    await signinPage.clickSubmit()

    // Assert

    await expect(page).toHaveURL('/')
  })

  test('should show error if email and password do not match', async ({
    page,
  }) => {
    // Arrange

    const signinPage = new SigninPage(page)

    const email = faker.internet.email()
    const password = faker.internet.password()

    const signupPage = new SignupPage(page)
    await signupPage.createUser(email, password)

    await signinPage.goto()

    // Act

    await signinPage.fillEmail(email)
    await signinPage.fillPassword('wrong-password')

    await signinPage.clickSubmit()

    // Assert

    await expect(page.getByTestId('login-error')).toHaveText(
      'Wrong email or password.',
    )
  })
})
