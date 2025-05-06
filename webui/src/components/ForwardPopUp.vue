<template>
	<div class="popup-overlay">
		<div class="popup-content">
			<img
				src="../assets/icons/close-big.png"
				alt="Chiudi"
				class="close-popup-btn"
				@click="$emit('cancel')"
			/>

			<h2 class="popup-title">Inoltra a:</h2>

			<div class="user-list">
				<li
					v-for="conv in conversations"
					:key="conv.conversationId"
					@click="forwardToConversation(conv)"
					class="user-item"
				>
					<img
						:src="conv.photo || defaultPhoto"
						alt="photo"
						class="user-photo"
					/>
					<span class="username">{{ conv.name }}</span>
				</li>

				<li
					v-for="user in users"
					:key="user.userId"
					@click="forwardToUser(user)"
					class="user-item"
				>
					<img
						:src="user.photo || defaultPhoto"
						alt="user-photo"
						class="user-photo"
					/>
					<span class="username">{{ user.username }}</span>
				</li>
			</div>
		</div>
	</div>
</template>

<script>
import defaultPhoto from "../assets/images/default-user-48.png";
export default {
	props: {
		message: Object,
		userId: Number,
	},

	data() {
		return {
			defaultPhoto,
			conversations: [],
			users: [],
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
			this.conversations = this.conversations.filter(
				(conv) => conv.type === "group"
			);
		} catch (err) {
			console.error("Errore nel recupero conversazioni", err);
		}

		try {
			const res = await this.$axios.get("/users", {
				headers: { Authorization: `Bearer ${this.userId}` },
			});
			this.users = res.data;
		} catch (err) {
			console.error("Errore nel recupero degli utenti", err);
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

		async forwardToUser(targetUser) {
			try {
				// Verifico se già esiste una conversazione con questo utente
				const res = await this.$axios.get(
					`/privateChat/${targetUser.userId}`,
					{
						headers: { Authorization: `Bearer ${this.userId}` },
					}
				);

				// Se la conversazione esiste, allora ottengo il conversationId, inoltro il messaggio e cambio View
				const conversationId = res.data.conversationId;

				try {
					await this.$axios.post(
						`/messages/${this.message.messageId}/forward`,
						{
							targetChatId: conversationId,
						},
						{
							headers: { Authorization: `Bearer ${this.userId}` },
						}
					);
					this.$emit("forwarded");
					this.$router.push({
						path: `/conversations/${conversationId}`,
						query: {
							type: "chat",
							name: targetUser.username,
							photo: targetUser.photo,
						},
					});
				} catch (err) {
					console.error("Errore durante inoltro", err);
				}
			} catch (err) {
				if (err.response && err.response.status === 404) {
					// Se la conversazione non esiste, ricevo un error 404, allora la creo, inoltro il messaggio e cambio view
					try {
						const createRes = await this.$axios.post(
							"/conversations",
							{
								name: "",
								photo: "",
								members: [targetUser.userId],
							},
							{
								headers: {
									Authorization: `Bearer ${this.userId}`,
								},
							}
						);

						const conversationId = createRes.data.conversationId;

						try {
							await this.$axios.post(
								`/messages/${this.message.messageId}/forward`,
								{
									targetChatId: conversationId,
								},
								{
									headers: {
										Authorization: `Bearer ${this.userId}`,
									},
								}
							);
							this.$emit("forwarded");
							this.$router.push({
								path: `/conversations/${conversationId}`,
								query: {
									type: "chat",
									name: targetUser.username,
									photo: targetUser.photo,
								},
							});
						} catch (err) {
							console.error("Errore durante inoltro", err);
						}
						this.$router.push({
							path: `/conversations/${conversationId}`,
							query: {
								type: "chat",
								name: targetUser.username,
								photo: targetUser.photo,
							},
						});
					} catch (err) {
						console.error(
							"Errore nella creazione della conversazione",
							err
						);
					}
				} else {
					console.error(
						"Errore nella creazione/recupero conversazione",
						err
					);
				}
			}
		},
	},
};
</script>

<style>
.popup-overlay {
	position: fixed;
	top: 0;
	left: 0;
	right: 0;
	bottom: 0;
	background-color: rgba(0, 0, 0, 0.5);
	display: flex;
	justify-content: center;
	align-items: center;
	z-index: 1000;
}

.popup-content {
	background: white;
	border-radius: 20px;
	padding: 30px;
	width: 90%;
	max-width: 500px;
	position: relative;
}

.close-popup-btn {
	position: absolute;
	top: 16px;
	right: 16px;
	width: 24px;
	height: 24px;
	cursor: pointer;
	opacity: 0.7;
}

.close-popup-btn:hover {
	opacity: 1;
}

.popup-title {
	font-size: 1.8rem;
	font-weight: bold;
	margin-bottom: 20px;
	text-align: center;
}

.user-list {
	max-height: 300px;
	overflow-y: auto;
	display: flex;
	flex-direction: column;
	gap: 12px;
	margin-bottom: 20px;
}

.user-item {
	display: flex;
	align-items: center;
	padding: 10px;
	border-radius: 3px;
	cursor: pointer;
	transition: background 0.2s;
	border: 1px solid transparent;
	border-bottom: 1px solid #73bfa1;
}

.user-item:hover {
	background-color: rgba(0, 128, 128, 0.1);
	box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.user-photo {
	width: 40px;
	height: 40px;
	border-radius: 50%;
	object-fit: cover;
	margin-right: 12px;
	box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.username {
	font-size: 1rem;
}
</style>
