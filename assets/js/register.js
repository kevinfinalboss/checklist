
document.addEventListener("DOMContentLoaded", function() {
    const form = document.getElementById("register-form");

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
                console.log("User successfully registered:", data);
            } else {
                const data = await response.json();
                console.log("Error registering user:", data);
            }
        } catch (error) {
            console.log("There was a problem sending the request:", error);
        }
    });
});
