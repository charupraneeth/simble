(function () {
    // 1. Get the current script tag and extract the data-domain attribute
    const scriptTag = document.currentScript;
    const domain = scriptTag.getAttribute('data-domain');

    // 2. Send the analytics hit
    fetch("https://simble.up.railway.app/api/event", {
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

})()