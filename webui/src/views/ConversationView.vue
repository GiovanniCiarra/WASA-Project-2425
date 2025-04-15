<template>
	<div class="conversation-view">
		<!-- Lista conversazioni a sinistra -->
		<div class="sidebar-conversations">
			<input
				v-model="searchQuery"
				type="text"
				placeholder="Cerca conversazioni..."
				class="search-bar"
			/>
			<div class="scrollable">
				<div
					v-for="conv in filteredConversations"
					:key="conv.conversationId"
				>
					<ConversationItem
						:conversation="conv"
						@click="
							$router.push({
								path: `/conversations/${conv.conversationId}`,
								query: {
									type: conv.type,
									name: conv.name,
									photo: conv.photo,
								},
							})
						"
					/>
				</div>
			</div>
		</div>

		<!-- Finestra messaggi -->
		<ConversationWindow
			:conversationId="conversationId"
			@messageSent="messageSent = true"
		/>
	</div>
</template>

<script>
import ConversationItem from "../components/ConversationItem.vue";
import ConversationWindow from "../components/ConversationWindow.vue";

export default {
	components: { ConversationItem, ConversationWindow },

	data() {
		return {
			conversations: [],
			userId: null,
			messageSent: false,
			searchQuery: "",
		};
	},
	computed: {
		conversationId() {
			return this.$route.params.conversationId;
		},

		filteredConversations() {
			return this.conversations.filter((conv) =>
				conv.name.toLowerCase().includes(this.searchQuery.toLowerCase())
			);
		},
	},

	created() {
		this.fetchConversations();
	},

	methods: {
		async fetchConversations() {
			this.userId = sessionStorage.getItem("userId");
			try {
				const res = await this.$axios.get(
					"http://localhost:3000/conversations",
					{
						headers: {
							Authorization: `Bearer ${this.userId}`,
						},
					}
				);
				this.conversations = res.data;
				console.log(this.conversations);
			} catch (err) {
				console.error("Errore nel caricamento conversazioni", err);
			}
		},
	},

	beforeRouteLeave(_, __, next) {
		if (!this.messageSent) {
			console.log("CIAONE");
			try {
				this.$axios.delete(
					`http://localhost:3000/conversations/${this.conversationId}`,
					{
						headers: {
							Authorization: `Bearer ${this.userId}`,
						},
					}
				);
				console.log("Conversazione eliminata perché vuota");
				next();
			} catch (err) {
				console.error("Errore nella cancellazione conversazione", err);
				next();
			}
		} else {
			next();
		}
	},
};
</script>

<style>
.conversation-view {
	display: flex;
	height: 100%;
	width: 100%;
	background-color: #ece5dd;
}

.search-bar {
	border: 2px solid transparent;
	width: 100%;
	height: 2.5em;
	padding-left: 0.8em;
	outline: none;
	background-color: #f3f3f3;
	border-radius: 10px;
	transition: all 0.5s;
	margin-top: 10px;
	box-sizing: border-box;
}

.search-bar:hover,
.search-bar:focus {
	border: 2px solid #4a9dec;
	background-color: white;
}

.sidebar-conversations {
	width: 340px;
	background-color: white;
	border-right: 1px solid #ddd;
	overflow-y: auto;
	padding: 0 10px;
}

.scrollable {
	padding: 5px;
	display: flex;
	flex-direction: column;
	gap: 8px;
}
</style>
