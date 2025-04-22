<template>
	<MessageToast :message="toastMessage" />
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
import { RouterView } from "vue-router";
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
		};
	},

	methods: {
		showToast(msg) {
			this.toastMessage = "";
			this.$nextTick(() => {
				this.toastMessage = msg;
			});
		},
	},
};
</script>

<style>
.background {
	background: repeating-linear-gradient(
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
</style>
