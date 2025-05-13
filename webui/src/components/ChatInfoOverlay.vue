<template>
	<transition name="slide" @after-leave="$emit('closeOverlay')">
		<div class="info-overlay" v-if="visible">
			<img
				@click="visible = false"
				class="close-btn"
				src="../assets/icons/close-big.png"
				alt="Close"
			/>

			<div class="info-content">
				<div class="photo-name-container">
					<img
						:src="chatPhoto || defaultPhoto"
						alt="Foto Chat"
						class="info-photo"
					/>
					<h2 class="info-name">{{ chatName }}</h2>
				</div>
			</div>
		</div>
	</transition>
</template>

<script>
import defaultUserPhoto from "../assets/images/default-user-100.png";
import closeIcon from "../assets/icons/close-big.png";

export default {
	props: {
		chatName: String,
		chatPhoto: String,
	},
	emits: ["closeOverlay"],

	data() {
		return {
			defaultPhoto: defaultUserPhoto,
			closeIcon,
			visible: true,
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
}

.slide-enter-active,
.slide-leave-active {
	transition: all 0.4s ease;
}

.slide-enter-from,
.slide-leave-to {
	transform: translateY(-20px);
	opacity: 0;
}

.slide-enter-to,
.slide-leave-from {
	transform: translateY(0);
	opacity: 1;
}

.close-btn {
	width: 48px;
	height: 48px;
	cursor: pointer;
	opacity: 0.7;
	align-self: flex-end;
}

.close-btn:hover {
	opacity: 1;
}

.info-content {
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
	flex-grow: 1;
	margin-bottom: 70px;
}

.photo-name-container {
	background: rgba(255, 255, 255, 0.7);
	backdrop-filter: blur(10px);
	border-radius: 20px;
	padding: 40px 50px;
	box-shadow: 0 12px 32px rgba(0, 0, 0, 0.25);
	text-align: center;
	max-width: 400px;
	width: 100%;
}

.info-photo {
	width: 160px;
	height: 160px;
	border-radius: 50%;
	object-fit: cover;
	margin-bottom: 20px;
	box-shadow: 0 6px 18px rgba(0, 0, 0, 0.35);
	transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.info-photo:hover {
	transform: scale(1.08);
	box-shadow: 0 10px 28px rgba(0, 0, 0, 0.4);
}

.info-name {
	font-size: 2.2rem;
	font-weight: 700;
	color: #333;
	margin-bottom: 12px;
}
</style>
