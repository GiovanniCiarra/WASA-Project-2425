<template>
	<div class="user-item" @click="$emit('select', user.userId)">
		<img
			:src="user.photo || '/images/profile.png'"
			alt="Avatar"
			class="user-avatar"
		/>
		<span class="user-name">{{ user.username }}</span>
	</div>
</template>

<script>
export default {
	props: {
		user: Object,
	},
};
</script>

<style scoped>
.user-item {
	display: flex;
	align-items: center;
	padding: 10px;
	border-bottom: 1px solid #ddd;
	cursor: pointer;
}
.user-avatar {
	width: 40px;
	height: 40px;
	border-radius: 50%;
	margin-right: 10px;
}
</style>
