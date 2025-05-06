<template>
	<MessageToast :message="toastMessage" :type="toastType" />
	<div class="background">
		<div class="app-container">
			<!-- Sidebar visibile ovunque tranne nella Login -->
			<Sidebar v-if="$route.name !== 'Login'" />

			<!-- Spazio per tutte le view dinamiche -->
			<main class="main-content">
				<RouterView />
			</main>
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
	background-size: cover; /* Ritaglia e adatta l’immagine al contenitore */
	background-position: center; /* Centra l’immagine */
	background-repeat: no-repeat; /* Non ripete l’immagine */
	width: 100%;
	height: 100vh; /* altezza 100% viewport */
	display: flex;
	justify-content: center;
	align-items: center;
	overflow: hidden;
	/*background: repeating-linear-gradient(
		45deg,
		#008080,
		#008080 10px,
		#20b2aa 10px,
		#20b2aa 20px
	);
	min-height: 100vh;
	display: flex;
	justify-content: center;
	align-items: center;
	overflow: hidden;*/
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
</style>
