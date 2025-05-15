<template>
	<div class="conversation-view">
		<!-- Lista conversazioni a sinistra -->
		<div class="sidebar-conversations">
			<input
				v-model="searchQuery"
				type="text"
				placeholder="Cerca conversazioni..."
				class="search-bar"
			/>
			<div class="scrollable">
				<div
					v-for="conv in filteredConversations"
					:key="conv.conversationId"
				>
					<ConversationItem
						:conversation="conv"
						:currentUser="userId"
						:isSelected="conv.conversationId == conversationId"
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
			</div>
		</div>

		<!-- Finestra messaggi a destra -->
		<ConversationWindow
			ref="conversationWindow"
			:conversationId="conversationId"
			:members="members"
			:messages="messages"
			:conversationName="conversationName"
			:conversationPhoto="conversationPhoto"
			:conversationType="conversationType"
			:commentPopupMessage="commentPopupMessage"
			:conversationInfoOpen="conversationInfoOpen"
			@updateMessages="updateMessages"
			@updateGroupName="updateGroupName"
			@updateGroupPhoto="updateGroupPhoto"
			@updateGroupMembers="updateGroupMembers"
			@addComment="addComment"
			@deleteComment="deleteComment"
			@openCommentPopup="openCommentPopup"
			@closeCommentPopup="closeCommentPopup"
			@openCloseConversationInfo="
				conversationInfoOpen = !conversationInfoOpen
			"
			@messageForwarded="messageForwarded"
		/>
	</div>
</template>

<script>
import ConversationItem from "../components/ConversationItem.vue";
import ConversationWindow from "../components/ConversationWindow.vue";

