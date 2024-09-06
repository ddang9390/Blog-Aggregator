import { test, expect } from '@playwright/test';

var username = "newuser"
var baseURL = 'http://localhost:8080/login'

test.beforeEach(async ({page}) => {
    await page.goto(baseURL);
    await page.getByPlaceholder('Enter username').click();
    await page.getByPlaceholder('Enter username').fill(username);
    await page.getByPlaceholder('Enter password').click();
});

test('test invalid login', async ({ page }) => {
  await page.getByPlaceholder('Enter password').fill('123grdhbdff456');
  await page.getByRole('button', { name: 'Login' }).click();

  await expect(page.locator("text=Invalid username or password")).toBeVisible()
});


test('test valid login', async ({ page }) => {
  await page.getByPlaceholder('Enter password').fill('123456');
  await page.getByRole('button', { name: 'Login' }).click();

  await expect(page.locator("text=Welcome " + username)).toBeVisible()
});