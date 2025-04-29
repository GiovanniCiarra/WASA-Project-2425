<template>
	<div class="info-overlay">
		<div class="top-bar">
			<img
				@click="$emit('closeOverlay')"
				class="close-btn"
				src="../assets/icons/close-big.png"
				alt="Close"
			/>
			<h1 class="title">Info Gruppo</h1>
		</div>

		<div class="info-content">
			<div class="group-header">
				<div class="photo-wrapper" @click="onChangePhoto">
					<img
						:src="newGroupPhoto || groupPhoto || defaultGroupPhoto"
						alt="Foto Gruppo"
						class="info-photo"
					/>
					<input
						type="file"
						ref="fileInput"
						@change="onFileSelected"
						accept="image/*"
						class="hidden-input"
					/>
				</div>
				<div class="name-edit">
					<input
						v-model="newGroupName"
						class="name-input"
						:placeholder="groupName"
					/>
					<button
						class="save-btn"
						@click="saveGroupName"
						:disabled="!canSave"
					>
						Salva
					</button>
				</div>
			</div>

			<div class="participants-header">Partecipanti:</div>

			<div class="member-list">
				<MemberItem
					v-for="member in members"
					:key="member.userId"
					:member="member"
					:currentUser="currentUser"
					@leaveGroup="leaveGroup"
				/>
			</div>

			<button class="add-member-btn" @click="showAddPopup = true">
				<img src="../assets/icons/più.png" alt="Add" />
				Aggiungi Partecipanti
			</button>

			<AddMembersPopup
				v-if="showAddPopup"
				:users="availableUsers"
				@close="showAddPopup = false"
				@confirm="handleAddMembers"
			/>
		</div>
	</div>
</template>

<script>
import defaultGroupPhoto from "../assets/images/default-user.png";
import MemberItem from "../components/MemberItem.vue";
import AddMembersPopup from "../components/AddMembersPopup.vue";

export default {
	components: {
		MemberItem,
		AddMembersPopup,
	},

	props: {
		groupName: String,
		groupPhoto: String,
		members: Array,
		currentUser: Number,
	},
	emits: [
		"closeOverlay",
		"leaveGroup",
		"addMembers",
		"updateGroupName",
		"updateGroupPhoto",
	],

	data() {
		return {
			defaultGroupPhoto,
			newGroupName: "",
			newGroupPhoto: null,
			showAddPopup: false,
			availableUsers: [],
		};
	},

	computed: {
		canSave() {
			return (
				(this.newGroupName.trim() &&
					this.newGroupName.trim() !== groupName) ||
				this.newGroupPhoto
			);
		},
	},

	created() {
		this.fetchAvailableUsers();
	},

	methods: {
		async fetchAvailableUsers() {
			// Faccio il fetch degli utenti
			try {
				const res = await this.$axios.get("/users", {
					headers: {
						Authorization: `Bearer ${this.currentUser}`,
					},
				});

				// Filtro solo gli utenti non già presenti nel gruppo
				const memberIds = this.members.map((m) => m.userId);
				this.availableUsers = res.data.filter(
					(u) => !memberIds.includes(u.userId)
				);
			} catch (err) {
				console.error("Errore nel caricamento degli utenti", err);
			}
		},

		handleAddMembers(userIds) {
			this.$emit("addMembers", userIds);
			this.showAddPopup = false;
		},

		leaveGroup() {
			this.$emit("leaveGroup");
		},

		saveGroupName() {
			if (this.newGroupName.trim() !== "") {
				this.emit("updateGroupName", this.newGroupName.trim());
			}
		},

		onChangePhoto() {
			this.$refs.fileInput.click();
		},
		onFileSelected(event) {
			const file = event.target.files[0];
			if (file) {
				this.newGroupPhoto = file;
				this.$emit("updateGroupPhoto", file);
			}
		},
	},
};
</script>

<style scoped>
.info-overlay {
	background: rgb(242, 242, 242);
	padding: 24px;
	overflow-y: auto;
	flex-grow: 1;
	display: flex;
	flex-direction: column;
	animation: slideDown 0.4s ease forwards;
}

@keyframes slideDown {
	from {
		transform: translateY(-20px);
		opacity: 0;
	}
	to {
		transform: translateY(0);
		opacity: 1;
	}
}

.top-bar {
	display: flex;
	align-items: center;
	margin-bottom: 20px;
}

.close-btn {
	font-size: 1.8rem;
	background: none;
	border: none;
	cursor: pointer;
	margin-right: 20px;
	opacity: 0.6;
}

.close-btn:hover {
	opacity: 1;
}

.title {
	font-size: 2rem;
	font-weight: bold;
}

.info-content {
	display: flex;
	flex-direction: column;
	align-items: center;
}

.group-header {
	display: flex;
	align-items: center;
	gap: 20px;
	margin-bottom: 20px;
}

.photo-wrapper {
	position: relative;
	cursor: pointer;
}

.info-photo {
	width: 100px;
	height: 100px;
	border-radius: 50%;
	object-fit: cover;
	border: 2px solid #2fa27a;
}

.hidden-input {
	display: none;
}

.name-edit {
	display: flex;
	align-items: center;
	gap: 10px;
}

.name-input {
	font-size: 1.4rem;
	padding: 8px 12px;
	border: 1px solid #ccc;
	border-radius: 8px;
	width: 300px;
}

.save-btn {
	padding: 8px 14px;
	background-color: #2fa27a;
	color: white;
	border: none;
	border-radius: 8px;
	cursor: pointer;
	font-size: 1rem;
}

.save-btn:hover {
	background-color: #27986d;
}

.save-btn:disabled {
	background-color: #a5d6c5;
	cursor: not-allowed;
}

.participants-header {
	align-self: flex-start;
	font-size: 1.4rem;
	margin: 70px 0 10px 0;
	font-weight: 600;
	color: #2fa27a;
}

.member-list {
	width: 100%;
	display: flex;
	flex-direction: column;
}

.add-member-btn {
	margin-top: 20px;
	padding: 1.3em 3em;
	font-size: 12px;
	font: sans-serif;
	text-transform: uppercase;
	letter-spacing: 2.5px;
	font-weight: 500;
	color: #000;
	background-color: #fff;
	border: none;
	border-radius: 45px;
	box-shadow: 0px 8px 15px rgba(0, 0, 0, 0.1);
	transition: all 0.3s ease 0s;
	cursor: pointer;
	outline: none;
}

.add-member-btn:hover {
	background-color: #23c483;
	box-shadow: 0px 15px 20px rgba(46, 229, 157, 0.4);
	color: #fff;
	transform: translateY(-7px);
}

.add-member-btn:active {
	transform: translateY(-1px);
}
</style>
