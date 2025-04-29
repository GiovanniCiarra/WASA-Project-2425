<template>
	<div class="input-container">
		<!-- Immagine -->
		<label class="icon-button">
			<input type="file" hidden @change="handleImageUpload" />
			<img src="../assets/icons/image.png" alt="Image" />
		</label>

		<!-- Campo messaggio -->
		<input
			type="text"
			v-model="message"
			placeholder="Scrivi un messaggio"
			class="text-input"
			@keyup.enter="sendMessage"
		/>

		<!-- Invia -->
		<button class="send-button" @click="sendMessage">➤</button>
	</div>
</template>

<script>
export default {
	props: {
		conversationId: {
			type: [String, Number],
			required: true,
		},
	},
	data() {
		return {
			message: "",
			userId: null,
		};
	},

	created() {
		this.userId = sessionStorage.getItem("userId");
	},

	methods: {
		async sendMessage() {
			if (!this.message.trim()) return;

			try {
				const res = await this.$axios.post(
					`/conversations/${this.conversationId}/messages`,
					{
						content: this.message.trim(),
						type: "text",
					},
					{
						headers: {
							Authorization: `Bearer ${this.userId}`,
						},
					}
				);
				console.log(
					"Invio messaggio:",
					this.message,
					"in conversazione:",
					this.conversationId
				);
				this.message = ""; //Resetto il campo di input
				this.$emit("messageSent", res.data); //Comunico l'aggiunta di un messaggio a ConversationWindow
			} catch (err) {
				console.log("Errore nell'invio del messaggio", err);
			}
		},

		handleImageUpload(event) {
			const file = event.target.files[0];
			if (!file) return;

			console.log("Hai selezionato l'immagine:", file.name);
			// Qui gestirai l'invio dell'immagine
		},
	},
};
</script>

<style scoped>
.input-container {
	display: flex;
	align-items: center;
	padding: 10px 20px;
	background-color: #f0f0f0;
	border-top: 1px solid #ddd;
	width: 100%;
	box-sizing: border-box;
}

.icon-button {
	cursor: pointer;
	font-size: 1.4rem;
	margin-right: 12px;
}

.text-input {
	flex-grow: 1;
	padding: 10px 14px;
	border-radius: 20px;
	border: 1px solid #ccc;
	outline: none;
	font-size: 1rem;
}

.send-button {
	margin-left: 12px;
	background: none;
	border: none;
	font-size: 1.4rem;
	cursor: pointer;
	color: #008080;
}
</style>
