<template>
	<div class="container">
		<!-- Colonna sinistra per gli utenti, utilizza il componente UserList -->
		<div class="users-column">
			<UserList
				:users="users"
				:userError="userError"
				@userSelected="goToConversation"
			/>
		</div>
		<!-- Colonna destra per le conversazioni, utilizza il componente ConversationList -->
		<div class="conversation-column">
			<ConversationList
				:conversations="conversations"
				:conversationError="conversationError"
				@conversationSelected="openConversation"
			/>
		</div>
	</div>
</template>

<script>
import UserList from "../components/UsersList.vue";
import ConversationList from "../components/ConversationsList.vue";

export default {
	components: { UserList, ConversationList },
	data() {
		return {
			users: [],
			conversations: [],
			userError: "",
			conversationError: "",
			userId: null,
		};
	},
	methods: {
		async goToConversation(targetUserId) {
			try {
				console.log("debug");
				// Controllo se esiste già una conversazione con questo utente
				const response = await this.$axios.get(
					`http://localhost:3000/privateChat/${targetUserId}`,
					{
						headers: {
							Authorization: `Bearer ${this.userId}`,
						},
					}
				);

				// Se esiste, Navigo direttamente alla chat
				console.log("debug2");
				if (response.data && response.data.conversationId) {
					this.$router.push(
						`/conversation/id/${response.data.conversationId}`
					);
					return;
				}
			} catch (e) {
				console.log("debug3");
				// Se è Errore 404, allora la conversazione non esiste e la creo
				if (e.response && e.response.status === 404) {
					const message = prompt("Scrivi il primo messaggio");
					if (!message || message.trim() === "") return;

					let createConv = null;
					try {
						// Creo la nuova conversazione
						createConv = await this.$axios.post(
							"http://localhost:3000/conversations",
							{
								name: "",
								photo: "",
								members: [targetUserId],
							},
							{
								headers: {
									Authorization: `Bearer ${this.userId}`,
								},
							}
						);
					} catch (e) {
						console.log("debug4");
						this.handleError(e);
					}
					console.log("debug5");
					const convId = createConv.data.conversationId;

					//Invio il primo messaggio della conversazione
					try {
						console.log("debug6");
						const firstmessage = await this.$axios.post(
							`http://localhost:3000/conversations/${convId}/messages`,
							{
								content: "message",
								type: "text",
							},
							{
								headers: {
									Authorization: `Bearer ${this.userId}`,
								},
							}
						);
					} catch (e) {
						console.log("debug7");
						this.handleError(e);
					}

					// Navigo alla nuova conversazione
					this.$router.push(`/Conversations/id/${convId}`);
				} else {
					console.log("debug8");
					this.handleError(e);
				}
			}
		},

		async openConversation(conversationId) {
			// Navigo alla conversazione
			this.$router.push(`/Conversations/id/${conversationId}`);
		},

		handleError(e) {
			console.log("debug9");
			if (e.response) {
				// Il server risponde con un codice di stato
				switch (e.response.status) {
					case 400:
						this.userError = "Errore 400: Bad Request!";
						break;
					case 401:
						this.userError = "Errore 401: Unauthorized!";
						break;
					case 403:
						this.userError = "Errore 403: Forbidden!";
						break;
					case 500:
						this.userError = "Errore 500: Internal Server Error!";
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
		},
	},

	async created() {
		this.userId = sessionStorage.getItem("userId"); //Recupero lo userId dal session storage
		if (!this.userId) {
			console.error("userId non trovato nel sessionStorage");
			return;
		}
		try {
			//Fetch degli utenti
			let usersResponse = await this.$axios.get(
				"http://localhost:3000/users",
				{
					headers: {
						Authorization: `Bearer ${this.userId}`,
					},
				}
			);
			this.users = usersResponse.data; //Salvo i dati degli utenti
		} catch (e) {
			handleError(e);
		}

		try {
			//Fetch delle conversazioni
			let conversationsResponse = await this.$axios.get(
				"http://localhost:3000/conversations",
				{
					headers: {
						Authorization: `Bearer ${this.userId}`,
					},
				}
			);
			this.conversations = conversationsResponse.data; //Salvo i dati delle conversazioni
		} catch (e) {
			if (e.response) {
				// Il server risponde con un codice di stato
				switch (e.response.status) {
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

<style scooped>
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
</style>
