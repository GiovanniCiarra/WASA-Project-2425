<template>
	<div
		class="conversation-item"
		@click="$emit('conversationSelected', conversation.conversationId)"
	>
		<img
			:src="conversation.photo || '/images/default-user.png'"
			alt="Foto"
			class="conversation-photo"
		/>
		<div class="conversation-info">
			<div class="conversation-name">{{ conversation.name }}</div>
			<div class="conversation-snippet">
				{{ conversation.lastMessage?.content || "Nessun messaggio" }}
			</div>
		</div>
	</div>
</template>

<script>
export default {
	props: {
		conversation: Object,
	},
};
</script>

<style scoped>
.conversation-item {
	display: flex;
	align-items: center;
	padding: 10px;
	border-bottom: 1px solid #eee;
	cursor: pointer;
}

.conversation-photo {
	width: 40px;
	height: 40px;
	border-radius: 50%;
	margin-right: 10px;
	object-fit: cover;
}

.conversation-info {
	flex: 1;
}

.conversation-name {
	font-weight: bold;
}

.conversation-snippet {
	color: #666;
	font-size: 0.9rem;
	margin-top: 2px;
	overflow: hidden;
	text-overflow: ellipsis;
	white-space: nowrap;
	max-width: 180px;
}
</style>
