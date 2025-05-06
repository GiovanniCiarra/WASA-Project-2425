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
			<div v-if="message.isForwarded" class="forwarded-label">
				<img
					src="../assets/icons/forwarded.png"
					alt="Inoltrato"
					class="forwarded-icon"
				/>
				Inoltrato
			</div>
			<!-- Header -->
			<div class="message-header">
				<span class="sender">{{ message.senderUsername }}</span>
				<span class="timestamp">{{ formattedTimestamp }}</span>
			</div>

			<!-- Contenuto -->
			<div class="message-content">
				<!-- Blocco risposta -->
				<div
					v-if="message.isReply && message.isReply !== -1"
					class="reply-block"
				>
					<strong>{{ message.replyTo.senderUsername }}</strong>
					<div class="reply-content">
						<span v-if="message.replyTo.type === 'text'">{{
							message.replyTo.content
						}}</span>
						<span v-else>[Immagine]</span>
					</div>
				</div>

				<!-- Testo -->
				<p v-if="message.type === 'text'">{{ message.content }}</p>
				<!-- Immagine -->
				<img
					v-else-if="message.type === 'photo'"
					:src="message.content"
					alt="Foto inviata"
					class="message-image"
				/>
			</div>

			<!-- Commenti -->
			<div
				class="message-comments"
				v-if="message.comments && message.comments.length > 0"
			>
				<span
					v-for="(c, index) in displayedComments"
					:key="c.commentId"
					class="emoji-comment"
					@click="$emit('clickComment', c)"
					:style="{ zIndex: message.comments.length - index }"
				>
					{{ c.content }}
				</span>
				<img
					v-if="message.comments.length > 3"
					class="emoji-comment more-comments"
					src="../assets/icons/more.png"
				/>
			</div>

			<!-- Status -->
			<div v-if="isMine" class="status-icon">
				<!-- Sent -->
				<img
					v-if="message.status === 'sent'"
					:src="checkGray"
					alt="Sent"
					class="status-icon-img"
				/>

				<!-- Received -->
				<div
					v-else-if="message.status === 'received'"
					class="status-double-icon"
				>
					<img :src="checkGray" alt="Received" />
					<img :src="checkGray" alt="Received" />
				</div>

				<!-- Read -->
				<div
					v-else-if="message.status === 'read'"
					class="status-double-icon"
				>
					<img :src="checkBlue" alt="Read" />
					<img :src="checkBlue" alt="Read" />
				</div>
			</div>
		</div>

		<!-- Menu Azioni -->
		<div class="message-actions">
			<button @click="showActions = !showActions">
				<img
					:src="showActions ? closeActions : dotsMenu"
					alt="Azioni"
				/>
			</button>

			<transition name="fade-slide">
				<!-- Griglia Azioni -->
				<div v-if="showActions" class="message-actions-grid">
					<!-- Reply -->
					<button
						@click="$emit('reply', message)"
						title="Rispondi"
						@mouseenter="hoveredReply = true"
						@mouseleave="hoveredReply = false"
					>
						<img
							:src="hoveredReply ? replyHover : replyDefault"
							alt="Reply"
						/>
					</button>

					<!-- Forward -->
					<button
						@click="$emit('forward', message)"
						title="Inoltra"
						@mouseenter="hoveredForward = true"
						@mouseleave="hoveredForward = false"
					>
						<img
							:src="
								hoveredForward ? forwardHover : forwardDefault
							"
							alt="Forward"
						/>
					</button>

					<!-- Comment -->
					<button
						@click="$emit('comment', message, $event)"
						title="Commenta"
						@mouseenter="hoveredComment = true"
						@mouseleave="hoveredComment = false"
					>
						<img
							:src="
								hoveredComment ? commentHover : commentDefault
							"
							alt="Comment"
						/>
					</button>
					<!-- Delete -->
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
			</transition>
		</div>
	</div>
</template>

<script>
import dotsMenu from "../assets/icons/menu-actions.png";
import deleteRed from "../assets/icons/delete-red.png";
import deleteBlack from "../assets/icons/delete.png";
import forwardBlue from "../assets/icons/forward-blue.png";
import forwardBlack from "../assets/icons/forward.png";
import commentGreen from "../assets/icons/comment-green.png";
import commentBlack from "../assets/icons/comment.png";
import replyBlack from "../assets/icons/reply.png";
import replyYellow from "../assets/icons/reply-yellow.png";
import forwarded from "../assets/icons/forwarded.png";
import checkGray from "../assets/icons/checkGray.png";
import checkBlue from "../assets/icons/checkBlue.png";
import closeActions from "../assets/icons/close.png";

