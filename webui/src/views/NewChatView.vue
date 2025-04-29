<template>
	<div class="new-chat-container">
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
			users: [],
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
	methods: {
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

		async handleUserClick(targetUser) {
			try {
				// Verifico se già esiste una conversazione con questo utente
				const res = await this.$axios.get(
					`/privateChat/${targetUser.userId}`,
					{
						headers: { Authorization: `Bearer ${this.userId}` },
					}
				);

				// Se la conversazione esiste, allora ottendo il conversationId e cambio View
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
							"/conversations",
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
					}
				} else {
					console.error(
						"Errore nella creazione/recupero conversazione",
						err
					);
				}
			}
		},
	},
	created() {
		this.fetchUsers();
	},
};
</script>

<style scoped>
.new-chat-container {
	flex-grow: 1;
	display: flex;
	justify-content: center;
	align-items: center;
	background-color: #d9dbd5;
	height: 100%;
	padding: 20px;
	box-sizing: border-box;
}

.new-chat-panel {
	background-color: white;
	width: 100%;
	max-width: 900px;
	height: 85vh;
	border-radius: 12px;
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
	border: 2px solid #4a9dec;
	background-color: white;
}

.user-list {
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
