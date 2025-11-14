import { test, expect } from '@playwright/test';

test.describe('Frontend - Clientes', () => {
  test('exibe cabeçalho e CTA', async ({ page }) => {
    await page.goto('/');

    await expect(page.getByRole('heading', { name: /Clientes/i })).toBeVisible();
    await expect(
      page.getByRole('button', { name: /Adicionar Cliente/i }),
    ).toBeVisible();
  });

  test('renderiza colunas da tabela', async ({ page }) => {
    await page.goto('/');

    const columns = ['Nome', 'Email', 'Telefone', 'Ações'];
    for (const column of columns) {
      await expect(
        page.getByRole('columnheader', { name: column }),
      ).toBeVisible();
    }
  });
});
