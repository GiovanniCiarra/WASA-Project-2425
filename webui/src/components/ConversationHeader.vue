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

		<div class="header-text">
			<span class="header-name">{{ name }}</span>
			<span v-if="type === 'group'" class="header-members">
				{{ truncatedMemberList }}
			</span>
		</div>
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
		members: Array,
	},

	emits: ["openInfo"],

	data() {
		return {
			defaultGroupPhoto,
			defaultUserPhoto,
		};
	},

	computed: {
		truncatedMemberList() {
			const usernames = this.members.map((m) => m.username);
			const result = usernames.join(", ");
			return result.length > 100 ? result.slice(0, 100) + "..." : result;
		},
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

.header-text {
	display: flex;
	flex-direction: column;
	overflow: hidden;
}

.header-name {
	font-size: 1.1rem;
	font-weight: 600;
	color: #333;
}

.header-members {
	font-size: 0.85rem;
	color: #666;
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
	max-width: 900px; /* regola a seconda dello spazio disponibile */
}
</style>
