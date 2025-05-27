<template>
	<div class="home-container">
		<div class="page-title">Le tue Conversazioni</div>
		<div class="conversation-panel">
			<!-- barra di ricerca -->
			<input
				v-model="searchQuery"
				type="text"
				class="search-bar"
				placeholder="Cerca conversazioni..."
			/>

			<!-- Lista conversazioni -->
			<div
				v-if="filteredConversations && filteredConversations.length > 0"
				class="conversation-list"
			>
				<ConversationItem
					v-for="conv in filteredConversations"
					:key="conv.conversationId"
					:conversation="conv"
					:currentUser="userId"
					@click="
						$router.push({
							path: `/conversations/${conv.conversationId}`,
							query: {
								type: conv.type,
								name: conv.name,
								photo: conv.photo,
							},
						})
					"
				/>
			</div>

			<!-- Stato vuoto -->
			<div v-else class="empty-state">
				<p>Nessuna conversazione trovata.</p>
				<button @click="startNewConversation">
					Avvia una nuova chat
				</button>
			</div>
		</div>
	</div>
</template>

<script>
import ConversationItem from "../components/ConversationItem.vue";

export default {
	components: { ConversationItem },

	data() {
		return {
			conversations: [],
			searchQuery: "",
			userId: null,
			pollingInterval: null,
		};
	},

	mounted() {
		this.userId = Number(sessionStorage.getItem("userId"));
		if (!this.userId) {
			this.$router.push("/Login");
		}
		this.fetchConversations();

		// Avvio il polling per il fetch delle conversazioni
		this.pollingInterval = setInterval(() => {
			this.fetchConversations();
		}, 5000);
	},

	beforeUnmount() {
		this.clearPolling();
	},

	beforeRouteLeave(_, __, next) {
		this.clearPolling();
		next();
	},

	computed: {
		filteredConversations() {
			if (this.conversations) {
				return this.conversations
					.filter(
						(conv) =>
							conv.name
								.toLowerCase()
								.includes(this.searchQuery.toLowerCase()) &&
							conv.lastMessage.messageId !== 0
					)
					.sort((a, b) => {
						const timeA = new Date(
							a.lastMessage.timestamp
						).getTime();
						const timeB = new Date(
							b.lastMessage.timestamp
						).getTime();
						return timeB - timeA; // Ordine decrescente, dal più recente
					});
			}
			return;
		},
	},

	methods: {
		// Push alla view per creazione chat
		startNewConversation() {
			this.$router.push("/NewChat");
		},

		// Fetch delle conversazioni
		async fetchConversations() {
			try {
				const res = await this.$axios.get("/conversations", {
					headers: {
						Authorization: `Bearer ${this.userId}`,
					},
				});
				this.conversations = res.data;
			} catch (err) {
				console.log("Errore nel caricamento delle conversazioni!");
				this.handleError(err);
			}
		},

		// Elimino il polling
		clearPolling() {
			if (this.pollingInterval) {
				clearInterval(this.pollingInterval);
				this.pollingInterval = null;
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

			switch (status) {
				case 401:
					this.$root.showToast(
						"Errore 401: Non autorizzato, effettua di nuovo il login",
						"error"
					);
					this.$router.push("/Login");
					break;
				case 500:
					this.$root.showToast(
						"Errore 500:\n Impossibile caricare le conversazioni!",
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
.home-container {
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
	font-family: "Poppins", sans-serif;
	text-align: center;
}

.conversation-panel {
	background-color: rgb(249, 249, 249);
	width: 100%;
	max-width: 900px;
	height: 75vh;
	border-radius: 16px;
	padding: 20px;
	padding-top: 10px;
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
	overflow: hidden;
	background-color: #f3f3f3;
	border-radius: 10px;
	transition: all 0.5s;
	margin-bottom: 12px;
	margin-top: 10px;
}

.search-bar:hover,
.search-bar:focus {
	border: 2px solid #73bfa1;
	box-shadow: 0px 0px 0px 7px rgb(74, 157, 236, 20%);
	background-color: white;
}

.conversation-list {
	flex-grow: 1;
	overflow-y: auto;
	display: flex;
	flex-direction: column;
	gap: 10px;
}

.empty-state {
	margin-top: 2rem;
	text-align: center;
}

.empty-state button {
	margin-top: 1rem;
	padding: 10px 16px;
	border: none;
	border-radius: 8px;
	background-color: #008080;
	color: white;
	cursor: pointer;
}

button {
	padding: 1.3em 3em;
	font-size: 12px;
	font: sans-serif;
	text-transform: uppercase;
	letter-spacing: 2.5px;
	font-weight: 500;
	color: #000;
	background-color: #fff;
	border: none;
	border-radius: 45px;
	box-shadow: 0px 8px 15px rgba(0, 0, 0, 0.1);
	transition: all 0.3s ease 0s;
	cursor: pointer;
	outline: none;
}

button:hover {
	background-color: #23c483;
	box-shadow: 0px 15px 20px rgba(46, 229, 157, 0.4);
	color: #fff;
	transform: translateY(-7px);
}

button:active {
	transform: translateY(-1px);
}
</style>
