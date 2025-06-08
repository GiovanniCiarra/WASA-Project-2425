<template>
	<div class="login-wrapper">
		<div class="login-box">
			<h1 class="title">WASA Text</h1>
			<h2 class="subtitle">Accedi al tuo account</h2>
			<form class="login-form" @submit.prevent="submitUsername">
				<input
					v-model="username"
					required
					type="text"
					placeholder="Username"
					:class="{ 'error-input': hasError }"
				/>

				<div v-if="hasError" class="input-error-message">
					{{ errorMessage }}
				</div>
				<button>
					{{ "Entra" }}
				</button>
			</form>
		</div>
	</div>
</template>

<script>
export default {
	data() {
		return {
			username: null, // username inserito dall'utente
			hasError: false, // variabile per errore username
		};
	},

	methods: {
		async submitUsername() {
			try {
				const response = await this.$axios.post("/session", {
					name: this.username,
				});

				// Salvo lo userId in sessionStorage
				sessionStorage.setItem("userId", response.data.identifier);

				// Cambio schermata e vado alla HomeView
				if (response.status === 200) {
					this.$root.showToast(
						"Login effettuato con successo",
						"success"
					);
				}
				if (response.status === 201) {
					this.$root.showToast(
						"Utente creato con successo",
						"success"
					);
				}
				this.$router.push("/Home");
			} catch (err) {
				this.handleError(err);
			}
		},

		handleError(err) {
			if (!err.response) {
				this.$root.showToast(
					"Errore di rete: controlla la connessione.",
					"error"
				);
				return;
			}

			const status = err.response.status;

			switch (status) {
				case 400:
					this.hasError = true;
					this.errorMessage =
						"Lo username deve avere lunghezza di 3-16 caratteri e contenere almeno un carattere alfanumerico";
					const detail = err.response.data || "Richiesta non valida";
					this.$root.showToast(`Errore 400: ${detail}`, "error");
					setTimeout(() => {
						this.hasError = false;
						this.errorMessage = "";
					}, 3000);
					break;
				case 500:
					this.$root.showToast(
						"Errore 500:\n Internal Server Error!",
						"error"
					);
					break;
				default:
					this.$root.showToast("Errore sconosciuto", "error");
			}
		},
	},
};
</script>

<style scoped>
.login-wrapper {
	width: 100%;
	display: flex;
	justify-content: center;
	align-items: center;
	height: 100%;
}

.login-box {
	background: rgba(255, 255, 255, 0.1);
	backdrop-filter: blur(14px);
	border-radius: 20px;
	padding: 60px 50px;
	box-shadow: 0 8px 32px rgba(0, 0, 0, 0.4);
	text-align: center;
	width: 100%;
	max-width: 500px;
}

.title {
	font-size: 3.2rem;
	font-weight: bold;
	color: #ffffff;
	text-shadow: 0 4px 12px rgba(0, 0, 0, 0.7);
	margin-bottom: 12px;
}

.subtitle {
	font-size: 1.6rem;
	color: #ffffffcc;
	margin-bottom: 35px;
}

.login-form {
	display: flex;
	flex-direction: column;
	gap: 15px;
}

.login-form input {
	padding: 14px 20px;
	border: none;
	border-radius: 30px;
	background: rgba(255, 255, 255, 0.85);
	font-size: 1.1rem;
	color: #333;
	transition: box-shadow 0.3s ease;
}

.login-form input:focus {
	outline: none;
	box-shadow: 0 0 0 3px #2fa27a;
}

.login-form button {
	padding: 14px 20px;
	border: none;
	border-radius: 30px;
	background: #2fa27a;
	color: #fff;
	font-size: 1.1rem;
	font-weight: 600;
	cursor: pointer;
	transition: background 0.3s ease, transform 0.2s ease;
}

.login-form button:hover {
	background: #24996d;
	transform: translateY(-2px);
}

.login-form button:disabled {
	background: #88cbb8;
	cursor: not-allowed;
}

.error-input {
	border: 2px solid #ff4d4f;
	box-shadow: 0 0 0 3px rgba(255, 77, 79, 0.5);
	animation: shake 0.3s;
}

@keyframes shake {
	0% {
		transform: translateX(0);
	}
	25% {
		transform: translateX(-5px);
	}
	50% {
		transform: translateX(5px);
	}
	75% {
		transform: translateX(-5px);
	}
	100% {
		transform: translateX(0);
	}
}

.input-error-message {
	color: #ff4d4f;
	font-size: 0.85rem;
	margin-top: 4px;
	text-align: left;
	padding-left: 10px;
}
</style>
