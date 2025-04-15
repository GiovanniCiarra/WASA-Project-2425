<template>
	<div class="header-container" @click="handleClick">
		<img
			:src="photo || '/images/default-user.png'"
			alt="Foto"
			class="header-photo"
		/>
		<span class="header-name">{{ name }}</span>
	</div>
</template>

<script>
export default {
	props: {
		name: {
			type: String,
			required: true,
		},
		photo: {
			type: String,
			default: "",
		},
		type: {
			type: String,
			required: true,
		},
	},
	methods: {
		handleClick() {
			if (this.type === "group") {
				this.$router.push(
					`/groups/${this.$route.params.conversationId}/info`
				);
			} else {
				// eventualmente niente o profilo utente
				console.log("Conversazione privata - nessuna azione");
			}
		},
	},
};
</script>

<style scoped>
.header-container {
	display: flex;
	align-items: center;
	padding: 12px 20px;
	background-color: #f8f8f8;
	border-bottom: 1px solid #ddd;
	cursor: pointer;
}

.header-photo {
	width: 40px;
	height: 40px;
	border-radius: 50%;
	object-fit: cover;
	margin-right: 12px;
}

.header-name {
	font-size: 1.1rem;
	font-weight: 600;
	color: #333;
}
</style>
