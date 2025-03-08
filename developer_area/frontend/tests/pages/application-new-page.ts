import type { Locator, Page } from '@playwright/test'

export class ApplicationNewPage {
  private readonly name: Locator
  private readonly submit: Locator

  constructor(public readonly page: Page) {
    this.name = this.page.getByTestId('name')
    this.submit = this.page.getByTestId('submit')
  }

  async goto() {
    const link = this.page.getByTestId('link-to-create-application')

    if (await link.isVisible()) {
      await link.click()
      return
    }

    await this.page.goto('/applications/new')
  }

  async fillName(text: string) {
    await this.name.fill(text)
  }

  async clickSubmit(opts?: { noWait: boolean }) {
    await this.submit.click()

    const { noWait } = opts || {}
    if (noWait) {
      return
    }

    await this.page.waitForURL('/applications')
  }

  async create(name: string): Promise<void> {
    await this.goto()
    await this.fillName(name)
    await this.clickSubmit()
  }
}
