<template>
	<div class="home-container">
		<div class="conversation-panel">
			<input
				v-model="searchQuery"
				type="text"
				class="search-bar"
				placeholder="Cerca conversazioni..."
			/>

			<div v-if="conversationError" class="error-message">
				{{ conversationError }}
			</div>

			<div
				v-if="filteredConversations.length > 0"
				class="conversation-list"
			>
				<ConversationItem
					v-for="conv in filteredConversations"
					:key="conv.conversationId"
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

			<div v-else class="empty-state">
				<p>Nessuna conversazione trovata.</p>
				<button @click="startNewConversation">
					Avvia una nuova chat
				</button>
			</div>
		</div>
	</div>
</template>

<script>
import ConversationItem from "@/components/ConversationItem.vue";

export default {
	components: { ConversationItem },

	data() {
		return {
			conversations: [],
			searchQuery: "",
			conversationError: "",
			userId: null,
		};
	},

	computed: {
		filteredConversations() {
			return this.conversations.filter((conv) =>
				conv.name.toLowerCase().includes(this.searchQuery.toLowerCase())
			);
		},
	},

	methods: {
		startNewConversation() {
			this.$router.push("/NewChat");
		},

		async fetchConversations() {
			this.userId = sessionStorage.getItem("userId");
			if (!this.userId) {
				this.$router.push("/login");
				return;
			}

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
			} catch (err) {
				console.log("Errore nel caricamento delle conversazioni!");
			}
		},
	},
	created() {
		this.fetchConversations();
	},
};
</script>

<style scoped>
.home-container {
	flex-grow: 1;
	display: flex;
	justify-content: center;
	align-items: center;
	background-color: #d9dbd5;
	height: 100%;
	padding: 20px;
	box-sizing: border-box;
}

.conversation-panel {
	background-color: white;
	width: 100%;
	max-width: 900px;
	height: 85vh;
	border-radius: 12px;
	padding: 20px;
	padding-top: 10px;
	box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
	display: flex;
	flex-direction: column;
	overflow: hidden;
}

.search-bar {
	border: 2px solid transparent;
	width: 100%;
	height: 2.5em;
	padding-left: 0.8em;
	outline: none;
	overflow: hidden;
	background-color: #f3f3f3;
	border-radius: 10px;
	transition: all 0.5s;
}

.search-bar:hover,
.search-bar:focus {
	border: 2px solid #4a9dec;
	box-shadow: 0px 0px 0px 7px rgb(74, 157, 236, 20%);
	background-color: white;
}

.conversation-list {
	flex-grow: 1;
	overflow-y: auto;
	display: flex;
	flex-direction: column;
	gap: 10px;
}

.empty-state {
	margin-top: 2rem;
	text-align: center;
}

.empty-state button {
	margin-top: 1rem;
	padding: 10px 16px;
	border: none;
	border-radius: 8px;
	background-color: #008080;
	color: white;
	cursor: pointer;
}

button {
	padding: 1.3em 3em;
	font-size: 12px;
	font: sans-serif;
	text-transform: uppercase;
	letter-spacing: 2.5px;
	font-weight: 500;
	color: #000;
	background-color: #fff;
	border: none;
	border-radius: 45px;
	box-shadow: 0px 8px 15px rgba(0, 0, 0, 0.1);
	transition: all 0.3s ease 0s;
	cursor: pointer;
	outline: none;
}

button:hover {
	background-color: #23c483;
	box-shadow: 0px 15px 20px rgba(46, 229, 157, 0.4);
	color: #fff;
	transform: translateY(-7px);
}

button:active {
	transform: translateY(-1px);
}

.error-message {
	color: red;
	text-align: center;
	margin-bottom: 10px;
}
</style>
