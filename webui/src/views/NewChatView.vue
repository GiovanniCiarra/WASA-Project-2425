<template>
	<div class="new-chat-container">
		<!-- page title -->
		<div class="page-title">Inizia una Nuova Chat</div>

		<!-- pannello principale -->
		<div class="new-chat-panel">
			<input
				v-model="searchQuery"
				class="search-bar"
				type="text"
				placeholder="Cerca utenti..."
			/>

			<div v-if="filteredUsers.length > 0" class="user-list">
				<UserItem
					v-for="user in filteredUsers"
					:key="user.userId"
					:user="user"
					@click="handleUserClick(user)"
				/>
			</div>

			<div v-else class="empty-state">Nessun utente trovato.</div>
		</div>
	</div>
</template>

<script>
import UserItem from "../components/UserItem.vue";

export default {
	components: { UserItem },
	data() {
		return {
			users: [], // array degli utenti registrati
			searchQuery: "",
			userId: null,
		};
	},
	computed: {
		filteredUsers() {
			return this.users.filter((u) =>
				u.username
					.toLowerCase()
					.includes(this.searchQuery.toLowerCase())
			);
		},
	},

	created() {
		this.fetchUsers();
	},

	methods: {
		// Fetch degli utenti
		async fetchUsers() {
			this.userId = sessionStorage.getItem("userId");
			if (!this.userId) return this.$router.push("/Login");

			try {
				const res = await this.$axios.get("/users", {
					headers: { Authorization: `Bearer ${this.userId}` },
				});
				this.users = res.data;
			} catch (err) {
				console.error("Errore nel recupero utenti", err);
			}
		},

		// handler del click su un utente, che crea una conversazione o la pusha se già esistente
		async handleUserClick(targetUser) {
			try {
				// Verifico se già esiste una conversazione con questo utente
				const res = await this.$axios.get(
					`/privateChat/${targetUser.userId}`,
					{
						headers: { Authorization: `Bearer ${this.userId}` },
					}
				);

				// Se la conversazione esiste, allora ottengo il conversationId e cambio View
				const conversationId = res.data.conversationId;
				this.$router.push({
					path: `/conversations/${conversationId}`,
					query: {
						type: "chat",
						name: targetUser.username,
						photo: targetUser.photo,
					},
				});
			} catch (err) {
				if (err.response && err.response.status === 404) {
					// Se la conversazione non esiste, ricevo un error 404, allora la creo
					try {
						const createRes = await this.$axios.post(
							"/conversations/create",
							{
								name: "",
								photo: "",
								members: [targetUser.userId],
							},
							{
								headers: {
									Authorization: `Bearer ${this.userId}`,
								},
							}
						);

						const conversationId = createRes.data.conversationId;
						this.$router.push({
							path: `/conversations/${conversationId}`,
							query: {
								type: "chat",
								name: targetUser.username,
								photo: targetUser.photo,
							},
						});
					} catch (err) {
						console.error(
							"Errore nella creazione della conversazione",
							err
						);
						this.handleError(err);
					}
				} else {
					console.error("Errore nel recupero conversazione", err);
					this.handleError(err);
				}
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
.new-chat-container {
	flex-grow: 1;
	display: flex;
	flex-direction: column;
	align-items: center;
	background: radial-gradient(circle at top left, #73bfa1, #a1ded0, #cdeee1);
	height: 100%;
	padding: 30px 20px;
	box-sizing: border-box;
}

.page-title {
	font-size: 2.2rem;
	font-weight: 700;
	color: white;
	margin-bottom: 20px;
	text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
	text-align: center;
	font-family: "Poppins", sans-serif;
}

.new-chat-panel {
	background-color: white;
	width: 100%;
	max-width: 900px;
	height: 75vh;
	border-radius: 16px;
	padding: 20px;
	box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
	display: flex;
	flex-direction: column;
	overflow: hidden;
}

.search-bar {
	border: 2px solid transparent;
	width: 100%;
	height: 2.5em;
	padding-left: 0.8em;
	outline: none;
	background-color: #f3f3f3;
	border-radius: 10px;
	margin-bottom: 15px;
	transition: all 0.5s;
}

.search-bar:hover,
.search-bar:focus {
	border: 2px solid #73bfa1;
	box-shadow: 0px 0px 0px 7px rgb(74, 157, 236, 20%);
	background-color: white;
}

.user-list {
	max-height: calc(85vh - 150px);
	flex-grow: 1;
	overflow-y: auto;
	display: flex;
	flex-direction: column;
	gap: 6px;
	scroll-behavior: smooth;
	scrollbar-width: thin;
	scrollbar-color: #ccc transparent;
	padding-right: 4px;
}

.empty-state {
	text-align: center;
	margin-top: 2rem;
}
</style>
