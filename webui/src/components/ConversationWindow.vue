<template>
	<div class="conversation-window">
		<div class="main-content">
			<ConversationHeader
				:name="conversationName"
				:photo="conversationPhoto"
				:type="conversationType"
				@openInfo="$emit('openCloseConversationInfo')"
			/>

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
					@forwarded="forwardingMessage = null"
				/>

				<div v-if="messages.length === 0" class="no-messages">
					Nessun messaggio
				</div>
			</div>

			<ConfirmDialog
				v-if="confirmingMessage"
				:message="'Vuoi davvero eliminare questo messaggio?'"
				@confirm="confirmDeleteMessage"
				@cancel="cancelDeleteMessage"
			/>

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
		messages: Array,
		members: Array,
		conversationName: String,
		conversationPhoto: String,
		conversationType: String,
		commentPopupMessage: Object,
		conversationInfoOpen: Boolean,
	},

	data() {
		return {
			userId: null,
			confirmingMessage: null,
			activeMessageForComment: null,
			pickerPosition: { top: 0, left: 0 },
			forwardingMessage: null,
			replyingToMessage: null,
		};
	},

	created() {
		this.userId = Number(sessionStorage.getItem("userId"));
	},

	methods: {
		handleMessageSent(message) {
			this.$emit("updateMessages", message, "add");
		},

		scrollToBottom() {
			this.$nextTick(() => {
				const container = this.$refs.messageArea;
				if (container) {
					container.scrollTop = container.scrollHeight;
				}
			});
		},

		/* -------------- GESTIONE AZIONI -------------- */

		/* --- Delete Message --- */
		handleDeleteMessage(message) {
			this.confirmingMessage = message;
		},

		async confirmDeleteMessage() {
			const message = this.confirmingMessage;
			this.confirmingMessage = null;
			try {
				await this.$axios.delete(
					`/conversations/${this.conversationId}/messages/${message.messageId}`,
					{
						headers: {
							Authorization: `Bearer ${this.userId}`,
						},
					}
				);
				this.$emit("updateMessages", message, "delete");
			} catch (err) {
				console.error("Errore durante la cancellazione:", err);
			}
		},

		async cancelDeleteMessage() {
			this.confirmingMessage = null;
		},

		/* --- Comment Message --- */
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

				// Aggiorno localmente il messaggio
				const newComment = res.data;
				this.$emit("addComment", { messageId, newComment });
			} catch (err) {
				console.error("Errore nel salvataggio commento", err);
			}
		},

		openCommentPopup(comment) {
			this.$emit("openCommentPopup", comment);
		},

		closeCommentPopup() {
			this.$emit("closeCommentPopup");
		},

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
			}
		},

		/* --- Forward Message --- */
		async handleForward(message) {
			this.forwardingMessage = message;
		},

		/* --- Reply Message --- */
		handleReply(message) {
			this.replyingToMessage = message;
		},

		clearReply() {
			this.replyingToMessage = null;
		},

		/* -------------- GESTIONE GRUPPI -------------- */

		async leaveGroup() {
			try {
				await this.$axios.delete(
					`/conversations/${this.conversationId}/leave`,
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
			}
		},

		async addMembers(newMembers) {
			for (const member of newMembers) {
				try {
					await this.$axios.post(
						`/conversations/${this.conversationId}/add`,
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
				}
			}
			if (newMembers.lenght === 1) {
				this.$root.showToast("Membro aggiunto con successo", "success");
			} else {
				this.$root.showToast("Membri aggiunti con successo");
			}
		},

		async updateGroupName(newName) {
			try {
				await this.$axios.put(
					`/conversations/${this.conversationId}/name`,
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
			}
		},

		async updateGroupPhoto(newPhoto) {
			try {
				await this.$axios.put(
					`/conversations/${this.conversationId}/photo`,
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
