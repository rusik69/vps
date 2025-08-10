document.addEventListener('DOMContentLoaded', function() {
    const form = document.getElementById('shortenForm');
    const urlInput = document.getElementById('url');
    const customCheckbox = document.getElementById('custom');
    const customCodeField = document.getElementById('customCodeField');
    const customCodeInput = document.getElementById('customCode');
    const resultDiv = document.getElementById('result');
    const shortUrlInput = document.getElementById('shortUrl');
    const copyButton = document.getElementById('copyButton');

    // Toggle custom code field
    customCheckbox.addEventListener('change', function() {
        if (customCheckbox.checked) {
            customCodeField.classList.remove('hidden');
        } else {
            customCodeField.classList.add('hidden');
            customCodeInput.value = '';
        }
    });

    // Handle form submission
    form.addEventListener('submit', async function(e) {
        e.preventDefault();

        const url = urlInput.value;
        const requestBody = { url };

        // Add custom code if checkbox is checked
        if (customCheckbox.checked && customCodeInput.value.trim()) {
            requestBody.custom_code = customCodeInput.value.trim();
        }

        try {
            // Make API request
            const response = await fetch('/api/shorten', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(requestBody)
            });

            if (!response.ok) {
                throw new Error('Failed to shorten URL');
            }

            const data = await response.json();
            shortUrlInput.value = data.full_url;
            resultDiv.classList.remove('hidden');
            
            copyButton.addEventListener('click', () => {
                navigator.clipboard.writeText(data.full_url)
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
