<template>
	<div class="conversation-window">
		<div class="main-content">
			<!-- Intestazione conversazione -->
			<ConversationHeader
				:name="conversationName"
				:photo="conversationPhoto"
				:type="conversationType"
				@openInfo="$emit('openCloseConversationInfo')"
			/>

			<!-- Overlay info (chat o gruppo) -->
			<div v-if="conversationInfoOpen" class="info-overlay-container">
				<!-- Se è un gruppo -->
				<GroupInfoOverlay
					v-if="conversationType === 'group'"
					:groupName="conversationName"
					:groupPhoto="conversationPhoto"
					:members="members"
					:currentUser="userId"
					@closeOverlay="$emit('openCloseConversationInfo')"
					@leaveGroup="leaveGroup"
					@addMembers="addMembers"
					@updateGroupPhoto="updateGroupPhoto"
					@updateGroupName="updateGroupName"
				/>

				<!-- Se è una chat -->
				<ChatInfoOverlay
					v-else
					:chatName="conversationName"
					:chatPhoto="conversationPhoto"
					@closeOverlay="$emit('openCloseConversationInfo')"
				/>
			</div>

			<!-- Area messaggi -->
			<div class="message-area" ref="messageArea">
				<Message
					v-for="msg in messages"
					:key="msg.messageId"
					:message="msg"
					:currentUser="this.userId"
					@forward="handleForward"
					@reply="handleReply"
					@comment="openCommentPicker"
					@delete="handleDeleteMessage"
					@clickComment="openCommentPopup"
				/>

				<!-- POPUP -->
				<CommentPicker
					v-if="activeMessageForComment"
					:style="{
						top: pickerPosition.top + 'px',
						left: pickerPosition.left + 'px',
					}"
					@select-emoji="addComment"
					@close-emoji="activeMessageForComment = null"
				/>

				<CommentPopUp
					v-if="commentPopupMessage"
					:comments="commentPopupMessage.comments"
					:currentUser="this.userId"
					@close="closeCommentPopup"
					@deleteComment="handleDeleteComment"
				/>

				<ForwardPopUp
					v-if="forwardingMessage"
					:message="forwardingMessage"
					:userId="userId"
					@cancel="forwardingMessage = null"
					@forwarded="messageForwarded"
				/>

				<!-- Messaggio informativo se conversazione vuota -->
				<div v-if="messages.length === 0" class="no-messages">
					Nessun messaggio
				</div>
			</div>

			<!-- Confirm Dialog per DeleteMessage -->
			<ConfirmDialog
				v-if="confirmingMessage"
				:message="'Vuoi davvero eliminare questo messaggio?'"
				@confirm="confirmDeleteMessage"
				@cancel="cancelDeleteMessage"
			/>

			<!-- MessageInput -->
			<MessageInput
				:conversationId="conversationId"
				:replyingTo="replyingToMessage"
				@messageSent="handleMessageSent"
				@cancelReply="clearReply"
			/>
		</div>
	</div>
</template>

<script>
import ConversationHeader from "../components/ConversationHeader.vue";
import MessageInput from "../components/MessageInput.vue";
import Message from "../components/Message.vue";
import ConfirmDialog from "../components/ConfirmDialog.vue";
import CommentPicker from "../components/CommentPicker.vue";
import CommentPopUp from "../components/CommentPopUp.vue";
import ForwardPopUp from "../components/ForwardPopUp.vue";
import GroupInfoOverlay from "../components/GroupInfoOverlay.vue";
import ChatInfoOverlay from "../components/ChatInfoOverlay.vue";

