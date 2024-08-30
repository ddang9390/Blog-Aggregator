import { test, expect } from '@playwright/test';

var username = "newuser"
test('test invalid login', async ({ page }) => {
  await page.goto('http://localhost:8080/login');
  await page.getByPlaceholder('Enter username').click();
  await page.getByPlaceholder('Enter username').fill(username);
  await page.getByPlaceholder('Enter password').click();
  await page.getByPlaceholder('Enter password').fill('123grdhbdff456');
  await page.getByRole('button', { name: 'Login' }).click();

  await expect(page.locator("text=Invalid username or password")).toBeVisible()
});

test('test invalid registration', async ({ page }) => {
  await page.goto('http://localhost:8080/login');
  await page.getByRole('button', { name: 'Register' }).click();
  await page.getByPlaceholder('Enter username').click();
  await page.getByPlaceholder('Enter username').fill(username);
  await page.getByPlaceholder('Enter password').click();
  await page.getByPlaceholder('Enter password').fill('123456');
  await page.getByRole('button', { name: 'Register' }).click();

  await expect(page.locator("text=Username is already being used")).toBeVisible()
});

test('test valid login', async ({ page }) => {
  await page.goto('http://localhost:8080/login');
  await page.getByPlaceholder('Enter username').click();
  await page.getByPlaceholder('Enter username').fill(username);
  await page.getByPlaceholder('Enter password').click();
  await page.getByPlaceholder('Enter password').fill('123456');
  await page.getByRole('button', { name: 'Login' }).click();

  await expect(page.locator("text=Welcome " + username)).toBeVisible()
});