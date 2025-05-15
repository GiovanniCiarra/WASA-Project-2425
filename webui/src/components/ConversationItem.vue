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
						<span class="snippet-with-icon">
							<img
								:src="imageSnippet"
								alt="Foto"
								class="snippet-icon"
							/>
							<span>Image</span>
						</span>
					</template>
					<template v-else>
						<template v-if="conversation.type === 'group'">
							<strong
								>{{ conversation.lastMessage.senderUsername }}:
							</strong>
							<span>
								{{ conversation.lastMessage.content }}
							</span>
						</template>
						<template v-else>
							{{ conversation.lastMessage.content }}
						</template>
					</template>
				</template>
				<template v-else> Nessun messaggio </template>
			</div>
		</div>

		<!-- Notifica di messaggio non letto -->
		<img
			v-if="
				conversation.lastMessage &&
				conversation.lastMessage.status !== 'read' &&
				conversation.lastMessage.senderId !== currentUser &&
				conversation.type !== 'group' &&
				($route.name !== 'Conversation' ||
					conversation.conversationId !==
						Number($route.params.conversationId))
			"
			:src="newMessageDot"
			alt="Non letto"
			class="unread-dot"
		/>
	</div>
</template>

<script>
import defaultGroupPhoto from "../assets/images/default-group-48.png";
import defaultUserPhoto from "../assets/images/default-user-48.png";
import imageSnippet from "../assets/icons/image-snippet.png";
import newMessageDot from "../assets/icons/new-message.png";

export default {
	props: {
		conversation: Object,
		isSelected: Boolean,
		currentUser: Number,
	},

	data() {
		return {
			defaultGroupPhoto,
			defaultUserPhoto,
			imageSnippet,
			newMessageDot,
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
	max-width: 170px;
}

.snippet-icon {
	width: 16px;
	height: 16px;
	margin-right: 2px;
	vertical-align: middle;
}

.snippet-with-icon {
	display: inline-flex;
	align-items: center;
	gap: 4px;
}

.snippet-icon {
	width: 16px;
	height: 16px;
}

.selected {
	background-color: #d2f1e6;
	border-left: 4px solid #2fa27a;
}

.unread-dot {
	width: 12px;
	height: 12px;
	margin-left: 8px;
}
</style>
