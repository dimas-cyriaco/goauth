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

  async clickSubmit() {
    await this.submit.click()
  }

  async createUser(email: string, password: string): Promise<void> {
    await this.goto()
    await this.fillEmail(email)
    await this.fillPassword(password)
    await this.fillPasswordConfirmation(password)
    await this.clickSubmit()
  }
}
