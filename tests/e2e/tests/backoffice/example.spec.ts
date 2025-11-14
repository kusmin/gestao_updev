import { test, expect } from '@playwright/test';

test.describe('Backoffice - Login', () => {
  test('redireciona usuários sem token para /login', async ({ page }) => {
    await page.goto('/');

    await expect(page).toHaveURL(/\/login$/);
    await expect(page.getByRole('heading', { name: /Login/i })).toBeVisible();
    await expect(page.getByLabel('Email')).toBeVisible();
    await expect(page.getByLabel('Password')).toBeVisible();
  });

  test('permite preencher o formulário de login', async ({ page }) => {
    await page.goto('/login');

    await page.getByLabel('Email').fill('admin@gestao.com');
    await page.getByLabel('Password').fill('supersecret');
    await expect(page.getByRole('button', { name: /Sign In/i })).toBeEnabled();
  });
});
