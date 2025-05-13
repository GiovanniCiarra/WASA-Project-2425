<template>
	<transition name="slide" @after-leave="$emit('closeOverlay')">
		<div class="info-overlay" v-if="visible">
			<div class="top-bar">
				<img
					@click="visible = false"
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
							:src="
								newGroupPhoto || groupPhoto || defaultGroupPhoto
							"
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
							@keyup.enter="saveGroupName"
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

				<div class="member-list-scrollable">
					<MemberItem
						v-for="member in sortedMembers"
						:key="member.userId"
						:member="member"
						:currentUser="currentUser"
						@leaveGroup="leaveGroup"
					/>
				</div>

				<div class="add-member-container">
					<button class="add-member-btn" @click="showAddPopup = true">
						<img src="../assets/icons/più.png" alt="Add" />
						Aggiungi Partecipanti
					</button>
				</div>

				<AddMembersPopup
					v-if="showAddPopup"
					:users="availableUsers"
					@close="showAddPopup = false"
					@confirm="handleAddMembers"
				/>
			</div>
		</div>
	</transition>
</template>

<script>
import defaultGroupPhoto from "../assets/images/default-group-100.png";
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
			visible: true,
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
				this.newGroupName.trim() &&
				this.newGroupName.trim() !== this.groupName
			);
		},

		// Ordino i membri per avere quello corrente sempre in cima
		sortedMembers() {
			if (!this.members) return [];

			// Divido il membro corrente e gli altri
			const current = this.members.find(
				(member) => member.userId === this.currentUser
			);
			const others = this.members.filter(
				(member) => member.userId !== this.currentUser
			);

			// Ritorno lista con l'utente corrente per primo
			return current ? [current, ...others] : others;
		},
	},

	watch: {
		members: {
			immediate: true,
			handler(newVal, oldVal) {
				if (JSON.stringify(newVal) !== JSON.stringify(oldVal)) {
					this.fetchAvailableUsers();
				}
			},
		},
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
				this.handleError(err);
			}
		},

		handleAddMembers(newUsers) {
			this.$emit("addMembers", newUsers);
			this.showAddPopup = false;
		},

		leaveGroup() {
			this.$emit("leaveGroup");
		},

		saveGroupName() {
			if (this.newGroupName.trim() !== "") {
				this.$emit("updateGroupName", this.newGroupName.trim());
			}
			this.newGroupName = "";
		},

		onChangePhoto() {
			this.$refs.fileInput.click();
		},

		onFileSelected(event) {
			const file = event.target.files[0];
			if (file) {
				const reader = new FileReader();
				reader.onload = () => {
					this.newGroupPhoto = reader.result;
					this.$emit("updateGroupPhoto", this.newGroupPhoto);
				};
				reader.readAsDataURL(file);
			}
		},

		handleError(err) {
			if (!err.response) {
				this.$root.showToast(
					"Errore di rete: controlla la connessione.",
					"error"
				);
				return;
			}

			const status = err.response.status;
			const detail = err.response.data || "ERRORE";

			switch (status) {
				case 401:
					this.$root.showToast(
						"Errore 401: Non autorizzato, effettua di nuovo il login",
						"error"
					);
					this.$router.push("/Login");
					break;
				case 500:
					this.$root.showToast(`Errore 500:\n ${detail}!`, "error");
					break;
				default:
					this.$root.showToast("Errore sconosciuto", "error");
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
	height: 100%;
	box-sizing: border-box;
}

.slide-enter-active,
.slide-leave-active {
	transition: all 0.4s ease;
}

.slide-enter-from,
.slide-leave-to {
	transform: translateY(-20px);
	opacity: 0;
}

.slide-enter-to,
.slide-leave-from {
	transform: translateY(0);
	opacity: 1;
}

.top-bar {
	display: flex;
	justify-content: flex-end;
	align-items: center;
	margin-bottom: 20px;
	position: relative;
}

.close-btn {
	width: 32px;
	height: 32px;
	cursor: pointer;
	opacity: 0.7;
}

.close-btn:hover {
	opacity: 1;
}

.title {
	position: absolute;
	left: 50%;
	transform: translateX(-50%);
	font-size: 2.2rem;
	font-weight: bold;
	color: rgb(69, 69, 69);
}

.info-content {
	display: flex;
	flex-direction: column;
	align-items: center;
	flex-grow: 1;
	width: 100%;
}

.group-header {
	display: flex;
	align-items: center;
	background: rgba(255, 255, 255, 0.2);
	backdrop-filter: blur(10px);
	padding: 20px 30px;
	border-radius: 20px;
	width: 100%;
	max-width: 700px;
	margin-bottom: 20px;
	box-shadow: 0 8px 32px rgba(0, 0, 0, 0.15);
	gap: 20px;
}

.photo-wrapper {
	position: relative;
	cursor: pointer;
}

.info-photo {
	width: 120px;
	height: 120px;
	border-radius: 50%;
	object-fit: cover;
	box-shadow: 0 4px 12px solid rgba(0, 0, 0, 0.2);
	transition: transofr 0.3s ease, box-shadow 0.3s ease;
}

.photo-wrapper:hover .info-photo {
	transform: scale(1.05);
	box-shadow: 0 6px 20px rgba(0, 0, 0, 0.3);
}

.hidden-input {
	display: none;
}

.name-edit {
	display: flex;
	flex-direction: row;
	align-items: center;
	width: 100%;
	gap: 10px;
}

.name-input {
	flex-grow: 1;
	width: auto;
	max-width: none;
	padding: 16px 20px;
	font-size: 1.1rem;
	border: 1px solid #ddd;
	border-radius: 12px;
	background-color: rgba(255, 255, 255, 0.9);
	box-shadow: 0 1px 4px rgba(0, 0, 0, 0.1);
	text-align: center;
	transition: all 0.3s ease;
	font-weight: 500;
}

.name-input:focus {
	border-color: #73bfa1;
	box-shadow: 0 0 8px rgba(74, 157, 236, 0.4);
	outline: none;
}

.name-input::placeholder {
	color: #aaa;
	text-align: center;
	font-weight: 500;
}

.save-btn {
	height: 52px;
	padding: 0 30px;
	font-size: 1.1rem;
	border: none;
	border-radius: 50px;
	background: linear-gradient(135deg, #2fa27a, #73bfa1);
	color: white;
	cursor: pointer;
	transition: all 0.3s ease;
	box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
}

.save-btn:hover {
	background: linear-gradient(135deg, #24996d, #5bb896);
	box-shadow: 0 6px 18px rgba(0, 0, 0, 0.3);
	transform: translateY(-2px);
}

.save-btn:disabled {
	background: #a5d6c5;
	cursor: not-allowed;
	box-shadow: none;
	transform: none;
}

.participants-header {
	align-self: flex-start;
	font-size: 1.4rem;
	margin: 20px 0 10px 0;
	font-weight: 600;
	color: #2fa27a;
	width: 100%;
}

.member-list-scrollable {
	width: 100%;
	flex-grow: 1;
	max-height: calc(100vh - 520px);
	overflow-y: auto;
	display: flex;
	flex-direction: column;
	margin-bottom: 10px;
	scroll-behavior: smooth;
}

.add-member-container {
	width: 100%;
	display: flex;
	justify-content: center;
	padding-top: 10px;
	margin-top: 10px;
}

.add-member-btn {
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
