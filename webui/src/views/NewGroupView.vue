<template>
	<div class="new-group-container">
		<div class="left-column">
			<div class="photo-wrapper">
				<label class="group-photo">
					<input
						type="file"
						accept="image/*"
						@change="handlePhotoUpload"
						hidden
					/>
					<!-- Immagine del Gruppo -->
					<img
						class="group-photo-image"
						:src="groupPhoto || defaultGroupImage"
						alt="Group"
					/>
					<!-- Icona modifica -->
					<img
						class="edit-icon"
						src="../assets/icons/pen.png"
						alt="Modifica"
					/>
				</label>
				<!--  Icona rimuovi -->
				<img
					v-if="groupPhoto"
					class="remove-icon"
					src="../assets/icons/close.png"
					alt="Rimuovi"
					@click.stop="removeGroupPhoto"
				/>
			</div>
			<input
				class="group-name"
				v-model="groupName"
				placeholder="Nome del gruppo"
			/>
		</div>

		<div class="right-column">
			<input
				v-model="searchQuery"
				type="text"
				class="search-bar"
				placeholder="Cerca utenti..."
			/>

			<div class="user-list">
				<UserItem
					v-for="user in filteredUsers"
					:key="user.userId"
					:user="user"
					:selected="isSelected(user)"
					@click="toggleUser(user)"
				/>
			</div>

			<button
				class="create-btn"
				@click="createGroup"
				:disabled="!canCreate"
			>
				Crea Gruppo
			</button>
		</div>
	</div>
</template>

<script>
import UserItem from "../components/UserItem.vue";
import defaultGroupImage from "../assets/images/default-user.png";

export default {
	components: { UserItem },

	data() {
		return {
			users: [],
			selectedUsers: [],
			searchQuery: "",
			groupName: "",
			groupPhoto: null,
			defaultGroupImage,
			userId: Number(sessionStorage.getItem("userId")),
		};
	},

	computed: {
		filteredUsers() {
			return this.users
				.filter((u) => u.userId !== this.userId)
				.filter((u) =>
					u.username
						.toLowerCase()
						.includes(this.searchQuery.toLowerCase())
				)
				.sort((a, b) => {
					const aSel = this.isSelected(a);
					const bSel = this.isSelected(b);
					return aSel === bSel ? 0 : aSel ? -1 : 1;
				});
		},

		canCreate() {
			return this.groupName.trim() && this.selectedUsers.length >= 2;
		},
	},

	async created() {
		try {
			const res = await this.$axios.get("/users", {
				headers: { Authorization: `Bearer ${this.userId}` },
			});
			this.users = res.data;
		} catch (err) {
			console.error("Errore nel caricamento utenti", err);
		}
	},

	methods: {
		isSelected(user) {
			return this.selectedUsers.some((u) => u.userId === user.userId);
		},

		toggleUser(user) {
			if (this.isSelected(user)) {
				this.selectedUsers = this.selectedUsers.filter(
					(u) => u.userId !== user.userId
				);
			} else {
				this.selectedUsers.push(user);
			}
		},

		handlePhotoUpload(event) {
			const file = event.target.files[0];
			if (file) {
				const reader = new FileReader();
				reader.onload = () => {
					this.groupPhoto = reader.result;
				};
				reader.readAsDataURL(file);
			}
		},

		removeGroupPhoto() {
			this.groupPhoto = null;
		},

		async createGroup() {
			try {
				const res = await this.$axios.post(
					"/conversations",
					{
						name: this.groupName,
						photo: this.groupPhoto || "",
						members: this.selectedUsers.map((u) => u.userId),
					},
					{
						headers: { Authorization: `Bearer ${this.userId}` },
					}
				);

				const { conversationId } = res.data;
				this.$router.push({
					path: `/conversations/${conversationId}`,
					query: {
						type: "group",
						name: this.groupName,
						photo: this.groupPhoto,
					},
				});
			} catch (err) {
				console.error("Errore nella creazione del gruppo", err);
			}
		},
	},
};
</script>

<style scoped>
.new-group-container {
	display: flex;
	gap: 30px;
	padding: 20px;
	height: 100%;
	box-sizing: border-box;
	background-color: #f0f0f0;
}

.left-column {
	width: 280px;
	display: flex;
	flex-direction: column;
	align-items: center;
	gap: 20px;
	margin-top: 40px;
}

.right-column {
	flex-grow: 1;
	display: flex;
	flex-direction: column;
	gap: 15px;
}

.group-photo {
	position: relative;
	cursor: pointer;
}

.group-photo-image {
	width: 120px;
	height: 120px;
	border-radius: 50%;
	object-fit: cover;
	border: 2px solid #ccc;
	box-shadow: 0 0 8px rgba(0, 0, 0, 0.2);
}

.edit-icon {
	position: absolute;
	bottom: 0;
	right: 0;
	width: 24px;
	height: 24px;
	object-fit: contain;
	border-radius: 50%;
	padding: 3px;
	box-shadow: 0 0 3px rgba(0, 0, 0, 0.2);
	transform: translate(25%, 25%);
}

.photo-wrapper {
	position: relative;
	width: fit-content;
}

.remove-icon {
	position: absolute;
	top: -12px;
	right: -12px;
	width: 22px;
	height: 22px;
	object-fit: contain;
	background: white;
	border-radius: 50%;
	padding: 3px;
	box-shadow: 0 0 3px rgba(0, 0, 0, 0.2);
	cursor: pointer;
	z-index: 3;
	opacity: 0.6;
}

.remove-icone:hover {
	opacity: 1;
}

.group-name {
	margin-top: 20px;
	width: 100%;
	padding: 10px;
	border-radius: 10px;
	border: 1px solid #ccc;
	text-align: center;
	font-size: 1rem;
}

.search-bar {
	padding: 10px;
	width: 100%;
	border-radius: 10px;
	border: 1px solid #ccc;
	font-size: 1rem;
	background-color: #f3f3f3;
	transition: all 0.3s ease;
}

.search-bar:focus {
	border: 1px solid #2fa27a;
	background-color: #fff;
	outline: none;
}

.user-list {
	flex-grow: 1;
	overflow-y: auto;
	display: flex;
	flex-direction: column;
	gap: 10px;
	max-height: 450px;
}

.create-btn {
	margin-top: auto;
	padding: 12px 24px;
	background-color: #2fa27a;
	color: white;
	border: none;
	border-radius: 10px;
	font-weight: bold;
	font-size: 1rem;
	cursor: pointer;
	transition: all 0.2s ease;
	box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
}

.create-btn:hover {
	background-color: #27a06c;
	transform: translateY(-2px);
}

.create-btn:disabled {
	background-color: #ccc;
	cursor: not-allowed;
	box-shadow: none;
}
</style>
