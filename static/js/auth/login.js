document.addEventListener('DOMContentLoaded', function() {
    const errorMessageEl = document.getElementById('errorMessage');
    const loginForm = document.getElementById('loginForm');

    loginForm.addEventListener('submit', async function(e) {
        e.preventDefault();

        const username = document.getElementById('username').value;
        const password = document.getElementById('password').value;

        try {
            const response = await fetch('/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    username: username,
                    password: password
                }),
            });

            const data = await response.json();

            if (response.ok) {
                localStorage.setItem('token', data.token);

                const form = document.createElement('form');
                form.method = 'GET';
                form.action = '/main';

                const tokenInput = document.createElement('input');
                tokenInput.type = 'hidden';
                tokenInput.name = 'token';
                tokenInput.value = data.token;

                form.appendChild(tokenInput);
                document.body.appendChild(form);
                form.submit();
            } else {
                errorMessageEl.textContent = data.error || 'Ошибка входа';
            }
        } catch (err) {
            errorMessageEl.textContent = 'Ошибка соединения с сервером';
            console.error('Login error:', err);
        }
    });
});