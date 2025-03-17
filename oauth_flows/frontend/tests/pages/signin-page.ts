import type { Locator, Page } from '@playwright/test'

export class SigninPage {
  private readonly email: Locator
  private readonly password: Locator
  private readonly submit: Locator

  constructor(public readonly page: Page) {
    this.email = this.page.getByTestId('email')
    this.password = this.page.getByTestId('password')
    this.submit = this.page.getByTestId('submit')
  }

  async goto() {
    const link = this.page.getByTestId('link-to-signin')

    if (await link.isVisible()) {
      await link.click()
      return
    }

    await this.page.goto('/signin')
  }

  async fillEmail(text: string) {
    await this.email.fill(text)
  }

  async fillPassword(text: string) {
    await this.password.fill(text)
  }

  async clickSubmit(opts?: { noWait: boolean }) {
    await this.submit.click()

    const { noWait } = opts || {}
    if (noWait) {
      return
    }

    await this.page.waitForURL('/')
  }

  async login(email: string, password: string): Promise<void> {
    await this.goto()
    await this.fillEmail(email)
    await this.fillPassword(password)
    await this.clickSubmit()
  }

  static async login(page: Page, email: string, password: string): Promise<void> {
    const signinPage = new SigninPage(page)

    await signinPage.goto()
    await signinPage.fillEmail(email)
    await signinPage.fillPassword(password)
    await signinPage.clickSubmit()
  }
}
