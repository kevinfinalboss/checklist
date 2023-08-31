const inputs = document.querySelectorAll(".input");

function addcl(){
	let parent = this.parentNode.parentNode;
	parent.classList.add("focus");
}

function remcl(){
	let parent = this.parentNode.parentNode;
	if(this.value == ""){
		parent.classList.remove("focus");
	}
}

inputs.forEach(input => {
	input.addEventListener("focus", addcl);
	input.addEventListener("blur", remcl);
});

window.onload = function() {
	const urlParams = new URLSearchParams(window.location.search);
	const sessionExpired = urlParams.get('session_expired');
	const invalidCredentials = urlParams.get('invalid_credentials');
	const loginSuccess = urlParams.get('login_success');

	if (sessionExpired === 'true') {
		showNotification('Sessão Expirada', 'red');
	} else if (invalidCredentials === 'true') {
		showNotification('Email ou Senha Inválida', 'red');
	} else if (loginSuccess === 'true') {
		showNotification('Login Realizado, Redirecionando...', 'green');
	}
};

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
