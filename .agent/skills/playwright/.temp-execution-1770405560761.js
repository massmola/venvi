const { chromium } = require('playwright');

const TARGET_URL = 'http://127.0.0.1:8090';

(async () => {
    const browser = await chromium.launch({ headless: true });
    const context = await browser.newContext();
    const page = await context.newPage();

    console.log('--- Testing "No Providers" Case ---');
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

    await page.goto(TARGET_URL);
    const loginBtn = page.locator('#login-btn');
    await loginBtn.click();

    const notification = page.locator('#auth-notification');
    await notification.waitFor({ state: 'visible', timeout: 5000 });
    const msg = await page.locator('#notification-msg').textContent();
    console.log('Notification message:', msg);

    if (msg.includes('No OAuth2 providers configured')) {
        console.log('✅ Pass: Correct message shown');
    } else {
        console.log('❌ Fail: Unexpected message:', msg);
    }

    console.log('\n--- Testing Loading State ---');
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

    // Mock the OAuth2 flow to delay
    await page.route('**/api/collections/users/auth-with-oauth2', async route => {
        await new Promise(resolve => setTimeout(resolve, 2000));
        await route.fulfill({
            status: 200,
            contentType: 'application/json',
            body: JSON.stringify({ token: 'fake-token', record: { id: '123', name: 'Test User' } })
        });
    });

    await page.reload();
    await loginBtn.click();

    const loginText = await page.locator('#login-text').textContent();
    console.log('Button text during auth:', loginText);
    const isDisabled = await loginBtn.isDisabled();
    console.log('Button is disabled:', isDisabled);

    if (loginText === 'Authorizing...' && isDisabled) {
        console.log('✅ Pass: Loading state correctly displayed');
    } else {
        console.log('❌ Fail: Expected "Authorizing..." and disabled, got:', loginText, isDisabled);
    }

    await browser.close();
})();
