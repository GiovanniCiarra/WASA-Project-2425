<template>
	<div class="info-overlay">
		<button class="close-btn" @click="$emit('closeOverlay')">✖</button>

		<div class="info-content">
			<img
				:src="chatPhoto || defaultPhoto"
				alt="Foto Chat"
				class="info-photo"
			/>
			<h2 class="info-name">{{ chatName }}</h2>
		</div>
	</div>
</template>

<script>
import defaultUserPhoto from "../assets/images/default-user-100.png";

export default {
	props: {
		chatName: String,
		chatPhoto: String,
	},
	emits: ["closeOverlay"],

	data() {
		return {
			defaultPhoto: defaultUserPhoto,
		};
	},
};
</script>

<style scoped>
.info-overlay {
	background: white;
	padding: 20px;
	overflow-y: auto;
	flex-grow: 1;
	display: flex;
	flex-direction: column;
	animation: slideDown 0.4s ease forwards;
}

@keyframes slideDown {
	from {
		transform: translateY(-20px);
		opacity: 0;
	}
	to {
		transform: translateY(0);
		opacity: 1;
	}
}

.close-btn {
	align-self: flex-end;
	font-size: 1.5rem;
	background: none;
	border: none;
	cursor: pointer;
}

.info-content {
	display: flex;
	flex-direction: column;
	align-items: center;
}

.info-photo {
	width: 120px;
	height: 120px;
	border-radius: 50%;
	object-fit: cover;
	margin-bottom: 10px;
}

.info-name {
	font-size: 1.8rem;
	margin-bottom: 20px;
}
</style>
