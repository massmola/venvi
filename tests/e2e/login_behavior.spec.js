const { test, expect } = require('@playwright/test');

test.describe('Login Feature', () => {
    test.beforeEach(async ({ page }) => {
        // Navigate to the app
        await page.goto('http://127.0.0.1:8090');
    });

    test('should show informative message when no OAuth2 providers are configured', async ({ page }) => {
        // Mock PocketBase API call for auth methods to return empty
        await page.route('**/api/collections/users/auth-methods', async route => {
            await route.fulfill({
                status: 200,
                contentType: 'application/json',
                body: JSON.stringify({
                    authProviders: [],
                    usernamePassword: true
                })
            });
        });

        const loginBtn = page.locator('#login-btn');
        await expect(loginBtn).toBeVisible();
        await expect(loginBtn).toHaveText('Login');

        // Click login
        await loginBtn.click();

        // Check for notification
        const notification = page.locator('#auth-notification');
        await expect(notification).toBeVisible();
        await expect(notification).toContainText('No OAuth2 providers configured');

        // Wait for it to disappear or just check initial state
        await expect(loginBtn).not.toBeDisabled();
        await expect(loginBtn).toHaveText('Login');
    });

    test('should show loading state during auth flow', async ({ page }) => {
        // Mock PocketBase API call for auth methods to return one provider
        await page.route('**/api/collections/users/auth-methods', async route => {
            await route.fulfill({
                status: 200,
                contentType: 'application/json',
                body: JSON.stringify({
                    authProviders: [{ name: 'github', state: '123' }],
                    usernamePassword: true
                })
            });
        });

        // Mock the OAuth2 flow to hang/delayed
        await page.route('**/api/collections/users/auth-with-oauth2', async route => {
            // Wait 1 second before fulfilling to see loading state
            await new Promise(resolve => setTimeout(resolve, 1000));
            await route.fulfill({
                status: 200,
                contentType: 'application/json',
                body: JSON.stringify({ token: 'fake-token', record: { id: '123', name: 'Test User' } })
            });
        });

        const loginBtn = page.locator('#login-btn');
        const loginText = page.locator('#login-text');

        await loginBtn.click();

        // Check for loading state text change
        // Note: It might change fast to "Authorizing..."
        await expect(loginText).not.toHaveText('Login');
        await expect(loginBtn).toBeDisabled();

        // After fulfillment (handled by Playwright waiting)
        // The page would reload in reality, but we can just check it stayed in loading for a bit
    });
});
