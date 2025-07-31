document.addEventListener('DOMContentLoaded', () => {
    const form = document.getElementById('shortenForm');
    const resultDiv = document.getElementById('result');
    const errorDiv = document.getElementById('error');
    const captchaContainer = document.getElementById('captchaContainer');
    const captchaImage = document.getElementById('captchaImage');
    const captchaAnswer = document.getElementById('captchaAnswer');
    const shortenedUrl = document.getElementById('shortenedUrl');

    // Initialize captcha
    updateCaptcha();

    form.addEventListener('submit', async (e) => {
        e.preventDefault();
        
        const url = document.getElementById('url').value;
        const captcha = captchaAnswer.value;

        try {
            const response = await fetch('/api/shorten', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ url, captcha }),
            });

            const data = await response.json();

            if (!response.ok) {
                throw new Error(data.error || 'Failed to shorten URL');
            }

            // Update UI with success
            resultDiv.classList.remove('hidden');
            shortenedUrl.href = data.shortUrl;
            shortenedUrl.textContent = data.shortUrl;
            errorDiv.classList.add('hidden');
        } catch (error) {
            // Update UI with error
            errorDiv.textContent = error.message;
            errorDiv.classList.remove('hidden');
            resultDiv.classList.add('hidden');
        } finally {
            // Reset form and update captcha
            form.reset();
            updateCaptcha();
        }
    });

    function updateCaptcha() {
        fetch('/api/captcha')
            .then(response => response.json())
            .then(data => {
                captchaImage.innerHTML = `<img src="data:image/png;base64,${data.image}" alt="Captcha" onclick="updateCaptcha()">`;
                captchaAnswer.value = '';
                captchaAnswer.focus();
            })
            .catch(error => {
                console.error('Failed to get captcha:', error);
            });
    }
});