export default {
	components: {
		ConversationHeader,
		MessageInput,
		Message,
		ConfirmDialog,
		CommentPicker,
		CommentPopUp,
		ForwardPopUp,
		GroupInfoOverlay,
		ChatInfoOverlay,
	},

	props: {
		conversationId: {
			type: [String, Number],
			required: true,
		},
		messages: Array, // Array dei messaggi
		members: Array, // Array dei membri
		conversationName: String, // Nome della conversazione
		conversationPhoto: String, // Foto della conversazione
		conversationType: String, // Tipo della conversazione
		commentPopupMessage: Object, // Messaggio relativo al commentPopup da aprire
		conversationInfoOpen: Boolean, // Boolean per aprire/chiudere gli overlay
	},

	data() {
		return {
			userId: null,
			confirmingMessage: null, // Messaggio da eliminare (apre ConfirmDialg)
			activeMessageForComment: null, // messaggio da commentare (Apre CommentPicker)
			pickerPosition: { top: 0, left: 0 }, // Decide la posizione del CommentPicker
			forwardingMessage: null, // Messaggio da inoltrare (Apre ForwardPopup)
			replyingToMessage: null, // Messaggio a cui rispondere (Modifica ConversationInput)
		};
	},

	created() {
		this.userId = Number(sessionStorage.getItem("userId"));
	},

	methods: {
		/* Riceve l'invio di un messaggio e lo inoltra a ConversationView per update locale */
		handleMessageSent(message) {
			this.$emit("updateMessages", message, "add");
		},

		/* Scrolla i messaggi */
		scrollToBottom() {
			this.$nextTick(() => {
				const container = this.$refs.messageArea;
				if (container) {
					container.scrollTop = container.scrollHeight;
				}
			});
		},

		/* Gestore degli errori */
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
				case 403:
					this.$root.showToast(`Errore 403: ${detail}`, "error");
					break;
				case 404:
					this.$root.showToast(`Errore 404: ${detail}`, "error");
					break;
				case 409:
					this.$root.showToast(
						"Errore 404: Hai già inserito un commento a questo messaggio",
						"error"
					);
					break;
				case 500:
					this.$root.showToast(`Errore 500:\n ${detail}!`, "error");
					break;
				default:
					this.$root.showToast("Errore sconosciuto", "error");
			}
		},

		/* -------------- GESTIONE AZIONI -------------- */

		/* *** Delete Message *** */

		/* Riceve l'evento di cancellazione di un messaggio e apre il ConfirmDialog */
		handleDeleteMessage(message) {
			this.confirmingMessage = message;
		},

		/* Riceve la conferma dell'eliminazione di un messaggio,
			 invia la richiesta al backend 
			 inoltra a ConversationView per update locale */
		async confirmDeleteMessage() {
			const message = this.confirmingMessage;
			this.confirmingMessage = null;
			try {
				await this.$axios.delete(
					`/conversation-messages/${this.conversationId}/${message.messageId}`,
					{
						headers: {
							Authorization: `Bearer ${this.userId}`,
						},
					}
				);
				this.$emit("updateMessages", message, "delete");
				this.$root.showToast(
					"Messaggio eliminato con successo",
					"success"
				);
			} catch (err) {
				console.error("Errore durante la cancellazione:", err);
				this.handleError(err);
			}
		},

		/* Riceve l'annullamento della cancellazione di un messaggio e chiude il ConfirmDialog */
		async cancelDeleteMessage() {
			this.confirmingMessage = null;
		},

		/* *** Comment Message ***  */

		/* Riceve l'evento di commento e apre il commentPicker nella posizione corretta */
		openCommentPicker(message, event) {
			this.activeMessageForComment = message.messageId;

			this.$nextTick(() => {
				const target = event.target;
				const messageArea = this.$refs.messageArea;

				if (target && messageArea) {
					const areaRect = messageArea.getBoundingClientRect();
					const targetRect = target.getBoundingClientRect();

					this.pickerPosition = {
						top: targetRect.bottom - areaRect.top + 70,
						left:
							targetRect.left -
							areaRect.left +
							targetRect.width / 2 -
							100,
					};
				}
			});
		},

		/* Riceve l'evento da CommentPicker,
		   invia la richiesta al backend,
			 inoltra l'evento a ConversationView per aggiornamento locale */
		async addComment(emoji) {
			const messageId = this.activeMessageForComment;
			this.activeMessageForComment = null;

			try {
				const res = await this.$axios.post(
					`/messages/${messageId}/comments`,
					{ content: emoji },
					{
						headers: {
							Authorization: `Bearer ${this.userId}`,
						},
					}
				);

				// Aggiorno localmente i commenti del messaggio
				const newComment = res.data;
				this.$emit("addComment", { messageId, newComment });
			} catch (err) {
				console.error("Errore nel salvataggio commento", err);
				this.handleError(err);
			}
		},

		/* Riceve l'evento da CommentPopup,
			invia la richiesta al backend,
			inoltra l'evento a ConversationView per aggiornamento locale */
		async handleDeleteComment(comment) {
			try {
				await this.$axios.delete(
					`/messages/${comment.messageId}/comments`,
					{
						headers: {
							Authorization: `Bearer ${this.userId}`,
						},
					}
				);

				// aggiorno la lista locale
				this.$emit("deleteComment", comment);
			} catch (err) {
				console.error("Errore nella cancellazione del commento", err);
				this.handleError(err);
			}
		},

		/* Apre il CommentPopup */
		openCommentPopup(comment) {
			this.$emit("openCommentPopup", comment);
		},

		/* Chiude il CommentPopup */
		closeCommentPopup() {
			this.$emit("closeCommentPopup");
		},

		/* *** Forward Message ***  */

		/* Riceve l'evento di inoltro e apre ForwardPopup  */
		async handleForward(message) {
			this.forwardingMessage = message;
		},

		/* Riceve l'evento da ForwardPopup e lo inoltra a ConversationView Per aggiornamento locale */
		messageForwarded() {
			this.forwardingMessage = null;
			this.$emit("messageForwarded");
		},

		/* *** Reply Message *** */

		/* Riceve l'evento di risposta e modifica MessageInput  */
		handleReply(message) {
			this.replyingToMessage = message;
		},

		/* Riceve l'evento da MessageInput e lo modifica  */
		clearReply() {
			this.replyingToMessage = null;
		},

		/* -------------- GESTIONE GRUPPI -------------- */

		/* Rimuove l'utente corrente dal gruppo (se era l'ultimo, pusha alla Home */
		async leaveGroup() {
			try {
				await this.$axios.delete(
					`/conversation-members/${this.conversationId}/me`,
					{
						headers: {
							Authorization: `Bearer ${this.userId}`,
						},
					}
				);
				this.$router.push("/Home");
				if (this.members.length === 0) {
					this.$root.showToast(
						"Hai abbandonato il gruppo! \n Gruppo eliminato perchè vuoto!"
					);
				} else {
					this.$root.showToast(
						"Hai abbandonato il gruppo!",
						"success"
					);
				}
			} catch (err) {
				console.error("Errore nell'abbandono del gruppo", err);
				this.handleError(err);
			}
		},

		/* Aggiunge una serie di utenti al gruppo e inoltra a ConversationView per update locale */
		async addMembers(newMembers) {
			for (const member of newMembers) {
				try {
					await this.$axios.post(
						`/conversation-add/${this.conversationId}`,
						{ userId: member.userId },
						{
							headers: {
								Authorization: `Bearer ${this.userId}`,
							},
						}
					);

					// Aggiorno i membri localmente
					this.$emit("updateGroupMembers", member);
				} catch (err) {
					console.error("Errore nell'inserimento al gruppo", err);
					this.handleError(err);
				}
			}
			if (newMembers.lenght === 1) {
				this.$root.showToast("Membro aggiunto con successo", "success");
			} else {
				this.$root.showToast("Membri aggiunti con successo");
			}
		},

		/* Aggiorna il nome del gruppo e inoltra a ConversationView per update locale */
		async updateGroupName(newName) {
			try {
				await this.$axios.put(
					`/conversation-name/${this.conversationId}`,
					{ name: newName },
					{
						headers: {
							Authorization: `Bearer ${this.userId}`,
						},
					}
				);

				this.$root.showToast(
					"Nome del gruppo modificato con successo!",
					"success"
				);
				// Aggiorno il nome localmente
				this.$emit("updateGroupName", newName);
			} catch (err) {
				console.error(
					"Errore nell'aggiornamento del nome del gruppo",
					err
				);
				this.handleError(err);
			}
		},

		/* Aggiorna la foto del gruppo e inoltra a ConversationView per update locale */
		async updateGroupPhoto(newPhoto) {
			try {
				await this.$axios.put(
					`/conversation-photo/${this.conversationId}`,
					{ photo: newPhoto },
					{
						headers: {
							Authorization: `Bearer ${this.userId}`,
						},
					}
				);

				this.$root.showToast(
					"Foto del gruppo modificata con successo!",
					"success"
				);

				// Aggiorno la foto localmente
				this.$emit("updateGroupPhoto", newPhoto);
			} catch (err) {
				console.error(
					"Errore nell'aggiornamento della foto del gruppo",
					err
				);
				this.handleError(err);
			}
		},
	},
};
</script>

<style scoped>
.conversation-window {
	flex-grow: 1;
	display: flex;
	flex-direction: column;
	background-color: #f0f0f0;
	width: 100%;
	max-width: 100%;
	height: 100%;
	overflow: hidden;
}

.main-content {
	position: relative;
	flex-grow: 1;
	display: flex;
	flex-direction: column;
	height: 100%;
	overflow: hidden;
}

.message-area {
	flex-grow: 1;
	overflow-y: auto;
	overflow-x: hidden;
	padding: 20px;
	display: flex;
	align-items: flex-start;
	justify-content: flex-start;
	flex-direction: column;
	color: #777;
	font-size: 1rem;
	scroll-behavior: smooth;
}

.no-messages {
	color: #999;
}

.info-overlay-container {
	position: absolute;
	top: 0;
	left: 0;
	width: 100%;
	height: 100%;
	background: white;
	z-index: 10;
	display: flex;
	flex-direction: column;
	animation: slideDown 0.4s ease forwards;
}

@keyframes slideDown {
	from {
		transform: translateY(-20px);
		opacity: 0;
	}
	to {
		transform: translateY(0);
		opacity: 1;
	}
}
</style>
