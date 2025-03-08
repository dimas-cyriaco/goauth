import { faker } from '@faker-js/faker'
import { expect, test } from '@playwright/test'

import { ApplicationNewPage } from './pages/application-new-page'
import { SigninPage } from './pages/signin-page'
import { SignupPage } from './pages/signup-page'

test.describe('New Applicatino Page', () => {
  test('should create new Application', async ({ page }) => {
    // Arrange

    const { email, password } = await SignupPage.createUser(page)
    await SigninPage.login(page, email, password)

    const name = faker.company.name()

    const newPage = new ApplicationNewPage(page)

    // Act

    // TODO:Implement!
    // TODO:Stop sending email on test
    await newPage.goto()
    await newPage.fillName(name)
    await newPage.clickSubmit()

    // Assert

    await expect(page).toHaveURL('/applications')

    await expect(page.getByTestId('application-name')).toHaveText(name)
  })
})
