import { faker } from '@faker-js/faker'
import { expect, test } from '@playwright/test'

import { SignupPage } from './pages/signup-page'

test.describe('Signup Page', () => {
  test('should redirect to Signin page on success', async ({ page }) => {
    // Arrange

    const signupPage = new SignupPage(page)

    await signupPage.goto()

    const password = faker.internet.password()

    // Act

    await signupPage.fillEmail(faker.internet.email())
    await signupPage.fillPassword(password)
    await signupPage.fillPasswordConfirmation(password)

    await signupPage.clickSubmit()

    // Assert

    await expect(page).toHaveURL('/signin')
  })

  test('should show email format error', async ({ page }) => {
    // Arrange

    const signupPage = new SignupPage(page)

    await signupPage.goto()

    const password = faker.internet.password()

    // Act

    await signupPage.fillEmail('invalid@email')
    await signupPage.fillPassword(password)
    await signupPage.fillPasswordConfirmation(password)

    await signupPage.clickSubmit()

    // Assert

    await expect(page.getByTestId('email-error')).toHaveText(
      'Email must be a valid email address',
    )
  })

  test('should show email taken error', async ({ page }) => {
    // Arrange

    const signupPage = new SignupPage(page)

    const email = faker.internet.email()
    const password = faker.internet.password()

    await signupPage.createUser(email, password)

    await signupPage.goto()

    // Act

    await signupPage.fillEmail(email)
    await signupPage.fillPassword(password)
    await signupPage.fillPasswordConfirmation(password)
    await signupPage.clickSubmit()

    // Assert

    await expect(page.getByTestId('email-error')).toHaveText(
      'Email already taken',
    )
  })

  test('should show password too short error', async ({ page }) => {
    // Arrange

    const signupPage = new SignupPage(page)
    await signupPage.goto()

    const password = '123'

    // Act

    await signupPage.fillEmail(faker.internet.email())
    await signupPage.fillPassword(password)
    await signupPage.fillPasswordConfirmation(password)
    await signupPage.clickSubmit()

    // Assert

    await expect(page.getByTestId('password-error')).toHaveText(
      'Password must be at least 6 characters in length',
    )
  })

  test('should show password too long error', async ({ page }) => {
    // Arrange

    const signupPage = new SignupPage(page)
    await signupPage.goto()

    const password = 'a'.repeat(73)

    // Act

    await signupPage.fillEmail(faker.internet.email())
    await signupPage.fillPassword(password)
    await signupPage.fillPasswordConfirmation(password)
    await signupPage.clickSubmit()

    // Assert

    await expect(page.getByTestId('password-error')).toHaveText(
      'Password must be a maximum of 72 characters in length',
    )
  })

  test('should show password confirmation error', async ({ page }) => {
    // Arrange

    const signupPage = new SignupPage(page)
    await signupPage.goto()

    // Act

    await signupPage.fillEmail('test@example.com')
    await signupPage.fillPassword('password123')
    await signupPage.fillPasswordConfirmation('wrong')
    await signupPage.clickSubmit()

    // Assert

    await expect(page.getByTestId('password-confirmation-error')).toHaveText(
      'PasswordConfirmation must be equal to Password',
    )
  })
})
