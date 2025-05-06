<template>
	<div
		class="member-item"
		:class="{ 'current-user': member.userId === currentUser }"
	>
		<img
			:src="member.photo || defaultUserPhoto"
			alt="Foto Membro"
			class="member-photo"
		/>
		<span class="member-username"
			>{{ member.username }}
			<span v-if="member.userId === currentUser" class="you-label"
				>(tu)</span
			>
		</span>
		<button
			v-if="member.userId === currentUser"
			class="leave-btn"
			@click="$emit('leaveGroup')"
		>
			Abbandona
		</button>
	</div>
</template>

<script>
import defaultUserPhoto from "../assets/images/default-user-48.png";

export default {
	props: {
		member: Object,
		currentUser: Number,
	},
	emits: ["leaveGroup"],
	data() {
		return {
			defaultUserPhoto,
		};
	},
};
</script>

<style scoped>
.member-item {
	display: flex;
	align-items: center;
	margin-bottom: 4px;
	padding: 8px;
	border: 1px solid #e0e0e0;
	border-radius: 8px;
	transition: background-color 0.2s, border-color 0.2s;
}

.member-item:hover {
	background-color: #f9f9f9;
	border-color: #c0c0c0;
}

.member-photo {
	width: 44px;
	height: 44px;
	border-radius: 50%;
	object-fit: cover;
	margin-right: 12px;
	box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.member-username {
	flex-grow: 1;
	font-weight: 500;
}

.you-label {
	font-size: 0.8rem;
	color: #888;
	margin-left: 6px;
}

.leave-btn {
	padding: 6px 12px;
	background-color: rgb(255, 228, 228); /* Sfondo bianco */
	color: #ff5b42; /* Testo arancione */
	border: 2px solid #ff5b42; /* Bordo arancione */
	border-radius: 6px;
	cursor: pointer;
	font-size: 0.9rem;
	transition: background-color 0.3s, color 0.3s;
}

.leave-btn:hover {
	background-color: #ff5b42; /* Sfondo si colora */
	color: white; /* Testo diventa bianco */
}
</style>
