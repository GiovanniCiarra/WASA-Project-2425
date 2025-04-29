<template>
	<div class="conversation-window">
		<div class="main-content">
			<ConversationHeader
				:name="name"
				:photo="photo"
				:type="type"
				@openInfo="infoOpen = true"
			/>

			<div v-if="infoOpen" class="info-overlay-container">
				<!-- Se è un gruppo -->
				<GroupInfoOverlay
					v-if="type === 'group'"
					:groupName="name"
					:groupPhoto="photo"
					:members="members"
					:currentUser="userId"
					@closeOverlay="infoOpen = false"
					@leaveGroup="leaveGroup"
					@addMembers="addMembers"
					@updateGroupPhoto="updateGroupPhoto"
					@updateGroupName="updateGroupName"
				/>

				<!-- Se è una chat -->
				<ChatInfoOverlay
					v-else
					:chatName="name"
					:chatPhoto="photo"
					@closeOverlay="infoOpen = false"
				/>
			</div>
			<div class="message-area" ref="messageArea">
				<Message
					v-for="msg in messages"
					:key="msg.messageId"
					:message="msg"
					:currentUser="this.userId"
					@forward="handleForward"
					@comment="handleComment"
					@delete="handleDelete"
					@clickComment="handleClickComment"
				/>
				<CommentPicker
					v-if="activeMessageForComment"
					:style="{
						top: pickerPosition.top + 'px',
						left: pickerPosition.left + 'px',
					}"
					@select-emoji="handleEmojiSelect"
					@close-emoji="activeMessageForComment = null"
				/>

				<CommentPopUp
					v-if="commentPopup"
					:comments="commentPopup.comments"
					:currentUser="this.userId"
					@close="commentPopup = null"
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
				@confirm="confirmDelete"
				@cancel="cancelDelete"
			/>

			<MessageInput
				:conversationId="conversationId"
				@messageSent="handleMessageSent"
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
	},

	data() {
		return {
			messages: [],
			userId: null,
			members: [],
			infoOpen: false,
			confirmingMessage: null,
			activeMessageForComment: null,
			pickerPosition: { top: 0, left: 0 },
			commentPopup: null,
			forwardingMessage: null,
		};
	},

	computed: {
		name() {
			return this.$route.query.name;
		},
		photo() {
			return this.$route.query.photo;
		},
		type() {
			return this.$route.query.type;
		},
	},

	created() {
		this.userId = Number(sessionStorage.getItem("userId"));
		this.fetchMessages();
		if (this.type === "group") {
			this.fetchMembers();
		}
	},

	watch: {
		"$route.params.conversationId"(newId, oldId) {
			if (newId !== oldId) {
				this.fetchMessages();
				if (this.type === "group") {
					this.fetchMembers();
				}
				this.infoOpen = false;
			}
		},
	},

	methods: {
		async fetchMessages() {
			try {
				const res = await this.$axios.get(
					`/conversations/${this.conversationId}`,
					{
						headers: {
							Authorization: `Bearer ${this.userId}`,
						},
					}
				);
				if (res.data) {
					this.messages = res.data;
					this.$emit("messageSent");
					this.scrollToBottom();
				}
			} catch (err) {
				console.error("Errore nel caricamento dei messaggi", err);
			}
		},

		async fetchMembers() {
			try {
				const res = await this.$axios.get(
					`/conversations/${this.conversationId}/members`,
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
			}
		},

		addMessage(newMessage) {
			this.messages.push(newMessage);
			this.scrollToBottom();
		},

		handleMessageSent(message) {
			this.addMessage(message);
			this.$emit("messageSent");
		},

		scrollToBottom() {
			this.$nextTick(() => {
				const container = this.$refs.messageArea;
				if (container) {
					container.scrollTop = container.scrollHeight;
				}
			});
		},

		handleDelete(message) {
			this.confirmingMessage = message;
		},

		async confirmDelete() {
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
				// Rimuovo il messaggio nell'array corrente
				this.messages = this.messages.filter(
					(m) => m.messageId !== message.messageId
				);

				// Se il messaggio eliminato era l'ultimo della conversazione, torno alla Home
				if (this.messages.length === 0) {
					this.$router.push("/Home");
					this.$root.showToast(
						"Conversazione eliminata perchè vuota!"
					);
				}
			} catch (err) {
				console.error("Errore durante la cancellazione:", err);
			}
		},

		async cancelDelete() {
			this.confirmingMessage = null;
		},

		handleComment(message, event) {
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

		async handleEmojiSelect(emoji) {
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
				const msg = this.messages.find(
					(m) => m.messageId === messageId
				);
				if (msg) {
					msg.comments = msg.comments || [];
					msg.comments.push(newComment);
				}
			} catch (err) {
				console.error("Errore nel salvataggio commento", err);
			}
		},

		handleClickComment(comment) {
			const message = this.messages.find(
				(m) => m.messageId === comment.messageId
			);

			if (message) {
				this.commentPopup = message;
			}
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
				this.commentPopup.comments = this.commentPopup.comments.filter(
					(c) => c.commentId !== comment.commentId
				);
			} catch (err) {
				console.error("Errore nella cancellazione del commento", err);
			}
		},

		async handleForward(message) {
			this.forwardingMessage = message;
		},

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
				// aggiorno la lista locale
				this.members = this.members.filter(
					(m) => m.userId !== this.userId
				);
				this.$router.push("/Home");
				this.$root.showToast("hai abbandonato il gruppo!");
				if (this.members.length === 0) {
					this.$root.showToast("Gruppo eliminato perchè vuoto!");
				}
			} catch (err) {
				console.error("Errore nell'abbandono del gruppo", err);
			}
		},

		async addMembers(memberId) {},
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
