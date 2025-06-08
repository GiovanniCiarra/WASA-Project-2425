<template>
	<div class="comment-details-overlay">
		<div class="comment-details">
			<div class="comment-details-header">
				<!-- popup title -->
				<span class="title">Commenti al messaggio</span>

				<!-- icona chiusura popup -->
				<img
					@click="$emit('close')"
					class="close-btn"
					src="../assets/icons/close.png"
					alt="Close"
				/>
			</div>

			<!-- label in caso di mancanza di commenti -->
			<div v-if="comments.length === 0" class="empty">
				Nessun commento
			</div>
			<ul v-else>
				<!-- Lista di commenti -->
				<li v-for="comment in comments" :key="comment.commentId">
					<span class="emoji">{{ comment.content }}</span>
					<div class="user-info">
						<span class="username">{{ comment.sender }}</span>
						<span
							v-if="Number(comment.senderId) === currentUser"
							class="you-label"
							>(tu)</span
						>
					</div>

					<!-- Delete button -->
					<button
						v-if="Number(comment.senderId) === currentUser"
						class="delete-btn"
						@click="$emit('deleteComment', comment)"
					>
						Elimina
					</button>
				</li>
			</ul>
		</div>
	</div>
</template>

<script>
export default {
	props: {
		comments: Array, // Array dei commenti al messaggio
		currentUser: Number, // userId dell'utente corrente
	},
};
</script>

<style scoped>
.comment-details-overlay {
	position: fixed;
	top: 0;
	left: 0;
	width: 100vw;
	height: 100vh;
	background-color: rgba(0, 0, 0, 0.35);
	display: flex;
	justify-content: center;
	align-items: center;
	z-index: 10000;
}

.comment-details {
	background: white;
	padding: 16px 20px;
	border-radius: 10px;
	width: 300px;
	max-height: 70vh;
	overflow-y: auto;
	box-shadow: 0 0 12px rgba(0, 0, 0, 0.2);
}

.comment-details-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-bottom: 12px;
}

.title {
	font-weight: bold;
}

ul {
	list-style: none;
	padding: 0;
}

li {
	display: flex;
	justify-content: space-between;
	align-items: center;
	padding: 6px 0;
	border-bottom: 1px solid #eee;
}

.emoji {
	font-size: 1.3rem;
}

.user-info {
	display: flex;
	align-items: center;
	flex-grow: 1;
	margin-left: 10px;
}

.username {
	color: #444;
}

.you-label {
	font-size: 0.8rem;
	color: #888;
	margin-left: 6px;
}

.delete-btn {
	background: transparent;
	border: none;
	color: red;
	font-size: 0.85rem;
	cursor: pointer;
}

.close-btn {
	background: none;
	border: none;
	font-size: 1.1rem;
	cursor: pointer;
	opacity: 0.6;
}

.close-btn:hover {
	opacity: 1;
}

.empty {
	margin-left: 25%;
}
</style>
