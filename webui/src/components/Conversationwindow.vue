<template>
	<div class="conversation-window">
		<ConversationHeader :name="name" :photo="photo" :type="type" />

		<div class="message-area" ref="messageArea">
			<Message
				v-for="msg in messages"
				:key="msg.messageId"
				:message="msg"
				:currentUser="this.userId"
				@forward="handleForward"
				@comment="handleComment"
				@delete="handleDelete"
			/>

			<div v-if="messages.length === 0" class="no-messages">
				Nessun messaggio
			</div>
		</div>

		<MessageInput
			:conversationId="conversationId"
			@messageSent="handleMessageSent"
		/>
	</div>
</template>

<script>
import ConversationHeader from "../components/ConversationHeader.vue";
import MessageInput from "../components/MessageInput.vue";
import Message from "../components/Message.vue";

export default {
	components: {
		ConversationHeader,
		MessageInput,
		Message,
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
	},

	watch: {
		"$route.params.conversationId"(newId, oldId) {
			if (newId !== oldId) {
				this.fetchMessages();
			}
		},
	},

	methods: {
		async fetchMessages() {
			try {
				const res = await this.$axios.get(
					`http://localhost:3000/conversations/${this.conversationId}`,
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
</style>
