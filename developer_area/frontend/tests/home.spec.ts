import { expect, test } from '@playwright/test'

test.describe('Home Page', () => {
  test('should have title', async ({ page }) => {
    // Act

    await page.goto('/signup')

    // Assert

    await expect(page).toHaveTitle('GOAuth')
  })
})
