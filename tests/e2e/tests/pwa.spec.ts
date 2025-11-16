import { test, expect } from '@playwright/test';

test.describe('Funcionalidade PWA do Backoffice', () => {
  const BACKOFFICE_URL = 'http://localhost:5174'; // Ajuste para a URL correta do seu backoffice

  test('deve carregar offline após a visita inicial', async ({ page }) => {
    // 1. Visitar a página online para que o Service Worker seja registrado e os assets cacheados
    await page.goto(BACKOFFICE_URL);
    await page.waitForLoadState('networkidle'); // Espera a rede ficar ociosa para garantir que o SW registrou e cacheou

    // 2. Simular modo offline
    await page.context().setOffline(true);

    // 3. Recarregar a página e verificar se ela carrega
    await page.reload();
    // Verificar se o elemento raiz da aplicação está visível, indicando que carregou
    await expect(page.locator('#root')).toBeVisible();
    // Você pode adicionar verificações mais específicas aqui, como a presença de um texto ou componente chave
    // await expect(page.getByText('Bem-vindo ao Backoffice')).toBeVisible(); // Exemplo, descomente e ajuste se houver um texto específico
  });

  test('deve ter um Web App Manifest válido', async ({ page }) => {
    await page.goto(BACKOFFICE_URL);

    // Obter a URL do manifest
    const manifestLink = page.locator('link[rel="manifest"]');
    await expect(manifestLink).toBeAttached(); // Garante que a tag manifest existe
    const manifestUrl = await manifestLink.getAttribute('href');
    expect(manifestUrl).toBeDefined();

    // Fazer uma requisição para o manifest e verificar seu conteúdo
    const manifestResponse = await page.goto(new URL(manifestUrl!, BACKOFFICE_URL).toString());
    const manifest = await manifestResponse?.json();

    expect(manifest).toBeDefined();
    expect(manifest.name).toBe('Gestão UpDev Backoffice');
    expect(manifest.short_name).toBe('Backoffice');
    expect(manifest.description).toBe('Aplicação de backoffice para Gestão UpDev');
    expect(manifest.theme_color).toBe('#ffffff');
    expect(manifest.icons).toBeInstanceOf(Array);
    expect(manifest.icons.length).toBeGreaterThan(0);
    expect(manifest.icons.some((icon: any) => icon.sizes === '512x512' && icon.purpose === 'any maskable')).toBeTruthy();
  });
});

test.describe('Funcionalidade PWA do Frontend', () => {
  const FRONTEND_URL = 'http://localhost:5173'; // Ajuste para a URL correta do seu frontend

  test('deve carregar offline após a visita inicial', async ({ page }) => {
    await page.goto(FRONTEND_URL);
    await page.waitForLoadState('networkidle');
    await page.context().setOffline(true);
    await page.reload();
    await expect(page.locator('#root')).toBeVisible();
    // await expect(page.getByText('Bem-vindo ao Gestão UpDev')).toBeVisible(); // Exemplo, descomente e ajuste se houver um texto específico
  });

  test('deve ter um Web App Manifest válido', async ({ page }) => {
    await page.goto(FRONTEND_URL);
    const manifestLink = page.locator('link[rel="manifest"]');
    await expect(manifestLink).toBeAttached();
    const manifestUrl = await manifestLink.getAttribute('href');
    expect(manifestUrl).toBeDefined();

    const manifestResponse = await page.goto(new URL(manifestUrl!, FRONTEND_URL).toString());
    const manifest = await manifestResponse?.json();

    expect(manifest).toBeDefined();
    expect(manifest.name).toBe('Gestão UpDev Frontend');
    expect(manifest.short_name).toBe('Frontend');
    expect(manifest.description).toBe('Aplicação frontend para Gestão UpDev');
    expect(manifest.theme_color).toBe('#ffffff');
    expect(manifest.icons).toBeInstanceOf(Array);
    expect(manifest.icons.length).toBeGreaterThan(0);
    expect(manifest.icons.some((icon: any) => icon.sizes === '512x512' && icon.purpose === 'any maskable')).toBeTruthy();
  });
});
