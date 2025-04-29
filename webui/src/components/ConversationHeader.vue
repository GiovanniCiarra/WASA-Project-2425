<template>
	<div class="header-container" @click="$emit('openInfo')">
		<img
			:src="photo || defaultConversationPhoto"
			alt="Foto"
			class="header-photo"
		/>
		<span class="header-name">{{ name }}</span>
	</div>
</template>

<script>
import defaultConversationPhoto from "../assets/images/default-user.png";

export default {
	props: {
		name: String,
		photo: String,
		type: String,
	},

	emits: ["openInfo"],

	data() {
		return {
			defaultConversationPhoto,
		};
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
