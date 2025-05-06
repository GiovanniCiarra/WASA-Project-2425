<template>
	<div class="header-container" @click="$emit('openInfo')">
		<img
			:src="
				photo ||
				(type === 'group' ? defaultGroupPhoto : defaultUserPhoto)
			"
			alt="Foto"
			class="header-photo"
		/>
		<span class="header-name">{{ name }}</span>
	</div>
</template>

<script>
import defaultGroupPhoto from "../assets/images/default-group-48.png";
import defaultUserPhoto from "../assets/images/default-user-48.png";

export default {
	props: {
		name: String,
		photo: String,
		type: String,
	},

	emits: ["openInfo"],

	data() {
		return {
			defaultGroupPhoto,
			defaultUserPhoto,
		};
	},
};
</script>

<style scoped>
.header-container {
	display: flex;
	align-items: center;
	padding: 12px 20px;
	background-color: #f8f8f8;
	border-bottom: 1px solid #ddd;
	cursor: pointer;
}

.header-photo {
	width: 40px;
	height: 40px;
	border-radius: 50%;
	object-fit: cover;
	margin-right: 12px;
}

.header-name {
	font-size: 1.1rem;
	font-weight: 600;
	color: #333;
}
</style>
