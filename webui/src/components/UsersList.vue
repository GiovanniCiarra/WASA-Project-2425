<template>
	<div class="user-list">
		<input
			v-model="searchUser"
			type="text"
			placeholder="Search User..."
			class="search-bar"
		/>
		<div v-if="userError" class="error-message">{{ userError }}</div>
		<div v-for="user in filteredUsers" :key="user.userId">
			<UserItem :user="user" @select="$emit('userSelected', $event)" />
		</div>
	</div>
</template>

<script>
import UserItem from "./UserItem.vue";

export default {
	components: { UserItem },
	props: {
		users: Array,
		userError: String,
	},
	data() {
		return {
			searchUser: "",
		};
	},
	computed: {
		filteredUsers() {
			return this.users.filter((user) =>
				user.username
					.toLowerCase()
					.includes(this.searchUser.toLowerCase())
			);
		},
	},
};
</script>

<style scoped>
.search-bar {
	width: 100%;
	padding: 8px;
	padding-left: 0.8em;
	outline: none;
	overflow: hidden;
	margin-bottom: 10px;
	border: 2px solid transparent;
	border-radius: 10px;
	background-color: #f3f3f3;
	transition: all 0.5s;
}

.search-bar:hover,
.search-bar:focus {
	border: 2px solid #4a9dec;
	box-shadow: 0px 0px 0px 7px rgb(74, 157, 236, 20%);
	background-color: white;
}

.user-list {
	max-height: 80vh;
	overflow-y: auto;
}

.error-message {
	color: red;
	font-weight: bold;
	text-align: center;
	margin-top: 10px;
	background: #ffe0e0;
	padding: 10px;
	border-radius: 5px;
}
</style>
