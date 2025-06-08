<template>
	<teleport to="body">
		<!-- Transizione -->
		<transition name="slide-up-fade">
			<!-- Stile in base alla classe-->
			<div v-if="visible" :class="['app-toast', typeClass]">
				<div class="toast-content">
					<img class="toast-icon" :src="iconSrc" alt="icon" />
					<div class="toast-text">{{ message }}</div>
				</div>
				<div class="accent-bar" :class="typeClass"></div>
			</div>
		</transition>
	</teleport>
</template>

<script>
import successIcon from "../assets/icons/success-green.png";
import warningIcon from "../assets/icons/warning-yellow.png";
import errorIcon from "../assets/icons/error-red.png";

export default {
	props: {
		message: String,
		type: {
			type: String,
			default: "success", // [success, warning, error]
		},
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
				}, 3000);
			}
		},
	},

	computed: {
		// Decide la classe in base al tipo
		typeClass() {
			return `toast-${this.type}`;
		},

		// Decide l'icona in base al tipo
		iconSrc() {
			switch (this.type) {
				case "error":
					return errorIcon;
				case "warning":
					return warningIcon;
				case "success":
				default:
					return successIcon;
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

.toast-success {
	border-left: 4px solid #2fa27a;
}
.toast-warning {
	border-left: 4px solid #fbc02d;
}
.toast-error {
	border-left: 4px solid #e53935;
}

.toast-content {
	display: flex;
	flex-direction: row;
	align-items: flex-start;
	gap: 10px;
}

.toast-icon {
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

.accent-bar {
	height: 4px;
	width: 100%;
	margin-top: 10px;
	border-bottom-left-radius: 8px;
	border-bottom-right-radius: 8px;
}
.accent-success {
	background-color: #2fa27a;
}
.accent-warning {
	background-color: #fbc02d;
}
.accent-error {
	background-color: #e53935;
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
