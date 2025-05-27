<template>
	<div class="conversation-list">
		<input
			v-model="search"
			placeholder="Cerca conversazioni..."
			class="search-input"
		/>

		<div v-if="filteredConversations.length > 0">
			<ConversationItem
				v-for="conv in filteredConversations"
				:key="conv.conversationId"
				:conversation="conv"
				@conversationSelected="
					$emit('conversationSelected', conv.conversationId)
				"
			/>
		</div>

		<div v-else class="no-conversations">
			Nessuna conversazione trovata.
			<button @click="$emit('newConversation')">
				+ Nuova conversazione
			</button>
		</div>
	</div>
</template>

<script>
import ConversationItem from "./ConversationItem.vue";

export default {
	components: { ConversationItem },
	props: {
		conversations: Array,
	},
	data() {
		return {
			search: "",
		};
	},
	computed: {
		filteredConversations() {
			if (this.conversations) {
				return this.conversations.filter((conv) =>
					conv.name.toLowerCase().includes(this.search.toLowerCase())
				);
			}
			return;
		},
	},
};
</script>

<style scoped>
.conversation-list {
	display: flex;
	flex-direction: column;
	height: 100%;
	padding: 10px;
}

.search-input {
	padding: 8px;
	border-radius: 6px;
	border: 1px solid #ccc;
	margin-bottom: 10px;
	width: 100%;
}

.no-conversations {
	text-align: center;
	color: #777;
	margin-top: 20px;
}

.no-conversations button {
	margin-top: 10px;
	padding: 8px 12px;
	border: none;
	background-color: #008080;
	color: white;
	border-radius: 6px;
	cursor: pointer;
}
</style>
