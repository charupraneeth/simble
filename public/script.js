(function () {
    // 1. Get the current script tag and extract the data-domain attribute
    const scriptTag = document.currentScript;
    const domain = scriptTag.getAttribute('data-domain');

    // Dynamically get the API URL based on where the script was loaded from
    const apiUrl = new URL('/api/event', scriptTag.src).href;

    // Helper to send the analytics payload
    function sendEvent() {
        fetch(apiUrl, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                "domain": domain,
                "name": "pageview",
                "referer": document.referrer || null,
                "url": window.location.href
            })
        }).catch((err) => {
            // Fail silently so we don't spam the website's console
        });
    }

    // 2. Track initial page load
    sendEvent();

    // 3. Track SPA (Single Page Application) route changes
    let lastUrl = window.location.href;

    function handleRouteChange() {
        if (lastUrl !== window.location.href) {
            lastUrl = window.location.href;
            sendEvent();
        }
    }

    // Listen for browser back/forward buttons
    window.addEventListener('popstate', handleRouteChange);

    // Monkey-patch pushState and replaceState to catch Vue/React/SPA routing
    const originalPushState = history.pushState;
    history.pushState = function (...args) {
        originalPushState.apply(this, args);
        handleRouteChange();
    };

    const originalReplaceState = history.replaceState;
    history.replaceState = function (...args) {
        originalReplaceState.apply(this, args);
        handleRouteChange();
    };})()