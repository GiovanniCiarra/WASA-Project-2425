<template>
	<div class="overlay">
		<div class="dialog">
			<p>{{ message }}</p>
			<div class="actions">
				<button class="cancel" @click="$emit('cancel')">Annulla</button>
				<button class="confirm" @click="$emit('confirm')">
					Elimina
				</button>
			</div>
		</div>
	</div>
</template>

<script>
export default {
	props: {
		message: {
			type: String,
			default: "Sei sicuro di voler eliminare?",
		},
	},
};
</script>

<style scoped>
.overlay {
	position: fixed;
	inset: 0;
	background-color: rgba(0, 0, 0, 0.4);
	display: flex;
	align-items: center;
	justify-content: center;
	z-index: 999;
}

.dialog {
	background: white;
	padding: 24px;
	border-radius: 12px;
	width: 300px;
	box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
	text-align: center;
	font-family: sans-serif;
}

.actions {
	margin-top: 20px;
	display: flex;
	justify-content: space-around;
}

button.cancel {
	background: #ccc;
	color: #333;
}

button.confirm {
	background: #d9534f;
	color: white;
}

button {
	padding: 8px 14px;
	border: none;
	border-radius: 6px;
	cursor: pointer;
	transition: 0.2s;
}

button:hover {
	opacity: 0.9;
}
</style>
