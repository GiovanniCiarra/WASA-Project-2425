<template>
	<div class="popup-overlay">
		<div class="popup-content">
			<img
				src="../assets/icons/close-big.png"
				alt="Chiudi"
				class="close-popup-btn"
				@click="$emit('close')"
			/>

			<h2 class="popup-title">Aggiungi Partecipanti</h2>

			<div class="user-list">
				<div
					v-for="user in users"
					:key="user.userId"
					class="user-item"
					@click="toggleSelection(user.userId)"
					:class="{ selected: selectedUserIds.includes(user.userId) }"
				>
					<img
						:src="user.photo || defaultUserPhoto"
						alt="Foto"
						class="user-photo"
					/>
					<span class="username">{{ user.username }}</span>
				</div>
			</div>

			<button class="confirm-btn" @click="confirmSelection">
				Conferma
			</button>
		</div>
	</div>
</template>

<script>
import defaultUserPhoto from "../assets/images/default-user.png";

export default {
	props: {
		users: Array,
	},
	emits: ["close", "confirm"],
	data() {
		return {
			defaultUserPhoto,
			selectedUserIds: [],
		};
	},
	methods: {
		toggleSelection(userId) {
			if (this.selectedUserIds.includes(userId)) {
				this.selectedUserIds = this.selectedUserIds.filter(
					(id) => id !== userId
				);
			} else {
				this.selectedUserIds.push(userId);
			}
		},
		confirmSelection() {
			this.$emit("confirm", this.selectedUserIds);
		},
	},
};
</script>

<style scoped>
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
	border-radius: 10px;
	cursor: pointer;
	transition: background 0.2s;
}

.user-item:hover {
	background: #f2f2f2;
}

.user-item.selected {
	background: #2fa27a;
	color: white;
}

.user-photo {
	width: 40px;
	height: 40px;
	border-radius: 50%;
	object-fit: cover;
	margin-right: 12px;
}

.username {
	font-size: 1rem;
}

.confirm-btn {
	width: 100%;
	padding: 12px;
	background-color: #2fa27a;
	color: white;
	border: none;
	border-radius: 10px;
	cursor: pointer;
	font-size: 1rem;
}

.confirm-btn:hover {
	background-color: #27986d;
}
</style>
