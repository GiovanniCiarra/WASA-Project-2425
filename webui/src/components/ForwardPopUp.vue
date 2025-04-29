<template>
	<div class="forward-popup-overlay">
		<div class="forward-popup">
			<h3>Inoltra a:</h3>
			<ul>
				<li
					v-for="conv in conversations"
					:key="conv.conversationId"
					@click="forwardToConversation(conv)"
				>
					<img
						:src="conv.photo || '/public/images/default-user.png'"
						alt="photo"
					/>
					<span>{{ conv.name }}</span>
				</li>
			</ul>
			<button @click="$emit('cancel')" class="cancel-btn">Annulla</button>
		</div>
	</div>
</template>

<script>
export default {
	props: {
		message: Object,
		userId: Number,
	},

	data() {
		return {
			conversations: [],
		};
	},

	async created() {
		try {
			const res = await this.$axios.get("/conversations", {
				headers: { Authorization: `Bearer ${this.userId}` },
			});
			this.conversations = res.data.filter(
				(conv) => conv.conversationId !== this.message.conversationId
			);
		} catch (err) {
			console.error("Errore nel recupero conversazioni", err);
		}
	},

	methods: {
		async forwardToConversation(targetConv) {
			try {
				await this.$axios.post(
					`/messages/${this.message.messageId}/forward`,
					{
						targetChatId: targetConv.conversationId,
					},
					{
						headers: { Authorization: `Bearer ${this.userId}` },
					}
				);
				this.$emit("forwarded");
				this.$router.push({
					path: `/conversations/${targetConv.conversationId}`,
					query: {
						type: targetConv.type,
						name: targetConv.name,
						photo: targetConv.photo,
					},
				});
			} catch (err) {
				console.error("Errore durante inoltro", err);
			}
		},
	},
};
</script>

<style scoped>
.forward-popup-overlay {
	position: fixed;
	inset: 0;
	background: rgba(0, 0, 0, 0.3);
	display: flex;
	justify-content: center;
	align-items: center;
	z-index: 9999;
}

.forward-popup {
	background: white;
	padding: 20px;
	border-radius: 10px;
	width: 320px;
	max-height: 80vh;
	overflow-y: auto;
	box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
}

ul {
	list-style: none;
	padding: 0;
	margin: 10px 0;
}

li {
	display: flex;
	align-items: center;
	padding: 10px;
	cursor: pointer;
	border-bottom: 1px solid #eee;
	gap: 10px;
}

li:hover {
	background-color: #f5f5f5;
}

img {
	width: 36px;
	height: 36px;
	border-radius: 50%;
	object-fit: cover;
}

.cancel-btn {
	margin-top: 12px;
	background: none;
	border: 1px solid #999;
	padding: 6px 12px;
	border-radius: 6px;
	cursor: pointer;
}
</style>
