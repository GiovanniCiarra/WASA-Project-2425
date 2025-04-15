<template>
	<div class="sidebar">
		<button @click="$router.push('/Home')">🏠</button>
		<button @click="$router.push('/NewChat')">💬</button>
		<button @click="$router.push('/NewGroup')">👥</button>
		<button @click="refresh">🔄</button>
		<button @click="logout">🚪</button>
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
</style>