export default {
	data() {
		return {
			dotsMenu,
			hoveredDelete: false,
			deleteHover: deleteRed,
			deleteDefault: deleteBlack,
			hoveredForward: false,
			forwardHover: forwardBlue,
			forwardDefault: forwardBlack,
			hoveredComment: false,
			commentHover: commentGreen,
			commentDefault: commentBlack,
			hoveredReply: false,
			replyHover: replyYellow,
			replyDefault: replyBlack,
			forwarded,
			checkGray,
			checkBlue,
			closeActions,
			showActions: false,
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

		displayedComments() {
			return this.message.comments.slice(0, 3);
		},
	},

	mounted() {
		document.addEventListener("click", this.closeActionsMenu);
	},
	beforeDestroy() {
		document.removeEventListener("click", this.closeActionsMenu);
	},
	methods: {
		closeActionsMenu(event) {
			const actionsGrid = this.$el.querySelector(".message-actions-grid");
			const actionsButton = this.$el.querySelector(
				".message-actions button"
			);
			if (
				this.showActions &&
				actionsGrid &&
				!actionsGrid.contains(event.target) &&
				!actionsButton.contains(event.target)
			) {
				this.showActions = false;
			}
		},
	},
};
</script>

<style scoped>
.message-wrapper {
	display: flex;
	gap: 5px;
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
	padding-bottom: 5px;
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

.forwarded-label {
	background-color: #f1f1f1;
	color: #555;
	font-size: 0.7rem;
	padding: 2px 8px;
	border-radius: 8px;
	margin-bottom: 4px;
	display: inline-flex;
	align-items: center;
	text-transform: uppercase;
	letter-spacing: 0.5px;
}

.forwarded-icon {
	width: 12px;
	height: 12px;
	margin-right: 6px;
	margin-bottom: 2px;
	vertical-align: middle;
}

.message-image {
	max-width: 100%;
	max-height: 300px;
	border-radius: 10px;
	margin-top: 5px;
	object-fit: cover;
	cursor: pointer;
	transition: transform 0.2s;
}

.message-image:hover {
	transform: scale(1.05);
}

.mine {
	align-self: flex-end;
	background-color: #33a07a;
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
	white-space: nowrap;
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
	align-self: flex-end;
}

.status-icon-img {
	width: 14px;
	height: 14px;
	margin-top: 5px;
}

.status-double-icon {
	display: flex;
}

.status-double-icon img {
	width: 14px;
	height: 14px;
}

.message-actions {
	display: flex;
	flex-direction: column;
	justify-content: center;
	padding-top: 5px;
	position: relative;
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

.message-actions-grid {
	display: grid;
	grid-template-columns: repeat(2, 1fr);
	position: absolute;
	top: 100%;
}

.message-wrapper.align-left .message-actions-grid {
	left: 0;
	right: auto;
}

.message-wrapper.align-right .message-actions-grid {
	right: 0;
	left: auto;
}

.fade-slide-enter-active,
.fade-slide-leave-active {
	transition: opacity 0.3s ease, transform 0.3s ease;
}

.fade-slide-enter-from,
.fade-slide-leave-to {
	opacity: 0;
	transform: translateY(-10px);
}

.fade-slide-enter-to,
.fade-slide-leave-from {
	opacity: 1;
	transform: translateY(0);
}

.message-comments {
	position: absolute;
	bottom: -10px;
	right: 6px;
	display: flex;
}

.mine .message-comments {
	right: auto;
	left: 6px;
}

.emoji-comment {
	font-size: 1rem;
	cursor: pointer;
	transition: transform 0.2s;
	background-color: rgba(255, 255, 255, 0.834);
	border-radius: 50%;
	padding: 2px;
	box-shadow: 0 0 2px rgba(0, 0, 0, 0.2);
	position: relative;
	margin-left: -6px;
}

.emoji-comment:first-child {
	margin-left: 0;
}

.emoji-comment:hover {
	transform: scale(1.2);
}

.more-comments {
	font-size: 1rem;
	color: #888;
	font-weight: bold;
}

.reply-block {
	background-color: rgba(0, 0, 0, 0.05);
	border-left: 3px solid #ccc;
	padding: 6px 10px;
	margin-bottom: 6px;
	border-radius: 8px;
	font-size: 0.85rem;
	color: #555;
}

.mine .reply-block {
	background-color: rgba(255, 255, 255, 0.2);
	border-left: 3px solid #ddd;
	color: #f1f1f1;
}

.reply-content {
	margin-top: 4px;
	font-style: italic;
	font-size: 0.8rem;
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
}
</style>
