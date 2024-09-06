import { test, expect } from '@playwright/test';

var username = "newuser"
var baseURL = 'http://localhost:8080/login'

test.beforeEach(async ({page}) => {
    await page.goto(baseURL);
    await page.getByRole('button', { name: 'Register' }).click();
    await page.getByPlaceholder('Enter username').click();
});

function makeid(length) {
    let result = '';
    const characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    const charactersLength = characters.length;
    let counter = 0;
    while (counter < length) {
      result += characters.charAt(Math.floor(Math.random() * charactersLength));
      counter += 1;
    }
    return result;
}

test('test duplicate registration', async ({ page }) => {
  await page.getByPlaceholder('Enter username').fill(username);
  await page.getByPlaceholder('Enter password').click();
  await page.getByPlaceholder('Enter password').fill('123456');
  await page.getByRole('button', { name: 'Register' }).click();

  await expect(page.locator("text=Username is already being used")).toBeVisible()
});

test('test valid registration', async ({page}) => {
    var randomUsername = makeid(6)
    await page.getByPlaceholder('Enter username').fill(randomUsername);
    await page.getByPlaceholder('Enter password').click();
    await page.getByPlaceholder('Enter password').fill('123456');
    await page.getByRole('button', { name: 'Register' }).click();
  
    await expect(page.locator("text=Blog Aggregator")).toBeVisible()
});

