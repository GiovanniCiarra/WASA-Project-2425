<template>
	<div class="new-group-container">
		<!-- page title -->
		<div class="page-title">Crea un Nuovo Gruppo</div>
		<div class="main-content">
			<!-- Colonna di sinistra -->
			<div class="left-column">
				<div class="photo-wrapper">
					<label class="group-photo">
						<input
							type="file"
							accept="image/*"
							@change="handlePhotoUpload"
							hidden
						/>
						<!-- Immagine del Gruppo -->
						<img
							class="group-photo-image"
							:src="groupPhoto || defaultGroupImage"
							alt="Group"
						/>
					</label>
					<!--  Icona rimuovi -->
					<img
						v-if="groupPhoto"
						class="remove-icon"
						src="../assets/icons/close.png"
						alt="Rimuovi"
						@click.stop="removeGroupPhoto"
					/>
				</div>

				<!-- Nome del gruppo -->
				<input
					class="group-name"
					v-model="groupName"
					placeholder="Nome del gruppo"
				/>
			</div>

			<!-- colonna destra -->
			<div class="right-column">
				<input
					v-model="searchQuery"
					type="text"
					class="search-bar"
					placeholder="Cerca utenti..."
				/>

				<div class="user-list">
					<UserItem
						v-for="user in filteredUsers"
						:key="user.userId"
						:user="user"
						:selected="isSelected(user)"
						@click="toggleUser(user)"
					/>
				</div>

				<button
					class="create-btn"
					@click="createGroup"
					:disabled="!canCreate"
				>
					Crea Gruppo
				</button>
			</div>
		</div>
	</div>
</template>

<script>
import UserItem from "../components/UserItem.vue";
import defaultGroupImage from "../assets/images/default-group-100.png";

export default {
	components: { UserItem },

	data() {
		return {
			users: [], // Array degli utenti registrati
			selectedUsers: [], // Array degli utenti selezionati
			searchQuery: "",
			groupName: "",
			groupPhoto: null,
			defaultGroupImage,
			userId: Number(sessionStorage.getItem("userId")),
		};
	},

	computed: {
		// filtro gli utenti in base alla searchQuery e li ordino se selezionati
		filteredUsers() {
			return this.users
				.filter((u) => u.userId !== this.userId)
				.filter((u) =>
					u.username
						.toLowerCase()
						.includes(this.searchQuery.toLowerCase())
				)
				.sort((a, b) => {
					const aSel = this.isSelected(a);
					const bSel = this.isSelected(b);
					return aSel === bSel ? 0 : aSel ? -1 : 1;
				});
		},

		// Abilita/Disabilita il bottone di creazione del gruppo
		canCreate() {
			return (
				this.groupName.trim() &&
				this.groupName.trim().length >= 3 &&
				this.groupName.trim().length <= 16 &&
				this.selectedUsers.length >= 2
			);
		},
	},

	// Alla creazione, esegue il fetch degli utenti
	async created() {
		try {
			const res = await this.$axios.get("/users", {
				headers: { Authorization: `Bearer ${this.userId}` },
			});
			this.users = res.data;
		} catch (err) {
			console.error("Errore nel caricamento utenti", err);
			this.handleError(err);
		}
	},

	methods: {
		// Verifica se un utente sia selezionato o meno
		isSelected(user) {
			return this.selectedUsers.some((u) => u.userId === user.userId);
		},

		// selettore degli utenti
		toggleUser(user) {
			if (this.isSelected(user)) {
				this.selectedUsers = this.selectedUsers.filter(
					(u) => u.userId !== user.userId
				);
			} else {
				this.selectedUsers.push(user);
			}
		},

		// photo handler
		handlePhotoUpload(event) {
			const file = event.target.files[0];
			if (file) {
				const reader = new FileReader();
				reader.onload = () => {
					this.groupPhoto = reader.result;
				};
				reader.readAsDataURL(file);
			}
		},

		// Rimuove la foto selezionata
		removeGroupPhoto() {
			this.groupPhoto = null;
		},

		// Crea il gruppo, specificando nome, foto e membri
		async createGroup() {
			try {
				const res = await this.$axios.post(
					"/conversations/create",
					{
						name: this.groupName,
						photo: this.groupPhoto || "",
						members: this.selectedUsers.map((u) => u.userId),
					},
					{
						headers: { Authorization: `Bearer ${this.userId}` },
					}
				);

				const { conversationId } = res.data;
				this.$router.push({
					path: `/conversations/${conversationId}`,
					query: {
						type: "group",
						name: this.groupName,
						photo: this.groupPhoto,
					},
				});
			} catch (err) {
				console.error("Errore nella creazione del gruppo", err);
				this.handleError(err);
			}
		},

		// Gestore degli errori
		handleError(err) {
			if (!err.response) {
				this.$root.showToast(
					"Errore di rete: controlla la connessione.",
					"error"
				);
				return;
			}

			const status = err.response.status;
			const detail = err.response.data || "ERRORE";

			switch (status) {
				case 400:
					this.$root.showToast(`Errore 400:\n ${detail}!`, "error");
					break;
				case 401:
					this.$root.showToast(
						"Errore 401: Non autorizzato, effettua di nuovo il login",
						"error"
					);
					this.$router.push("/Login");
					break;
				case 500:
					this.$root.showToast(`Errore 500:\n ${detail}!`, "error");
					break;
				default:
					this.$root.showToast("Errore sconosciuto", "error");
			}
		},
	},
};
</script>

