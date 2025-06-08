<template>
	<div class="comment-picker">
		<!-- Icona di chiusura del Picker -->
		<img
			src="../assets/icons/close.png"
			alt="close"
			class="close-icon"
			@click="$emit('close-emoji')"
		/>

		<!-- Griglia emoji -->
		<div class="emoji-grid">
			<span
				v-for="emoji in emojis"
				:key="emoji"
				@click="$emit('select-emoji', emoji)"
			>
				{{ emoji }}
			</span>
		</div>
	</div>
</template>

<script>
export default {
	data() {
		return {
			emojis: ["😄", "❤️", "👍", "👎", "😂", "😢", "🔥"],
		};
	},
};
</script>

<style scoped>
.comment-picker {
	position: absolute;
	background: white;
	border: 1px solid #ccc;
	border-radius: 8px;
	box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
	padding: 10px;
	display: flex;
	flex-direction: column;
	z-index: 9999;
	min-width: 180px;
}

.close-icon {
	width: 14px;
	height: 14px;
	align-self: flex-end;
	cursor: pointer;
	margin-bottom: 6px;
	opacity: 0.6;
}

.close-icon:hover {
	opacity: 1;
}

.emoji-grid {
	display: flex;
	flex-wrap: wrap;
	gap: 6px;
	justify-content: center;
}

.emoji-grid span {
	font-size: 1.1rem;
	cursor: pointer;
	transition: transform 0.2s;
}

.emoji-grid span:hover {
	transform: scale(1.3);
}
</style>
