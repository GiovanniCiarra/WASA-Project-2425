<template>
	<div
		class="conversation-item"
		:class="{ selected: isSelected }"
		@click="$emit('conversationSelected', conversation.conversationId)"
	>
		<img
			:src="
				conversation.photo ||
				(conversation.type === 'group'
					? defaultGroupPhoto
					: defaultUserPhoto)
			"
			alt="Foto"
			class="conversation-photo"
		/>
		<div class="conversation-info">
			<div class="conversation-name">{{ conversation.name }}</div>
			<div class="conversation-snippet">
				<template v-if="conversation.lastMessage">
					<template v-if="conversation.lastMessage.type === 'photo'">
						<img
							src="@/assets/icons/image-snippet.png"
							alt="Foto"
							class="snippet-icon"
						/>
						Image
					</template>
					<template v-else>
						{{ conversation.lastMessage.content }}
					</template>
				</template>
				<template v-else> Nessun messaggio </template>
			</div>
		</div>
	</div>
</template>

<script>
import defaultGroupPhoto from "../assets/images/default-group-48.png";
import defaultUserPhoto from "../assets/images/default-user-48.png";

export default {
	props: {
		conversation: Object,
		isSelected: Boolean,
	},

	data() {
		return {
			defaultGroupPhoto,
			defaultUserPhoto,
		};
	},
};
</script>

<style scoped>
.conversation-item {
	display: flex;
	align-items: center;
	padding: 10px;
	border-bottom: 1px solid #73bfa1;
	cursor: pointer;
	border-radius: 3px;
	transition: background-color 0.3s, box-shadow 0.3s;
}

.conversation-item:hover {
	background-color: rgba(0, 128, 128, 0.1);
	box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.conversation-photo {
	width: 40px;
	height: 40px;
	border-radius: 50%;
	margin-right: 10px;
	object-fit: cover;
	box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.conversation-info {
	flex: 1;
}

.conversation-name {
	font-weight: bold;
}

.conversation-snippet {
	color: #666;
	font-size: 0.9rem;
	margin-top: 2px;
	overflow: hidden;
	text-overflow: ellipsis;
	white-space: nowrap;
	max-width: 180px;
}

.snippet-icon {
	width: 16px;
	height: 16px;
	margin-right: 2px;
	vertical-align: middle;
}

.selected {
	background-color: #d2f1e6;
	border-left: 4px solid #2fa27a;
}
</style>