<style scoped>
.new-group-container {
	display: flex;
	flex-direction: column;
	gap: 20px;
	padding: 30px;
	height: 100%;
	box-sizing: border-box;
	background: radial-gradient(circle at top left, #73bfa1, #a1ded0, #cdeee1);
}

.page-title {
	font-size: 2.2rem;
	font-weight: 700;
	color: white;
	margin-bottom: 20px;
	text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
	font-family: "Poppins", sans-serif;
	text-align: center;
}

.main-content {
	display: flex;
	gap: 30px;
	flex-grow: 1;
	padding-top: 10px;
}

.left-column {
	width: 280px;
	height: 300px;
	display: flex;
	flex-direction: column;
	align-items: center;
	gap: 20px;
	background-color: rgba(255, 255, 255, 0.2);
	padding: 20px;
	border: 1px solid rgba(255, 255, 255, 0.3);
	border-radius: 20px;
	justify-content: center;
	box-shadow: 0 8px 32px rgba(0, 0, 0, 0.15);
	backdrop-filter: blur(10px);
	transition: transform 0.3s ease;
	margin-left: 10px;
}

.left-column:hover {
	transform: scale(1.02);
}

.right-column {
	flex-grow: 1;
	display: flex;
	flex-direction: column;
	gap: 15px;
}

.group-photo {
	position: relative;
	cursor: pointer;
}

.group-photo-image {
	width: 150px;
	height: 150px;
	border-radius: 50%;
	object-fit: cover;
	background-color: white;
	box-shadow: 0 4 12px rgba(0, 0, 0, 0.2);
	transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.group-photo:hover .group-photo-image {
	transform: scale(1.05);
	box-shadow: 0 6px 20px rgba(0, 0, 0, 0.3);
}

.photo-wrapper {
	position: relative;
	width: fit-content;
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

.group-name {
	width: 100%;
	max-width: 240px;
	padding: 16px 20px;
	font-size: 1.1rem;
	border: 1px solid #ddd;
	border-radius: 12px;
	background-color: rgba(255, 255, 255, 0.9);
	box-shadow: 0 1px 4px rgba(0, 0, 0, 0.1);
	transition: all 0.3s ease;
	text-align: center;
}

.group-name:focus {
	border-color: #73bfa1;
	box-shadow: 0 0 8px rgba(74, 157, 236, 0.4);
	outline: none;
}

.group-name::placeholder {
	color: #aaa;
	text-align: center;
}

.search-bar {
	padding: 10px;
	width: 100%;
	border-radius: 10px;
	border: 1px solid #ccc;
	font-size: 1rem;
	background-color: #f3f3f3;
	transition: all 0.3s ease;
}

.search-bar:hover,
.search-bar:focus {
	border: 2px solid #73bfa1;
	box-shadow: 0px 0px 0px 3px rgba(74, 101, 236, 0.2);
	background-color: white;
}

.user-list {
	flex-grow: 1;
	overflow-y: auto;
	display: flex;
	flex-direction: column;
	gap: 10px;
	max-height: 450px;
}

.create-btn {
	margin-top: auto;
	margin-bottom: 5px;
	padding: 14px 30px;
	background: linear-gradient(135deg, #2fa27a, #73bfa1);
	color: white;
	border: none;
	border-radius: 50px;
	font-weight: bold;
	font-size: 1.1rem;
	cursor: pointer;
	transition: all 0.3s ease;
	box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
}

.create-btn:hover {
	background: linear-gradient(135deg, #24996d, #5bb896);
	box-shadow: 0 6px 18px rgba(0, 0, 0, 0.3);
	transform: translateY(-2px);
}

.create-btn:disabled {
	background: #a5d6c5;
	cursor: not-allowed;
	box-shadow: none;
	transform: none;
}
</style>
