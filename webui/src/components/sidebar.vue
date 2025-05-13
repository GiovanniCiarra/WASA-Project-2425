<template>
	<div class="sidebar">
		<!-- Bottoni in alto, per navigazione principale -->
		<div class="top-buttons">
			<button
				:class="{ active: currentRouteName === 'Home' }"
				@click="$router.push('/Home')"
			>
				<img
					src="../assets/icons/conversations.png"
					alt="Conversations"
				/>
			</button>
			<button
				:class="{ active: currentRouteName === 'newChat' }"
				@click="$router.push('/NewChat')"
			>
				<img src="../assets/icons/chat.png" alt="NewChat" />
			</button>
			<button
				:class="{ active: currentRouteName === 'newGroup' }"
				@click="$router.push('/NewGroup')"
			>
				<img src="../assets/icons/group.png" alt="NewGroup" />
			</button>
			<button
				:class="{ active: currentRouteName === 'Profile' }"
				@click="$router.push('/Profile')"
			>
				<img src="../assets/icons/profile.png" alt="Profile" />
			</button>
		</div>

		<!-- Bottoni in basso, Refresh e Logout -->
		<div class="bottom-buttons">
			<button @click="refresh">
				<img src="../assets/icons/refresh.png" alt="Refresh" />
			</button>
			<button @click="logout">
				<img src="../assets/icons/logout.png" alt="Logout" />
			</button>
		</div>
	</div>
</template>

<script>
export default {
	computed: {
		// Prendo il nome della route attiva
		currentRouteName() {
			return this.$route.name;
		},
	},

	methods: {
		// Refresh della pagina
		refresh() {
			location.reload();
		},

		// Clear del sessionStorage e push al LoginView
		logout() {
			sessionStorage.removeItem("userId");
			this.$router.push("/Login");
		},
	},
};
</script>

<style scoped>
.sidebar {
	width: 80px;
	background-color: #f0f0f0;
	display: flex;
	flex-direction: column;
	flex-shrink: 0;
	align-items: center;
	padding: 10px 0;
	border-right: 2px solid #ccc;
	height: 100%;
}

.top-buttons {
	display: flex;
	flex-direction: column;
	align-items: center;
}

.bottom-buttons {
	margin-top: auto;
	display: flex;
	flex-direction: column;
	align-items: center;
}

button {
	width: 50px;
	height: 50px;
	border: none;
	border-radius: 14px;
	background-color: transparent;
	display: flex;
	align-items: center;
	justify-content: center;
	margin: 10px 0;
	cursor: pointer;
	transition: all 0.3s ease;
}

button img {
	width: 32px;
	height: 32px;
	transition: filter 0.3s ease;
}

button:hover {
	opacity: 1;
	transform: scale(1.1);
}

button.active {
	background-color: rgba(47, 162, 122, 0.15);
	border: 1px solid #2fa27a;
	box-shadow: 0 2px 6px rgba(47, 162, 122, 0.25);
}
</style>
