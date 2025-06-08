<template>
	<div class="popup-overlay">
		<div class="popup-content">
			<!-- Icona di chiusura popup -->
			<img
				src="../assets/icons/close-big.png"
				alt="Chiudi"
				class="close-popup-btn"
				@click="$emit('close')"
			/>

			<h2 class="popup-title">Aggiungi Partecipanti</h2>

			<!-- Lista degli utenti da poter aggiungere al gruppo -->
			<div class="user-list">
				<div
					v-for="user in users"
					:key="user.userId"
					class="user-item"
					@click="toggleSelection(user)"
					:class="{ selected: isSelected(user) }"
				>
					<img
						:src="user.photo || defaultUserPhoto"
						alt="Foto"
						class="user-photo"
					/>
					<span class="username">{{ user.username }}</span>
				</div>

				<!-- Label in caso di assenza di utenti -->
				<div v-if="users.length === 0" class="no-users-msg">
					Nessun utente disponibile da aggiungere
				</div>
			</div>

			<!-- Confirm button -->
			<button
				class="confirm-btn"
				@click="confirmSelection"
				:disabled="!canConfirm"
			>
				Conferma
			</button>
		</div>
	</div>
</template>

<script>
import defaultUserPhoto from "../assets/images/default-user-48.png";

export default {
	props: {
		users: Array,
	},

	emits: ["close", "confirm"],

	data() {
		return {
			defaultUserPhoto,
			selectedUsers: [], // Utenti selezionati da aggiungere
		};
	},

	computed: {
		// Abilita/disabilita il bottone di aggiunta al gruppo
		canConfirm() {
			return this.selectedUsers.length > 0;
		},
	},

	methods: {
		// Verifica se un utente sia selezionato o meno
		isSelected(user) {
			return this.selectedUsers.some((u) => u.userId === user.userId);
		},

		// Slettore degli utenti
		toggleSelection(user) {
			if (this.selectedUsers.some((u) => u.userId === user.userId)) {
				this.selectedUsers = this.selectedUsers.filter(
					(u) => u.userId !== user.userId
				);
			} else {
				this.selectedUsers = [...this.selectedUsers, user];
			}
		},

		// Conferma l'aggiunta
		confirmSelection() {
			this.$emit("confirm", this.selectedUsers);
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
	box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.username {
	font-size: 1rem;
}

.no-users-msg {
	text-align: center;
	font-size: 1rem;
	color: #888;
	padding: 20px;
	font-style: italic;
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

.confirm-btn:disabled {
	background-color: #a5d6c5;
	cursor: not-allowed;
}
</style>
