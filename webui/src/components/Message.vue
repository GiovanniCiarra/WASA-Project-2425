<template>
	<div
		class="message-wrapper"
		:class="{
			reverse: isMine,
			'align-right': isMine,
			'align-left': !isMine,
		}"
	>
		<!-- Box messaggio -->
		<div :class="['message', isMine ? 'mine' : 'received']">
			<!-- Header -->
			<div class="message-header">
				<span class="sender">{{ message.senderUsername }}</span>
				<span class="timestamp">{{ formattedTimestamp }}</span>
			</div>

			<!-- Contenuto -->
			<div class="message-content">
				<p v-if="message.type === 'text'">{{ message.content }}</p>
				<!-- In futuro: immagine -->
			</div>

			<!-- Commenti -->
			<div
				class="message-comments"
				v-if="message.comments && message.comments.length > 0"
			>
				<span
					v-for="c in message.comments"
					:key="c.commentId"
					class="emoji-comment"
					@click="$emit('clickComment', c)"
				>
					{{ c.content }}
				</span>
			</div>

			<!-- Status -->
			<div v-if="isMine" class="status-icon">
				<span v-if="message.status === 'sent'">✔</span>
				<span v-else-if="message.status === 'received'">✔✔</span>
				<span
					v-else-if="message.status === 'read'"
					style="color: #4fc3f7"
					>✔✔</span
				>
			</div>
		</div>

		<!-- Azioni -->
		<div class="message-actions">
			<button
				@click="$emit('forward', message)"
				title="Inoltra"
				@mouseenter="hoveredForward = true"
				@mouseleave="hoveredForward = false"
			>
				<img
					:src="hoveredForward ? forwardHover : forwardDefault"
					alt="Forward"
				/>
			</button>
			<button
				@click="$emit('comment', message, $event)"
				title="Commenta"
				@mouseenter="hoveredComment = true"
				@mouseleave="hoveredComment = false"
			>
				<img
					:src="hoveredComment ? commentHover : commentDefault"
					alt="Comment"
				/>
			</button>
			<button
				v-if="isMine"
				@click="$emit('delete', message)"
				title="Elimina"
				@mouseenter="hoveredDelete = true"
				@mouseleave="hoveredDelete = false"
			>
				<img
					:src="hoveredDelete ? deleteHover : deleteDefault"
					alt="Delete"
				/>
			</button>
		</div>
	</div>
</template>

<script>
import deleteRed from "../assets/icons/delete-red.png";
import deleteBlack from "../assets/icons/delete.png";
import forwardBlue from "../assets/icons/forward-blue.png";
import forwardBlack from "../assets/icons/forward.png";
import commentGreen from "../assets/icons/comment-green.png";
import commentBlack from "../assets/icons/comment.png";

export default {
	data() {
		return {
			hoveredDelete: false,
			deleteHover: deleteRed,
			deleteDefault: deleteBlack,
			hoveredForward: false,
			forwardHover: forwardBlue,
			forwardDefault: forwardBlack,
			hoveredComment: false,
			commentHover: commentGreen,
			commentDefault: commentBlack,
		};
	},

	props: {
		message: {
			type: Object,
			required: true,
		},
		currentUser: {
			type: Number,
			required: true,
		},
	},

	computed: {
		isMine() {
			return this.message.senderId === this.currentUser;
		},

		formattedTimestamp() {
			// Mostra solo ora e minuti per semplicità
			const date = new Date(this.message.timestamp);
			return date.toLocaleString("it-IT", {
				weekday: "short",
				day: "2-digit",
				month: "2-digit",
				hour: "2-digit",
				minute: "2-digit",
			});
		},
	},
};
</script>

<style scoped>
.message-wrapper {
	display: flex;
	gap: 6px;
	align-items: flex-start;
	max-width: 100%;
}

.message-wrapper.align-right {
	align-self: flex-end;
}

.message-wrapper.align-left {
	align-self: flex-start;
}

.message-wrapper.reverse {
	flex-direction: row-reverse;
}

.message {
	max-width: 85%;
	padding: 10px 14px;
	padding-bottom: 18px;
	margin: 6px 0;
	border-radius: 14px;
	display: flex;
	flex-direction: column;
	position: relative;
	box-shadow: 0 2px 6px rgba(0, 0, 0, 1);
	word-wrap: break-word;
	overflow-wrap: break-word;
	white-space: normal;
	box-sizing: border-box;
	font-family: sans-serif;
}

.mine {
	align-self: flex-end;
	background-color: #2fa27a;
	color: white;
	text-align: left;
	border-top-right-radius: 4px;
}

.received {
	align-self: flex-start;
	background-color: white;
	color: #333;
	text-align: left;
	border-top-left-radius: 4px;
}

.message-header {
	font-size: 0.8rem;
	color: #666;
	margin-bottom: 6px;
	display: flex;
	justify-content: space-between;
	gap: 8px;
	flex-wrap: nowrap;
}

.sender {
	font-weight: 700;
	font-size: 1rem;
}

.mine .sender {
	color: #fffafa;
}

.received .sender {
	color: #222;
}

.timestamp {
	font-size: 0.8rem;
	color: #d7d4d4;
	white-space: nowrap;
}

.message-content p {
	margin: 0;
	font-size: 1rem;
	color: inherit;
}

.status-icon {
	font-size: 0.9rem;
	margin-top: 4px;
	align-self: flex-end;
}

.message-actions {
	display: flex;
	flex-direction: column;
	gap: 6px;
	justify-content: center;
	padding-top: 10px;
}

.message-actions button {
	background: none;
	border: none;
	cursor: pointer;
	font-size: 0.9rem;
	padding: 2px;
	opacity: 0.6;
	transition: opacity 0.2s, transform 0.2s;
}

.message-actions button:hover {
	opacity: 1;
	transform: scale(1.1);
}

.message-comments {
	position: absolute;
	bottom: 0.5px;
	right: 6px;
	display: flex;
}

.mine .message-comments {
	right: auto;
	left: 6px;
}

.emoji-comment {
	font-size: 1.4rem;
	cursor: pointer;
	transition: transform 0.2s;
}

.emoji-comment:hover {
	transform: scale(1.2);
}
</style>
