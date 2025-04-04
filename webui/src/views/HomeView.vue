<script>
export default {
	data() {
		return {
			users: [],
			conversations: [],
			searchUser: "",
			searchConversation: "",
			userError: "",
			conversationError: "",
		};
	},

	//Ogni volta che cambia il valore di this.users o this.searchUser il valore ricercato viene ricalcolato automaticamente (stessa cosa per conversation)
	computed: {
		filteredUsers() {
			return this.users.filter((user) =>
				user.username
					.toLowerCase()
					.includes(this.searchUser.toLowerCase())
			);
		},
		filteredConversations() {
			return this.conversations.filter((conversation) =>
				conversation.name
					.toLowerCase()
					.includes(this.searchConversation.toLowerCase())
			);
		},
	},

	async created() {
		const userId = sessionStorage.getItem("userId"); //Recupero lo userId dal session storage
		if (!userId) {
			console.error("UserId non trovato nel sessionStorage");
			return;
		}
		try {
			let usersResponse = await this.$axios.get(
				"http://localhost:3000/users",
				{
					headers: {
						Authorization: `Bearer ${userId}`,
					},
				}
			);
			this.users = usersResponse.data; //Salvo i dati degli utenti
		} catch (e) {
			if (e.response) {
				// Il server ha risposto con un codice di stato
				switch (e.response.status) {
					case 401:
						this.userError = "Errore 401: Unauthorized!";
						break;
					case 500:
						this.userError =
							"Errore 500: Errore interno del server!";
						break;
					default:
						this.userError = `Errore ${e.response.status}: ${
							e.response.data.message ||
							"Si è verificato un errore!"
						}`;
				}
			} else {
				this.userError = "Si è verificato un errore!";
			}
		}

		try {
			let conversationsResponse = await this.$axios.get(
				"http://localhost:3000/conversations",
				{
					headers: {
						Authorization: `Bearer ${userId}`,
					},
				}
			);
			this.conversations = conversationsResponse.data; //Salvo i dati delle conversazioni
			/*console.log("Utenti ricevuti:", this.users);
			console.log("Conversazioni ricevute:", this.conversations);*/
		} catch (e) {
			if (e.response) {
				// Il server ha risposto con un codice di stato
				switch (e.response.status) {
					case 400:
						this.conversationError = "Errore 400: Bad Request!";
						break;
					case 401:
						this.conversationError = "Errore 401: Unauthorized!";
						break;
					case 500:
						this.conversationError =
							"Errore 500: Errore interno del server!";
						break;
					default:
						this.conversationError = `Errore ${
							e.response.status
						}: ${
							e.response.data.message ||
							"Si è verificato un errore!"
						}`;
				}
			} else {
				this.conversationError = "Si è verificato un errore!";
			}
		}
	},
};
</script>

<template>
	<div class="container">
		<!-- Colonna sinistra per la Lista utenti -->
		<div class="users-column">
			<input
				v-model="searchUser"
				type="text"
				placeholder="Search User..."
				class="search-bar"
			/>
			<div v-if="userError" class="error-message">{{ userError }}</div>
			<!--Se c'è uno userError lo mostro -->
			<div class="user-list">
				<div
					v-for="user in filteredUsers"
					:key="user.userId"
					class="user-item"
					@click="goToConversation(user.userId)"
				>
					<img
						:src="user.photo || '/images/profile.png'"
						alt="Avatar"
						class="user-avatar"
					/>
					<span class="user-name">{{ user.username }}</span>
				</div>
			</div>
		</div>

		<!-- Colonna destra per la Lista conversazioni -->
		<div class="conversation-column">
			<input
				v-model="searchConversation"
				type="text"
				placeholder="Search Conversation..."
				class="search-bar"
			/>
			<div v-if="conversationError" class="error-message">
				{{ conversationError }}
			</div>
			<!-- Se c'è un conversationError lo mostro -->
			<div class="conversation-list">
				<div
					v-for="conversation in filteredConversations"
					:key="conversation.conversationId"
					class="conversation-item"
				>
					<img
						:src="conversation.photo"
						alt="Conversation Photo"
						class="conversation-avatar"
					/>
					<div class="conversation-info">
						<span class="conversation-name">{{
							conversation.name
						}}</span>
						<span
							v-if="conversation.type === 'group'"
							class="conversation-sender"
							>{{
								conversation.lastMessage.SenderUsername
							}}:</span
						>
						<span class="conversation-snippet">{{
							conversation.lastMessage.Content
						}}</span>
					</div>
					<span class="conversation-status">{{
						conversation.status
					}}</span>
				</div>
			</div>
		</div>
	</div>
</template>

<style scoped>
.container {
	display: flex;
	width: 100vw;
	height: 100vh;
	background: repeating-linear-gradient(
		45deg,
		#008080,
		#008080 10px,
		#20b2aa 10px,
		#20b2aa 20px
	);
	padding: 20px;
	box-sizing: border-box;
	position: relative;
}

.users-column {
	width: 30%;
	background: #20b2aa;
	padding: 10px;
	border-radius: 10px;
	overflow-y: auto;
	border-right: 2px solid black;
}

.conversation-column {
	flex-grow: 1;
	background: #ffffff;
	padding: 10px;
	border-radius: 10px;
	overflow-y: auto;
	border-left: 2px solid black;
}

.search-bar {
	width: 100%;
	padding: 8px;
	padding-left: 0.8em;
	outline: none;
	overflow: hidden;
	margin-bottom: 10px;
	border: 2px solid transparent;
	border-radius: 10px;
	background-color: #f3f3f3;
	transition: all 0.5s;
}

.search-bar:hover,
.search-bar:focus {
	border: 2px solid #4a9dec;
	box-shadow: 0px 0px 0px 7px rgb(74, 157, 236, 20%);
	background-color: white;
}

.user-list,
.conversation-list {
	max-height: 80vh;
	overflow-y: auto;
}

.user-item,
.conversation-item {
	display: flex;
	align-items: center;
	padding: 10px;
	border-bottom: 1px solid #ddd;
}

.user-avatar,
.conversation-avatar {
	width: 40px;
	height: 40px;
	border-radius: 50%;
	margin-right: 10px;
}

.conversation-info {
	flex-grow: 1;
}

.settings-icon {
	position: absolute;
	top: 10px;
	right: 10px;
	font-size: 24px;
	cursor: pointer;
}

.error-message {
	color: red;
	font-weight: bold;
	text-align: center;
	margin-top: 10px;
	background: #ffe0e0;
	padding: 10px;
	border-radius: 5px;
}
</style>
