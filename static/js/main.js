document.addEventListener('DOMContentLoaded', function() {
    const logoutBtn = document.getElementById('logoutBtn');
    
    const token = localStorage.getItem('token');
    if (!token) {
        window.location.href = '/';
        return;
    }

   
    logoutBtn.addEventListener('click', function() {
        localStorage.removeItem('token');
        window.location.href = '/';
    });

   
});