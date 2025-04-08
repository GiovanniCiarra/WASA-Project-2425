<template>
	<div class="conversation-list">
		<input
			v-model="searchConversation"
			type="text"
			placeholder="Search Conversation..."
			class="search-bar"
		/>
		<div v-if="conversationError" class="error-message">
			{{ conversationError }}
		</div>
		<div
			v-for="conversation in filteredConversations"
			:key="conversation.conversationId"
		>
			<ConversationItem
				:conversation="conversation"
				@select="$emit('conversationSelected', $event)"
			/>
		</div>
	</div>
</template>

<script>
import ConversationItem from "./ConversationItem.vue";

export default {
	components: { ConversationItem },
	props: {
		conversations: Array,
		conversationError: String,
	},
	data() {
		return {
			searchConversation: "",
		};
	},
	computed: {
		filteredConversations() {
			return this.conversations.filter((conversation) =>
				conversation.name
					.toLowerCase()
					.includes(this.searchConversation.toLowerCase())
			);
		},
	},
};
</script>

<style scoped>
.search-bar {
	width: 100%;
	padding: 8px;
	padding-left: 0.8em;
	outline: none;
	overflow: hidden;
	margin-bottom: 10px;
	border: 2px solid transparent;
	border-radius: 10px;
	background-color: #f3f3f3;
	transition: all 0.5s;
}

.search-bar:hover,
.search-bar:focus {
	border: 2px solid #4a9dec;
	box-shadow: 0px 0px 0px 7px rgb(74, 157, 236, 20%);
	background-color: white;
}

.conversation-list {
	max-height: 80vh;
	overflow-y: auto;
}

.error-message {
	color: red;
	font-weight: bold;
	text-align: center;
	margin-top: 10px;
	background: #ffe0e0;
	padding: 10px;
	border-radius: 5px;
}
</style>
