const { chromium } = require('/root/EmpoweredPixels/frontend/node_modules/playwright');

(async () => {
    let browser;
    try {
        browser = await chromium.launch({ headless: true });
        const context = await browser.newContext();
        const page = await context.newPage();

        page.on('console', msg => {
            console.log('[CONSOLE][' + msg.type() + '] ' + msg.text());
        });

        page.on('response', response => {
             console.log('[RESPONSE] ' + response.status() + ' ' + response.url());
        });

        console.log('Navigating to leagues...');
        await page.goto('http://localhost:49100/leagues');
        
        console.log('Waiting for "Connection lost" indicator...');
        
        // Wait for potential error text to appear
        try {
            await page.waitForSelector('text=/Connection lost/i', { timeout: 15000 });
            console.log('FOUND: "Connection lost" error message on page.');
        } catch (e) {
            console.log('Could not find "Connection lost" text on page via selector.');
        }

        await page.screenshot({ path: '/root/EmpoweredPixels/league_debug_screenshot.png' });
        
        console.log('Poking the backend directly...');
        // Try fetching the league endpoint directly from the node environment
    } catch (error) {
        console.error('Execution failed:', error);
    } finally {
        if (browser) await browser.close();
    }
})();
