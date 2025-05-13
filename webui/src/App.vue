<template>
	<!-- Toast per messaggi di errore, warning, successo -->
	<MessageToast :message="toastMessage" :type="toastType" />

	<div class="background">
		<!-- Sidebar visibile ovunque tranne nella Login -->
		<div v-if="$route.name !== 'Login'" class="app-container">
			<Sidebar />

			<!-- Spazio per tutte le view dinamiche -->
			<main class="main-content">
				<RouterView />
			</main>
		</div>

		<!-- Se il LoginView, ho un contenitore diverso -->
		<div v-else class="login-container">
			<RouterView />
		</div>
	</div>
</template>

<script>
import Sidebar from "../src/components/sidebar.vue";
import MessageToast from "../src/components/messageToast.vue";

export default {
	components: {
		Sidebar,
		MessageToast,
	},

	data() {
		return {
			toastMessage: "",
			toastType: "success",
		};
	},

	methods: {
		// Metodo per mostrare un toast
		showToast(msg, type) {
			this.toastMessage = "";
			this.$nextTick(() => {
				this.toastMessage = msg;
				this.toastType = type;
			});
		},
	},
};
</script>

<style>
.background {
	background-image: url("../src/assets/images/sfondo.jpg");
	background-size: cover;
	background-position: center;
	background-repeat: no-repeat;
	width: 100%;
	height: 100vh;
	display: flex;
	justify-content: center;
	align-items: center;
	overflow: hidden;
}

.app-container {
	display: flex;
	width: 90%;
	max-width: 1300px;
	height: 90vh;
	background-color: white;
	border-radius: 16px;
	box-shadow: 0 0 20px rgba(0, 0, 0, 0.2);
	overflow: hidden;
}

.main-content {
	flex-grow: 1;
	height: 100%;
	overflow: hidden;
}

.login-container {
	width: 100%;
	height: 100%;
	display: flex;
	justify-content: center;
	align-items: center;
}
</style>
