<script>
export default {
	data: function () {
		return {
			errormsg: null,
			loading: false,
			username: null,
		};
	},
	methods: {
		async submitUsername() {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.post(
					"http://localhost:3000/session",
					{
						name: this.username,
					}
				);
				// Salvo lo userId in sessionStorage
				sessionStorage.setItem("userId", response.data.identifier);

				// Cambio schermata e vado alla HomeView
				this.$router.push("/Home");
			} catch (e) {
				if (e.response) {
					// Il server ha risposto con un codice di stato
					switch (e.response.status) {
						case 401:
							this.errormsg =
								"Errore 400: Lo username inserito non è valido.";
							break;
						case 500:
							this.errormsg =
								"Errore 500: Errore interno del server!";
							break;
						default:
							this.errormsg = `Errore ${e.response.status}: ${
								e.response.data.message ||
								"Si è verificato un errore."
							}`;
					}
				} else if (e.request) {
					// Nessuna risposta dal server
					this.errormsg =
						"Nessuna risposta dal server. Verifica la tua connessione.";
				}
			}
			this.loading = false;
		},
	},
};
</script>

<template>
	<div class="container">
		<h1 class="title">WASA Text</h1>
		<div class="login-box">
			<h1 class="login-label">Login</h1>
			<br /><br />
			<div>
				<form class="formField">
					<input v-model="username" required="" type="text" />
					<span>Username</span>
				</form>
			</div>
			<ErrorMsg v-if="errormsg" :msg="errormsg" class="error"></ErrorMsg>
			<button
				class="animated-button"
				@click="submitUsername"
				:disabled="loading"
			>
				<svg
					viewBox="0 0 24 24"
					class="arr-2"
					xmlns="http://www.w3.org/2000/svg"
				>
					<path
						d="M16.1716 10.9999L10.8076 5.63589L12.2218 4.22168L20 11.9999L12.2218 19.778L10.8076 18.3638L16.1716 12.9999H4V10.9999H16.1716Z"
					></path>
				</svg>
				<span class="text"
					>{{ loading ? "Loading..." : "SIGN IN" }}
				</span>
				<span class="circle"></span>
				<svg
					viewBox="0 0 24 24"
					class="arr-1"
					xmlns="http://www.w3.org/2000/svg"
				>
					<path
						d="M16.1716 10.9999L10.8076 5.63589L12.2218 4.22168L20 11.9999L12.2218 19.778L10.8076 18.3638L16.1716 12.9999H4V10.9999H16.1716Z"
					></path>
				</svg>
			</button>
			<br /><br />
		</div>
	</div>
</template>

<style scoped>
.container {
	display: flex;
	flex-direction: column;
	justify-content: center;
	align-items: center;
	height: 100vh;
	background: repeating-linear-gradient(
		45deg,
		#008080,
		#008080 10px,
		#20b2aa 10px,
		#20b2aa 20px
	);
}

.title {
	position: relative;
	text-align: center; /* Centra orizzontalmente il testo */
	/*top: 20px; /* Distanza dalla parte superiore */
	font-size: 100px; /* Dimensione del testo */
	font-weight: bold; /* Testo in grassetto */
	color: white; /* Colore del testo */
	margin-bottom: 20px;
}

.login-box {
	background-color: #20b2aa;
	padding: 30px;
	border-radius: 10px;
	box-shadow: 0px 4px 6px rgba(0, 0, 0, 0.1);
	text-align: center;
	width: 300px;
}

.login-label {
	font-size: 44px;
	margin-bottom: 20px;
	color: white;
}

.animated-button {
	position: relative;
	display: flex;
	align-items: center;
	gap: 4px;
	padding: 16px 36px;
	border: 4px solid;
	border-color: transparent;
	font-size: 16px;
	background-color: inherit;
	border-radius: 100px;
	font-weight: 600;
	color: rgb(47, 248, 255);
	box-shadow: 0 0 0 2px rgb(47, 255, 238);
	cursor: pointer;
	overflow: hidden;
	transition: all 0.6s cubic-bezier(0.23, 1, 0.32, 1);
}

.animated-button svg {
	position: absolute;
	width: 24px;
	fill: rgb(47, 255, 255);
	z-index: 9;
	transition: all 0.8s cubic-bezier(0.23, 1, 0.32, 1);
}

.animated-button .arr-1 {
	right: 16px;
}

.animated-button .arr-2 {
	left: -25%;
}

.animated-button .circle {
	position: absolute;
	top: 50%;
	left: 50%;
	transform: translate(-50%, -50%);
	width: 20px;
	height: 20px;
	background-color: rgb(47, 224, 255);
	border-radius: 50%;
	opacity: 0;
	transition: all 0.8s cubic-bezier(0.23, 1, 0.32, 1);
}

.animated-button .text {
	position: relative;
	z-index: 1;
	transform: translateX(-12px);
	transition: all 0.8s cubic-bezier(0.23, 1, 0.32, 1);
}

.animated-button:hover {
	box-shadow: 0 0 0 12px transparent;
	color: #212121;
	border-radius: 12px;
}

.animated-button:hover .arr-1 {
	right: -25%;
}

.animated-button:hover .arr-2 {
	left: 16px;
}

.animated-button:hover .text {
	transform: translateX(12px);
}

.animated-button:hover svg {
	fill: #212121;
}

.animated-button:active {
	scale: 0.95;
	box-shadow: 0 0 0 4px greenyellow;
}

.animated-button:hover .circle {
	width: 220px;
	height: 220px;
	opacity: 1;
}

.formField {
	margin: 10px;
	position: relative;
}

.formField input {
	padding: 10px 15px;
	outline: none;
	border: none;
	border-radius: 5px;
	background-color: #f1f1f1;
	color: #333;
	font-size: 16px;
	font-weight: 550;
	transition: 0.3s ease-in-out;
	box-shadow: 0 0 0 5px transparent;
}

.formField input:hover,
.formField input:focus {
	box-shadow: 0 0 0 2px #333;
}

.formField span {
	position: absolute;
	left: 0;
	top: 0;
	padding: 8px 15px;
	color: #333;
	font-size: 16px;
	font-weight: 600;
	transition: 0.3s ease-in-out;
	pointer-events: none;
}

.formField input:focus + span,
.formField input:valid + span {
	transform: translateY(-32px) translateX(-5px) scale(0.95);
	transition: 0.3s ease-in-out;
}

.error {
	color: red;
	font-weight: bold;
}
</style>
