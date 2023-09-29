
document.addEventListener("DOMContentLoaded", function() {
    const form = document.getElementById("register-form");

    form.addEventListener("submit", function(event) {
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

        // Aqui você pode adicionar o código para enviar esses dados ao seu servidor
        console.log("User data to be sent:", formData);
    });
});
