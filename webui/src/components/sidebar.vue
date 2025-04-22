<template>
	<div class="sidebar">
		<button @click="$router.push('/Home')">
			<img src="../assets/icons/conversations.png" alt="Conversations" />
		</button>
		<button @click="$router.push('/NewChat')">
			<img src="../assets/icons/chat.png" alt="NewChat" />
		</button>
		<button @click="$router.push('/NewGroup')">
			<img src="../assets/icons/group.png" alt="NewGroup" />
		</button>
		<button @click="$router.push('/Profile')">
			<img src="../assets/icons/profile.png" alt="Profile" />
		</button>
		<button @click="refresh">
			<img src="../assets/icons/refresh.png" alt="Refresh" />
		</button>
		<button @click="logout">
			<img src="../assets/icons/logout.png" alt="Logout" />
		</button>
	</div>
</template>

<script>
export default {
	methods: {
		refresh() {
			location.reload();
		},
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
}
button {
	margin: 10px 0;
	font-size: 1.2rem;
	background: none;
	border: none;
	cursor: pointer;
}

button:hover {
	opacity: 1;
	transform: scale(1.1);
}
</style>
