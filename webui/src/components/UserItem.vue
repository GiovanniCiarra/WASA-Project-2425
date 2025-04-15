<template>
	<div class="user-item">
		<img
			:src="user.photo || '/images/default-user.png'"
			alt="Foto profilo"
			class="user-photo"
		/>
		<span class="user-name">{{ user.username }}</span>
	</div>
</template>

<script>
export default {
	props: {
		user: {
			type: Object,
			required: true,
		},
	},
};
</script>

<style scoped>
.user-item {
	display: flex;
	align-items: center;
	padding: 12px 16px;
	border-bottom: 1px solid #e0e0e0;
	cursor: pointer;
	transition: transform 0, 2s ease, background-color 0.2s ease;
	border-radius: 8px;
}

.user-item:hover {
	background-color: #f5f5f5;
	transform: scale(1.02);
}

.user-photo {
	width: 48px;
	height: 48px;
	border-radius: 50%;
	object-fit: cover;
	margin-right: 14px;
}

.user-name {
	font-size: 1rem;
	font-weight: 500;
	color: #222;
}
</style>
