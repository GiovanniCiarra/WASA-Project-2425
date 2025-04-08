<script>
import ConversationList from "../components/ConversationsList.vue";
import ConversationWindow from "../components/ConversationWindow.vue";

export default {
	components: { ConversationList, ConversationWindow },

	data: function () {
		return {
			conversations: [],
			selectedConversation: null,
			userId: null,
		};
	},

	async created() {
		this.userId = sessionStorage.getItem("userId"); //Recupero lo userId dal sessionStorage
		if (!userId) {
			console.log("userId non trovato nel sessionStorage");
			return;
		}

		try {
			//Fetch delle conversazioni
			let conversationsResponse = await this.$axios.get(
				"http://localhost:3000/conversations",
				{
					headers: {
						Authorization: `Bearer ${this.userId}`,
					},
				}
			);
			this.conversations = conversationsResponse.data; //Salvo le conversazioni
		} catch (e) {
			if (e.response) {
				// Il server risponde con un codice di stato
				switch (e.response.status) {
					case 401:
						this.conversationError = "Errore 401: Unauthorized!";
						break;
					case 500:
						this.conversationError =
							"Errore 500: Errore interno del server!";
						break;
					default:
						this.conversationError = `Errore ${
							e.response.status
						}: ${
							e.response.data.message ||
							"Si è verificato un errore!"
						}`;
				}
			} else {
				this.conversationError = "Si è verificato un errore!";
			}
		}
	},

	methods: {
		async refresh() {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/");
				this.some_data = response.data;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
	},
	mounted() {
		this.refresh();
	},
};
</script>

<template>
	<div class="container">
		<ConversationList
			:conversations="conversations"
			@conversationSelected="goToConversation"
		/>
		<ConversationWindow :conversation="selectedConversation" />
	</div>
</template>

<stylec scooped></stylec>
