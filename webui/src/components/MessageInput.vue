<template>
	<div v-if="replyingTo" class="reply-preview-block">
		<div class="reply-block">
			<strong>{{ replyingTo.senderUsername }}</strong>
			<div class="reply-content">
				<span v-if="replyingTo.type === 'text'">
					{{ replyingTo.content }}
				</span>
				<span v-else class="snippet-with-icon">
					<img :src="imageSnippet" alt="Foto" class="snippet-icon" />
					<span>Image</span>
				</span>
			</div>
		</div>
		<button class="cancel-reply-btn" @click="$emit('cancelReply')">
			<img :src="closeIcon" alt="Annulla" />
		</button>
	</div>

	<div class="input-container">
		<!-- Immagine -->
		<div class="icon-button" @click="$refs.fileInput.click()">
			<input
				type="file"
				ref="fileInput"
				hidden
				@change="handleImageSelect"
				accept="image/png"
			/>
			<img src="../assets/icons/image.png" alt="Image" />
		</div>

		<!-- Campo messaggio -->
		<div class="input-wrapper">
			<input
				type="text"
				v-model="textMessage"
				placeholder="Scrivi un messaggio"
				class="text-input"
				@keyup.enter="handleSend"
				:disabled="selectedImage"
			/>

			<!-- Pillola Immagine -->
			<div v-if="selectedImage" class="image-pill">
				{{ selectedImage.name }}
				<span class="remove-image" @click="removeImage">
					<img :src="closeIcon" alt="close" />
				</span>
			</div>
		</div>

		<!-- Invia -->
		<button class="send-button" @click="handleSend">➤</button>
	</div>
</template>

<script>
import closeIcon from "../assets/icons/close.png";
import imageSnippet from "../assets/icons/image-snippet.png";

export default {
	props: {
		conversationId: {
			type: [String, Number],
			required: true,
		},
		replyingTo: Object,
	},

	emits: ["messageSent", "cancelReply"],

	data() {
		return {
			textMessage: "",
			userId: null,
			selectedImage: null,
			closeIcon,
			imageSnippet,
		};
	},

	created() {
		this.userId = sessionStorage.getItem("userId");
	},

	methods: {
		handleImageSelect(event) {
			const file = event.target.files[0];
			if (file) this.selectedImage = file;
		},

		removeImage() {
			this.selectedImage = null;
			if (this.$refs.fileInput) this.$refs.fileInput.value = "";
		},

		async handleSend() {
			if (this.selectedImage) {
				const reader = new FileReader();
				reader.onload = async () => {
					if (this.replyingTo) {
						await this.replyMessage(reader.result, "photo");
					} else {
						await this.sendMessage(reader.result, "photo");
					}
				};
				reader.readAsDataURL(this.selectedImage);
			} else if (this.textMessage.trim()) {
				if (this.replyingTo) {
					await this.replyMessage(this.textMessage.trim(), "text");
				} else {
					await this.sendMessage(this.textMessage.trim(), "text");
				}
			}
			this.resetInput();
		},

		resetInput() {
			this.textMessage = "";
			this.selectedImage = null;
			if (this.$refs.fileInput) this.$refs.fileInput.value = "";
		},

		async sendMessage(content, type) {
			try {
				const res = await this.$axios.post(
					`/conversation-messages/${this.conversationId}`,
					{
						content: content,
						type: type,
					},
					{
						headers: { Authorization: `Bearer ${this.userId}` },
					}
				);
				this.$emit("messageSent", res.data);
			} catch (err) {
				console.error("Errore nell'invio messaggio", err);
				this.handleError(err);
			}
		},

		async replyMessage(content, type) {
			try {
				const res = await this.$axios.post(
					`/messages/${this.replyingTo.messageId}/replies`,
					{
						content: content,
						type: type,
					},
					{
						headers: { Authorization: `Bearer ${this.userId}` },
					}
				);
				this.$emit("cancelReply");
				this.$emit("messageSent", res.data);
			} catch (err) {
				console.error("Errore nell'invio messaggio", err);
				this.handleError(err);
			}
		},

		handleError(err) {
			if (!err.response) {
				this.$root.showToast(
					"Errore di rete: controlla la connessione.",
					"error"
				);
				return;
			}

			const status = err.response.status;
			const detail = err.response.data || "ERRORE";

			switch (status) {
				case 400:
					this.$root.showToast(`Errore 400: ${detail}`, "error");
					break;
				case 401:
					this.$root.showToast(
						"Errore 401: Non autorizzato, effettua di nuovo il login",
						"error"
					);
					this.$router.push("/Login");
					break;
				case 403:
					this.$root.showToast(`Errore 403: ${detail}`, "error");
					break;
				case 500:
					this.$root.showToast(`Errore 500:\n ${detail}!`, "error");
					break;
				default:
					this.$root.showToast("Errore sconosciuto", "error");
			}
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

.input-wrapper {
	position: relative;
	flex-grow: 1;
}

.text-input {
	width: 100%;
	padding: 10px 14px;
	border-radius: 20px;
	border: 1px solid #ccc;
	outline: none;
	font-size: 1rem;
}

.image-pill {
	position: absolute;
	top: 50%;
	left: 12px;
	transform: translateY(-50%);
	background-color: #f1f1f1;
	color: #333;
	padding: 5px 10px;
	border-radius: 20px;
	font-size: 0.75rem;
	display: flex;
	align-items: center;
	max-width: 60%;
	min-width: 18%;
	overflow: hidden;
	white-space: nowrap;
	text-overflow: ellipsis;
	justify-content: space-between;
	box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
	border: 1px solid #ddd;
}

.remove-image {
	margin-left: 8px;
	cursor: pointer;
}

.send-button {
	margin-left: 12px;
	background: none;
	border: none;
	font-size: 1.4rem;
	cursor: pointer;
	color: #008080;
}

.reply-preview-block {
	display: flex;
	align-items: center;
	justify-content: space-between;
	background-color: #e0f7fa;
	border-left: 4px solid #26c6da;
	padding: 6px 10px;
	margin: 0 20px 8px 20px;
	border-radius: 8px;
	font-size: 0.85rem;
	color: #333;
	box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.cancel-reply-btn {
	background: none;
	border: none;
	cursor: pointer;
	color: #777;
	font-size: 1rem;
	margin-left: 8px;
}

.cancel-reply-btn:hover {
	color: #d32f2f;
}

.reply-preview-block .reply-block {
	background: none;
	border: none;
	padding: 0;
	margin: 0;
	border-radius: 0;
	box-shadow: none;
}

.reply-content {
	margin-top: 4px;
	font-style: italic;
	font-size: 0.8rem;
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
	max-width: 400px;
}

.snippet-with-icon {
	display: inline-flex;
	align-items: center;
	gap: 4px;
}

.snippet-icon {
	width: 16px;
	height: 16px;
}
</style>
