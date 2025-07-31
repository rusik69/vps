document.addEventListener('DOMContentLoaded', function() {
    const form = document.getElementById('shortenForm');
    const urlInput = document.getElementById('url');
    const customCheckbox = document.getElementById('custom');
    const captchaContainer = document.getElementById('captcha');
    const resultDiv = document.getElementById('result');
    const shortUrlInput = document.getElementById('shortUrl');
    const copyButton = document.getElementById('copyButton');

    // Initialize reCAPTCHA
    let grecaptcha;
    let recaptchaWidget;

    window.onloadCallback = function() {
        recaptchaWidget = grecaptcha.render('recaptcha-container', {
            'sitekey': 'YOUR_RECAPTCHA_SITE_KEY',
            'callback': function(response) {
                // Handle successful verification
            }
        });
    };

    // Handle form submission
    form.addEventListener('submit', async function(e) {
        e.preventDefault();

        const url = urlInput.value;
        const customCode = customCheckbox.checked ? urlInput.value.split('/').pop() : '';

        try {
            // Show CAPTCHA
            captchaContainer.classList.remove('hidden');
            
            // Wait for CAPTCHA verification
            const token = await new Promise((resolve) => {
                grecaptcha.ready(() => {
                    grecaptcha.execute('YOUR_RECAPTCHA_SITE_KEY', {action: 'shorten'})
                        .then(resolve);
                });
            });

            // Make API request
            const response = await fetch('/api/shorten', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    url,
                    customCode,
                    captchaToken: token
                })
            });

            if (!response.ok) {
                throw new Error('Failed to shorten URL');
            }

            const data = await response.json();
            shortUrlInput.value = data.short_url;
            resultDiv.classList.remove('hidden');
            copyButton.addEventListener('click', () => {
                navigator.clipboard.writeText(data.short_url)
                    .then(() => {
                        copyButton.textContent = 'Copied!';
                        setTimeout(() => {
                            copyButton.textContent = 'Copy URL';
                        }, 2000);
                    });
            });
        } catch (error) {
            alert('Error: ' + error.message);
        }
    });
});
