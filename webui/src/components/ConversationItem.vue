<template>
	<div
		class="conversation-item"
		@click="$emit('select', conversation.conversationId)"
	>
		<img
			:src="conversation.photo"
			alt="Conversation"
			class="conversation-avatar"
		/>
		<div class="conversation-info">
			<span class="conversation-name">{{ conversation.name }}</span>
			<span
				v-if="conversation.type === 'group'"
				class="conversation-sender"
			>
				{{ conversation.lastMessage.SenderUsername }}:
			</span>
			<span class="conversation-snippet">{{
				conversation.lastMessage.Content
			}}</span>
		</div>
		<span class="conversation-status">{{ conversation.status }}</span>
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
	border-bottom: 1px solid #ddd;
	cursor: pointer;
}
.conversation-avatar {
	width: 40px;
	height: 40px;
	border-radius: 50%;
	margin-right: 10px;
}
.conversation-info {
	flex-grow: 1;
}
</style>
