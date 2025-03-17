import { expect, test } from '@playwright/test'

import { SigninPage } from './pages/signin-page'
import { SignupPage } from './pages/signup-page'

test.describe('Signup Page', () => {
  test('should show login page if user is not logged in', async ({ page }) => {
    // Act

    await page.goto('/')

    // Assert

    await expect(page.getByTestId('link-to-signin')).toBeVisible()
  })

  test('should redirect to Home page on success', async ({ page }) => {
    // Arrange

    const { email, password } = await SignupPage.createUser(page)

    const signinPage = new SigninPage(page)
    await signinPage.goto()

    // Act

    await signinPage.fillEmail(email)
    await signinPage.fillPassword(password)

    await signinPage.clickSubmit()

    // Assert

    await expect(page).toHaveURL('/')
  })

  test('should not show login link if user is logged in', async ({ page }) => {
    // Arrange

    const { email, password } = await SignupPage.createUser(page)
    await SigninPage.login(page, email, password)

    // Act

    await page.getByTestId('link-to-home').click()

    // Assert

    await expect(page.getByTestId('link-to-signin')).not.toBeVisible()
  })

  test('should keep login state on page reload', async ({ page }) => {
    // Arrange

    const { email, password } = await SignupPage.createUser(page)
    await SigninPage.login(page, email, password)

    // Act

    await page.goto('/')

    // Assert

    await expect(page.getByTestId('link-to-signin')).not.toBeVisible()
  })

  test('should show error if email and password do not match', async ({
    page,
  }) => {
    // Arrange

    const { email } = await SignupPage.createUser(page)

    const signinPage = new SigninPage(page)
    await signinPage.goto()

    // Act

    await signinPage.fillEmail(email)
    await signinPage.fillPassword('wrong-password')

    await signinPage.clickSubmit({ noWait: true })

    // Assert

    await expect(page.getByTestId('login-error')).toHaveText(
      'Wrong email or password.',
    )
  })
})
