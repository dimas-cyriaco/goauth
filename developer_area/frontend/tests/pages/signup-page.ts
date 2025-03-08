import { faker } from '@faker-js/faker'
import type { Locator, Page } from '@playwright/test'

export class SignupPage {
  private readonly email: Locator
  private readonly password: Locator
  private readonly passwordConfirmation: Locator
  private readonly submit: Locator

  constructor(public readonly page: Page) {
    this.email = this.page.getByTestId('email')
    this.password = this.page.getByTestId('password')
    this.passwordConfirmation = this.page.getByTestId('password-confirmation')
    this.submit = this.page.getByTestId('submit')
  }

  async goto() {
    const link = this.page.getByTestId('link-to-signup')

    if (await link.isVisible()) {
      await link.click()
      return
    }

    await this.page.goto('/signup')
  }

  async fillEmail(text: string) {
    await this.email.fill(text)
  }

  async fillPassword(text: string) {
    await this.password.fill(text)
  }

  async fillPasswordConfirmation(text: string) {
    await this.passwordConfirmation.fill(text)
  }

  async clickSubmit(opts?: { noWait: boolean }) {
    await this.submit.click()

    const { noWait } = opts || {}
    if (noWait) {
      return
    }

    await this.page.waitForLoadState('networkidle')
  }

  async createUser(email: string, password: string): Promise<void> {
    await this.goto()
    await this.fillEmail(email)
    await this.fillPassword(password)
    await this.fillPasswordConfirmation(password)
    await this.clickSubmit()
  }

  static async createUser(
    page: Page,
  ): Promise<{ email: string; password: string }> {
    const signupPage = new SignupPage(page)

    const email = faker.internet.email()
    const password = faker.internet.password()

    await signupPage.goto()
    await signupPage.fillEmail(email)
    await signupPage.fillPassword(password)
    await signupPage.fillPasswordConfirmation(password)

    await signupPage.clickSubmit()

    return {
      email,
      password,
    }
  }
}