export default {
	components: { ConversationItem, ConversationWindow },

	data() {
		return {
			messages: [], // Lista dei messaggi
			oldMessages: [], // Copia dei messaggi pre fetchMessages()
			members: [], // Conversation Members se type = 'group'
			conversations: [], // Lista delle conversazioni
			userId: null, // userId dell'utente corrente
			messageSent: false, // Booleano per notificare che la conversazione corrente abbia almeno un messaggio
			searchQuery: "", // Input della search bar delle conversazioni
			pollingInterval: null, // Variabile per il polling
			conversationName: this.$route.query.name, // ConversationName
			conversationPhoto: this.$route.query.photo, // ConversationPhoto
			conversationType: this.$route.query.type, // ConversationType
			commentPopupMessage: null, // Messaggio che segnala l'apertura del commentPopup associato
			conversationInfoOpen: false, // Booleano che segnala l'apertura degli overlay
		};
	},
	computed: {
		// Recupero il conversationId dai route params
		conversationId() {
			return Number(this.$route.params.conversationId);
		},

		//Filtro per le conversazioni visibili ordinate in base al timestamp dell'ultimo messaggio
		filteredConversations() {
			if (this.conversations) {
				return this.conversations
					.filter((conv) =>
						conv.name
							.toLowerCase()
							.includes(this.searchQuery.toLowerCase())
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
		},
	},

	watch: {
		// Se il conversationId cambia, eseguo un refetch e aggiorno le variabili nella query
		"$route.params.conversationId": {
			async handler(newId, oldId) {
				if (oldId && oldId !== newId) {
					this.conversationName = this.$route.query.name;
					this.conversationPhoto = this.$route.query.photo;
					this.conversationType = this.$route.query.type;
					this.conversationInfoOpen = false;
					await this.deleteEmptyConversation(oldId);
					this.messageSent = false;
					if (this.conversationType === "group") {
						this.fetchMembers();
					}
					this.fetchConversations();
					await this.fetchMessages();

					// Dopo il fetchMessage iniziale faccio uno scroll in basso
					this.$nextTick(() => {
						if (
							this.$refs.conversationWindow &&
							this.$refs.conversationWindow.scrollToBottom
						) {
							this.$refs.conversationWindow.scrollToBottom();
						}
					});
				}
			},
		},
	},

	mounted() {
		// Eseguo il fetch di conversazioni, messaggi e memebri
		this.userId = Number(sessionStorage.getItem("userId"));
		if (!this.userId) {
			this.$router.push("/Login");
		}
		this.fetchConversations();
		if (this.conversationType === "group") {
			this.fetchMembers();
		}
		this.fetchMessages().then(() => {
			this.$nextTick(() => {
				if (
					this.$refs.conversationWindow &&
					this.$refs.conversationWindow.scrollToBottom
				) {
					this.$refs.conversationWindow.scrollToBottom();
				}
			});
		});

		// Avvio il polling per il fetch dei membri, conversazioni, messaggi e controllo dell'esistenza
		this.pollingInterval = setInterval(async () => {
			if (this.conversationType === "group") {
				await this.fetchMembers();
			}
			await this.fetchMessages();
			await this.fetchConversations();
			this.checkConversationInfoUpdates(this.conversations);
			this.checkConversationExistance();
		}, 5000);
	},

	beforeUnmount() {
		this.clearPolling();
	},

	beforeRouteLeave(_, __, next) {
		this.clearPolling();
		this.deleteEmptyConversation(this.conversationId).finally(() => next());
	},

	methods: {
		/* -------------- FETCH -------------- */
		async fetchConversations() {
			try {
				const res = await this.$axios.get("/conversations", {
					headers: {
						Authorization: `Bearer ${this.userId}`,
					},
				});
				this.conversations = res.data;
			} catch (err) {
				console.error("Errore nel caricamento conversazioni", err);
				this.handleError(err);
			}
		},

		async fetchMessages() {
			try {
				const res = await this.$axios.get(
					`/conversation-info/${this.conversationId}`,
					{
						headers: {
							Authorization: `Bearer ${this.userId}`,
						},
					}
				);
				if (res.data) {
					this.oldMessages = this.messages;
					this.messages = res.data;
					this.messageSent = true;
				}

				if (this.messages.length > this.oldMessages.length) {
					if (
						this.$refs.conversationWindow &&
						this.$refs.conversationWindow.scrollToBottom
					) {
						this.$refs.conversationWindow.scrollToBottom();
					}
				}
			} catch (err) {
				console.error("Errore nel caricamento dei messaggi", err);
				this.handleError(err);
			}
		},

		async fetchMembers() {
			try {
				const res = await this.$axios.get(
					`/conversation-members/${this.conversationId}`,
					{
						headers: {
							Authorization: `Bearer ${this.userId}`,
						},
					}
				);
				if (res.data) {
					this.members = res.data;
				}
			} catch (err) {
				console.error(
					"errore nel caricamento dei membri della conversazione",
					err
				);
				this.handleError(err);
			}
		},

		/* -------------- UTILITY -------------- */

		// Controllo periodico se una conversazione è stata eliminata a causa di azioni di altri utenti
		checkConversationExistance() {
			if (
				!this.conversationId ||
				isNaN(this.conversationId) ||
				this.conversationId === 0
			) {
				return;
			}

			if (!this.conversations || this.conversations.length === 0) {
				this.clearPolling();
				this.$root.showToast(
					"Conversazione eliminata perché senza messaggi!",
					"warning"
				);
				this.$router.push("/Home");
				return;
			}

			const conv = this.conversations.find(
				(c) => c.conversationId === this.conversationId
			);
			if (!conv && this.$route.name === "Conversation") {
				this.clearPolling();
				this.$root.showToast(
					"Conversazione eliminata perché senza messaggi!",
					"warning"
				);
				this.$router.push("/Home");
			}
		},

		// Controllo periodico se l'altro utente (o membro del gruppo) ha modificato le info (name e photo)
		checkConversationInfoUpdates(conversations) {
			if (!conversations || conversations.length === 0) return;
			for (const conv of conversations) {
				if (conv.conversationId === this.conversationId) {
					if (conv.name != this.conversationName)
						this.conversationName = conv.name;
					if (conv.photo != this.conversationPhoto)
						this.conversationPhoto = conv.photo;
				}
			}
		},

		// Elimino le conversazioni senza messaggi
		async deleteEmptyConversation(conversationId) {
			if (!this.messageSent) {
				try {
					await this.$axios.delete(
						`/conversation-delete/${conversationId}`,
						{
							headers: { Authorization: `Bearer ${this.userId}` },
						}
					);
					this.$root.showToast(
						"Conversazione eliminata perché senza messaggi!",
						"warning"
					);
				} catch (err) {
					console.error(
						"Errore nella cancellazione conversazione",
						err
					);
					this.handleError(err);
				}
			}
		},

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

		/* -------------- GESTIONE MESSAGGI -------------- */

		/* Aggiorno localmente la lista dei messaggi della conversazione (add or delete) */
		updateMessages(message, flag) {
			if (flag === "add") {
				this.messageSent = true;
				this.messages.push(message);
				if (
					this.$refs.conversationWindow &&
					this.$refs.conversationWindow.scrollToBottom
				) {
					this.$refs.conversationWindow.scrollToBottom();
				}
			} else if (flag === "delete") {
				this.messages = this.messages.filter(
					(m) => m.messageId !== message.messageId
				);

				// Se il messaggio eliminato era l'ultimo della conversazione, torno alla Home
				if (this.messages.length === 0) {
					this.$router.push("/Home");
					this.$root.showToast(
						"Conversazione eliminata perchè vuota!",
						"warning"
					);
				}
			}
			this.fetchConversations();
		},

		/* -------------- GESTIONE GRUPPI -------------- */

		/* Aggiorno localmente il nome di un gruppo */
		updateGroupName(newName) {
			this.conversationName = newName;
			this.$router.replace({
				query: {
					...this.$route.query,
					name: newName,
				},
			});

			// Refetch per update immediato
			this.fetchConversations();
		},

		/* Aggiorno localmente la foto di un gruppo */
		updateGroupPhoto(newPhoto) {
			this.conversationPhoto = newPhoto;
			this.$router.replace({
				query: {
					...this.$route.query,
					photo: newPhoto,
				},
			});

			// Refetch per update immediato
			this.fetchConversations();
		},

		/* Aggiorno localmente i membri di un gruppo */
		updateGroupMembers(newMember) {
			this.members.push(newMember);
		},

		/* -------------- GESTIONE COMMENTI -------------- */

		/* Aggiunge un commento localmente */
		addComment({ messageId, newComment }) {
			const msg = this.messages.find((m) => m.messageId === messageId);
			if (msg) {
				msg.comments = msg.comments || [];
				msg.comments.push(newComment);
			}
		},

		/* Rimuove un commento localmente */
		deleteComment(comment) {
			this.commentPopupMessage.comments =
				this.commentPopupMessage.comments.filter(
					(c) => c.commentId !== comment.commentId
				);
			this.$root.showToast("Commento eliminato con successo", "success");
			if (this.commentPopupMessage.comments.length === 0) {
				this.closeCommentPopup();
			}
			this.fetchMessages();
		},

		/* Apre il popup contenente i commenti del messaggio */
		openCommentPopup(comment) {
			const message = this.messages.find(
				(m) => m.messageId === comment.messageId
			);

			if (message) {
				this.commentPopupMessage = message;
			}
		},

		/* Chiude il popup contenente i commenti del messaggio */
		closeCommentPopup() {
			this.commentPopupMessage = null;
		},

		/* -------------- GESTIONE FORWARD -------------- */

		/* Refetch e scroll quando viene effettuato il forward, per evitare ritardi */
		messageForwarded() {
			this.fetchMessages();
			setTimeout(() => {
				if (
					this.$refs.conversationWindow &&
					this.$refs.conversationWindow.scrollToBottom
				) {
					this.$refs.conversationWindow.scrollToBottom();
				}
			}, 100);
		},
	},
};
</script>

<style scoped>
.conversation-view {
	display: flex;
	height: 100%;
	width: 100%;
	background-color: #ece5dd;
}

.search-bar {
	border: 2px solid transparent;
	width: 100%;
	height: 2.5em;
	padding-left: 0.8em;
	outline: none;
	background-color: #f3f3f3;
	border-radius: 10px;
	transition: all 0.5s;
	margin-top: 15px;
	margin-bottom: 10px;
	box-sizing: border-box;
}

.search-bar:hover,
.search-bar:focus {
	border: 2px solid #73bfa1;
	background-color: white;
}

.sidebar-conversations {
	width: 340px;
	background-color: white;
	border-right: 1px solid #ddd;
	overflow-y: auto;
	padding: 0 10px;
}

.scrollable {
	padding: 5px;
	display: flex;
	flex-direction: column;
	gap: 8px;
}
</style>
