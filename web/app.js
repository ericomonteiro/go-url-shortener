document.addEventListener('DOMContentLoaded', () => {
    const urlInput = document.getElementById('urlInput');
    const shortenBtn = document.getElementById('shortenBtn');
    const resultDiv = document.getElementById('result');

    // Only attach event listener if elements exist
    if (shortenBtn) {
        shortenBtn.addEventListener('click', async () => {
                const url = urlInput.value.trim();

                if (!url) {
                    showError('Please enter a URL');
                    return;
                }

                try {
                    const response = await fetch('/v1/shortener', {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json',
                        },
                        body: JSON.stringify({url}),
                    });

                    const data = await response.json();

                    if (response.ok) {
                        showSuccess(data.short_url);
                        urlInput.value = '';
                    } else {
                        showError(data.error || 'Failed to shorten URL');
                    }
                } catch (error) {
                    showError('An error occurred while processing your request');
                }
            }
        );

        function showSuccess(shortUrl) {
            const shortUrlDiv = document.getElementById('shortUrl');
            const copyBtn = document.getElementById('copyBtn');

            shortUrlDiv.textContent = `Short URL: ${shortUrl}`;
            resultDiv.className = 'result success';
            resultDiv.style.display = 'flex';

            // Add copy functionality
            copyBtn.addEventListener('click', () => {
                navigator.clipboard.writeText(shortUrl)
                    .then(() => {
                        copyBtn.textContent = 'Copied!';
                        copyBtn.className = 'copy-btn success';
                        setTimeout(() => {
                            copyBtn.textContent = 'Copy URL';
                            copyBtn.className = 'copy-btn';
                        }, 2000);
                    })
                    .catch(err => {
                        copyBtn.textContent = 'Error';
                        copyBtn.className = 'copy-btn error';
                        setTimeout(() => {
                            copyBtn.textContent = 'Copy URL';
                            copyBtn.className = 'copy-btn';
                        }, 2000);
                    });
            });
        }

        function showError(message) {
            resultDiv.textContent = message;
            resultDiv.className = 'result error';
            resultDiv.style.display = 'block';
        }
    }
})
