document.addEventListener("DOMContentLoaded", function() {
    const form = document.getElementById("register-form");

    function showNotification(message, color) {
        const notification = document.getElementById('notification');
        notification.textContent = message;
        notification.style.backgroundColor = color;
        notification.style.animation = 'fadein 3s, fadeout 1s 3s';
        notification.style.display = 'block';
        setTimeout(() => {
            notification.style.display = 'none';
        }, 4000);
    }

    function redirectToLogin() {
        setTimeout(() => {
            window.location.href = '/login';
        }, 3000);
    }

    form.addEventListener("submit", async function(event) {
        event.preventDefault();
        const formData = {
            "Name": form["Name"].value,
            "Email": form["Email"].value,
            "Password": form["Password"].value,
            "ConfirmPassword": form["ConfirmPassword"].value,
            "CPF": form["CPF"].value,
            "BirthDate": form["BirthDate"].value,
            "TelephoneNumber": form["TelephoneNumber"].value,
            "Address": form["Address"].value
        };

        try {
            const response = await fetch('/user/register', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(formData)
            });

            if (response.ok) {
                const data = await response.json();
                showNotification('Usuário criado com sucesso', 'green');
                redirectToLogin();
            } else {
                const data = await response.json();
                showNotification(data.error, 'red');
            }
        } catch (error) {
            showNotification('Houve um problema ao enviar a solicitação', 'red');
        }
    });
});
