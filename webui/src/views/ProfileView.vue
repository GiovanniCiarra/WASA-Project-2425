<template>
	<div class="profile-container">
		<div class="page-title">Profilo</div>
		<div class="profile-content">
			<div class="profile-photo-wrapper">
				<label class="photo-upload">
					<input
						type="file"
						accept="image/*"
						@change="handlePhotoUpload"
						hidden
					/>
					<img
						:src="newPhoto || photo || defaultPhoto"
						alt="Foto profilo"
						class="profile-photo"
					/>
				</label>

				<img
					v-if="newPhoto || photo"
					class="remove-icon"
					src="../assets/icons/close.png"
					alt="Rimuovi"
					@click="removeUserPhoto"
				/>
			</div>

			<input
				v-model="newUsername"
				class="username-input"
				:placeholder="username"
			/>

			<button class="save-btn" @click="saveProfile" :disabled="!canSave">
				Salva Modifiche
			</button>
		</div>
	</div>
</template>

<script>
import defaultUserPhoto from "../assets/images/default-user-100.png";

export default {
	data() {
		return {
			userId: Number(sessionStorage.getItem("userId")),
			username: "",
			photo: null,
			newUsername: "",
			newPhoto: null,
			defaultPhoto: defaultUserPhoto,
		};
	},
	computed: {
		canSave() {
			return (
				(this.newUsername.trim() &&
					this.newUsername.trim() !== this.username &&
					this.newUsername.trim().length >= 3 &&
					this.newUsername.trim().length <= 16) ||
				this.newPhoto
			);
		},
	},
	methods: {
		async fetchUserProfile() {
			try {
				const res = await this.$axios.get("/profile", {
					headers: { Authorization: `Bearer ${this.userId}` },
				});
				this.username = res.data.username;
				this.photo = res.data.photo;
			} catch (err) {
				console.error("Errore nel recupero profilo", err);
			}
		},

		handlePhotoUpload(event) {
			const file = event.target.files[0];
			if (file) {
				const reader = new FileReader();
				reader.onload = () => {
					this.newPhoto = reader.result;
				};
				reader.readAsDataURL(file);
			}
		},

		async saveProfile() {
			try {
				if (
					this.newUsername.trim() &&
					this.newUsername.trim() !== this.username
				) {
					await this.$axios.put(
						"/profile/username",
						{
							username: this.newUsername,
						},
						{
							headers: { Authorization: `Bearer ${this.userId}` },
						}
					);
					this.username = this.newUsername;
					this.newUsername = "";
					this.$root.showToast(
						"Profilo aggiornato con successo",
						"success"
					);
				}
			} catch (err) {
				this.handleError(err);
				console.error(
					"Errore durante il salvataggio dello username del profilo",
					err
				);
			}

			try {
				if (this.newPhoto) {
					await this.$axios.put(
						"/profile/photo",
						{
							photo: this.newPhoto,
						},
						{
							headers: { Authorization: `Bearer ${this.userId}` },
						}
					);
					this.photo = this.newPhoto;
					this.newPhoto = null;
					this.$root.showToast(
						"Profilo aggiornato con successo",
						"success"
					);
				}
			} catch (err) {
				this.handleError(err);
				console.error(
					"Errore durante il salvataggio della foto del profilo",
					err
				);
			}
		},

		removeUserPhoto() {
			this.photo = null;
			this.newPhoto = null;
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
					const detail = err.response.data || "Richiesta non valida";
					this.$root.showToast(`Errore 400: ${detail}`, "error");
					break;
				case 401:
					this.$root.showToast(
						"Errore 401: Non autorizzato, effettua di nuovo il login",
						"error"
					);
					this.$router.push("/Login");
					break;
				case 500:
					this.$root.showToast(
						"Errore 500:\n Internal Server Error!",
						"error"
					);
					break;
				default:
					this.$root.showToast("Errore", "error");
			}
		},
	},

	created() {
		this.fetchUserProfile();
	},
};
</script>

<style scoped>
.profile-container {
	display: flex;
	flex-direction: column;
	padding: 30px;
	width: 100%;
	height: 100%;
	background: linear-gradient(30deg, #73bfa1, #cdeee1, #73bfa1);
	overflow-y: auto;
	box-sizing: border-box;
}

/* Titolo alto */
.page-title {
	align-self: center;
	font-size: 3.2rem;
	font-weight: bold;
	color: white;
	margin-bottom: 20px;
	text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
	font-family: "Poppins", sans-serif;
}

/* Wrapper contenuto principale */
.profile-content {
	background: rgba(255, 255, 255, 0.2);
	backdrop-filter: blur(10px);
	-webkit-backdrop-filter: blur(10px);
	border: 1px solid rgba(255, 255, 255, 0.3);
	padding: 40px 20px;
	border-radius: 20px;
	box-shadow: 0 8px 32px rgba(0, 0, 0, 0.15);
	display: flex;
	flex-direction: column;
	align-items: center;
	width: 50%;
	left: 25%;
	position: relative;
	transition: transform 0.3s ease;
}

.profile-content:hover {
	transform: scale(1.02);
}

.remove-icon {
	position: absolute;
	top: -12px;
	right: -12px;
	width: 22px;
	height: 22px;
	object-fit: contain;
	background: white;
	border-radius: 50%;
	padding: 3px;
	box-shadow: 0 0 3px rgba(0, 0, 0, 0.2);
	cursor: pointer;
	z-index: 3;
	opacity: 0.6;
}

.remove-icone:hover {
	opacity: 1;
}

.profile-photo-wrapper {
	position: relative;
	margin-bottom: 30px;
}

.photo-upload {
	position: relative;
	cursor: pointer;
	display: inline-block;
}

.profile-photo {
	width: 180px;
	height: 180px;
	border-radius: 50%;
	object-fit: cover;
	box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
	background-color: #ffffff;
	transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.photo-upload:hover .profile-photo {
	transform: scale(1.05);
	box-shadow: 0 6px 20px rgba(0, 0, 0, 0.3);
}

.username-input {
	width: 100%;
	max-width: 400px;
	padding: 16px 20px;
	font-size: 1.1rem;
	font-weight: 500;
	border: 1px solid #ddd;
	border-radius: 12px;
	background-color: rgba(255, 255, 255, 0.9);
	margin-bottom: 30px;
	box-shadow: 0 1px 4px rgba(0, 0, 0, 0.1);
	transition: all 0.3s ease;
	text-align: center;
}

.username-input:focus {
	border-color: #73bfa1;
	box-shadow: 0 0 8px rgba(74, 157, 236, 0.4);
	outline: none;
}

.username-input::placeholder {
	color: #2a2a2a;
	text-align: center;
}

.save-btn {
	margin-top: 20px;
	padding: 14px 30px;
	font-size: 1.1rem;
	border: none;
	border-radius: 50px;
	background: linear-gradient(135deg, #2fa27a, #73bfa1);
	color: white;
	cursor: pointer;
	transition: all 0.3s ease;
	box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
}

.save-btn:hover {
	background: linear-gradient(135deg, #24996d, #5bb896);
	box-shadow: 0 6px 18px rgba(0, 0, 0, 0.3);
	transform: translateY(-2px);
}

.save-btn:disabled {
	background: #a5d6c5;
	cursor: not-allowed;
	box-shadow: none;
	transform: none;
}
</style>
