<template>
	<teleport to="body">
		<transition name="slide-up-fade">
			<div v-if="visible" class="app-toast">
				<div class="toast-content">
					<img
						class="toast-icon-inline"
						src="../assets/icons/warning-yellow.png"
						alt="icon"
					/>
					<div class="toast-text">{{ message }}</div>
				</div>
				<div class="accent-bar"></div>
			</div>
		</transition>
	</teleport>
</template>

<script>
export default {
	props: {
		message: String,
	},
	data() {
		return {
			visible: false,
		};
	},
	watch: {
		message(newVal) {
			if (newVal) {
				this.visible = true;
				setTimeout(() => {
					this.visible = false;
				}, 5000);
			}
		},
	},
};
</script>

<style>
.app-toast {
	position: fixed;
	bottom: 20px;
	right: 200px;
	background-color: white;
	color: #333;
	padding: 12px 20px;
	border-radius: 8px;
	font-size: 0.95rem;
	font-family: sans-serif;
	box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
	z-index: 99999;
	min-width: 240px;
	max-width: 320px;
	overflow: hidden;
	display: flex;
	flex-direction: column;
}

.toast-content {
	display: flex;
	flex-direction: row;
	align-items: flex-start;
	gap: 10px;
}

.toast-icon-inline {
	width: 18px;
	height: 18px;
	flex-shrink: 0;
	object-fit: contain;
	margin-top: 2px;
}

.toast-text {
	flex: 1;
	overflow-wrap: break-word;
	word-break: break-word;
	line-height: 1.4;
}

/* Riga decorativa */
.accent-bar {
	height: 4px;
	width: 100%;
	background-color: #2fa27a;
	margin-top: 10px;
	border-bottom-left-radius: 8px;
	border-bottom-right-radius: 8px;
}

/* Animazione */
.slide-up-fade-enter-active,
.slide-up-fade-leave-active {
	transition: opacity 0.6s ease, transform 0.6s ease;
}
.slide-up-fade-enter-from,
.slide-up-fade-leave-to {
	opacity: 0;
	transform: translateY(40px);
}
.slide-up-fade-enter-to,
.slide-up-fade-leave-from {
	opacity: 1;
	transform: translateY(0);
}
</style>
