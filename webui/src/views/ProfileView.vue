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
					<img
						class="edit-icon"
						src="../assets/icons/pen.png"
						alt="Modifica"
					/>
				</label>

				<img
					v-if="newPhoto || photo"
					class="remove-icon"
					src="../assets/icons/close.png"
					alt="Rimuovi"
					@click.stop="removeUserPhoto"
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
import defaultUserPhoto from "../assets/images/default-user.png";

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
					this.newUsername.trim() !== this.username) ||
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
					console.log("Username aggiornato con successo");
				}
			} catch (err) {
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
					console.log("Photo aggiornata con successo");
				}
			} catch (err) {
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
	background-color: #f0f2f5;
	overflow-y: auto;
	box-sizing: border-box;
}

/* Titolo alto */
.page-title {
	font-size: 2.2rem;
	font-weight: 700;
	color: white;
	margin-bottom: 30px;
	text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
	font-family: "Poppins", sans-serif;
}

/* Wrapper contenuto principale */
.profile-content {
	background-color: white;
	padding: 40px 20px;
	border-radius: 16px;
	box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
	display: flex;
	flex-direction: column;
	align-items: center;
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
	box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
	background-color: #e0e0e0;
}

.edit-icon {
	position: absolute;
	bottom: 0;
	right: 0;
	width: 34px;
	height: 34px;
	background: white;
	border-radius: 50%;
	padding: 6px;
	box-shadow: 0 0 5px rgba(0, 0, 0, 0.3);
	transform: translate(30%, 30%);
}

.username-input {
	width: 100%;
	max-width: 400px;
	padding: 14px 18px;
	font-size: 1.2rem;
	border: 1px solid #ccc;
	border-radius: 10px;
	background-color: white;
	margin-bottom: 30px;
	transition: border 0.3s;
}

.username-input:focus {
	border: 1px solid #4a9dec;
	outline: none;
}

.save-btn {
	margin-top: 20px;
	padding: 14px 28px;
	font-size: 1.1rem;
	border: none;
	border-radius: 10px;
	background-color: #2fa27a;
	color: white;
	cursor: pointer;
	transition: background-color 0.3s;
}

.save-btn:hover {
	background-color: #24996d;
}

.save-btn:disabled {
	background-color: #a5d6c5;
	cursor: not-allowed;
}
</style>
