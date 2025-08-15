document.addEventListener('DOMContentLoaded', function() {
    const errorMessageEl = document.getElementById('errorMessage');
    const signupForm = document.getElementById('signupForm');

    signupForm.addEventListener('submit', async function(e) {
        e.preventDefault();

        const username = document.getElementById('username').value;
        const email = document.getElementById('email').value;
        const password = document.getElementById('password').value;

        try {
            const response = await fetch('/register', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    username: username,
                    email: email,
                    password: password
                }),
            });

            const data = await response.json();

            if (response.ok) {
                alert('Регистрация прошла успешно!');
                window.location.href = '/login';
            } else {
                errorMessageEl.textContent = data.error || 'Ошибка регистрации';
            }
        } catch (err) {
            errorMessageEl.textContent = 'Ошибка соединения с сервером';
            console.error('Signup error:', err);
        }
    });
});