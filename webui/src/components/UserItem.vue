<template>
	<div class="user-item" :class="{ selected }">
		<img
			:src="user.photo || defaultUserPhoto"
			alt="Foto profilo"
			class="user-photo"
		/>
		<span class="user-name">{{ user.username }}</span>

		<!-- Icona rimozione se selezionato -->
		<img
			v-if="selected"
			class="remove-icon"
			src="../assets/icons/close.png"
			alt="Rimuovi"
		/>
	</div>
</template>

<script>
import defaultUserPhoto from "../assets/images/default-user-48.png";

export default {
	props: {
		user: {
			type: Object,
			required: true,
		},
		selected: Boolean,
	},

	data() {
		return {
			defaultUserPhoto,
		};
	},
};
</script>

<style scoped>
.user-item {
	display: flex;
	align-items: center;
	padding: 12px 16px;
	border: none;
	background-color: white;
	border-radius: 8px;
	width: 100%;
	text-align: left;
	cursor: pointer;
	transition: background-color 0.2s ease;
	gap: 10px;
	font: inherit;
	color: inherit;
	appearance: none;
}

.user-item:hover {
	background-color: #f3f3f3;
}

.user-item.selected {
	background-color: #c0ede1;
}

.user-photo {
	width: 42px;
	height: 42px;
	border-radius: 50%;
	object-fit: cover;
	box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.user-name {
	font-size: 1rem;
	font-weight: 500;
	color: #222;
	flex-grow: 1;
}

.remove-icon {
	width: 20px;
	height: 20px;
	object-fit: contain;
	pointer-events: none;
}
</style>
