import { createRouter, createWebHashHistory } from "vue-router";
import LoginView from "../views/LoginView.vue";
import HomeView from "../views/HomeView.vue";
import NewChatView from "../views/NewChatView.vue";
import NewGroupView from "../views/NewGroupView.vue";
import ProfileView from "../views/ProfileView.vue";
import ConversationView from "../views/ConversationView.vue";

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{ path: "/Login", name: "Login", component: LoginView },
		{ path: "/Home", name: "Home", component: HomeView },
		{ path: "/NewChat", name: "newChat", component: NewChatView },
		{ path: "/NewGroup", name: "newGroup", component: NewGroupView },
		{ path: "/Profile", name: "Profile", component: ProfileView },
		{
			path: "/Conversations/:conversationId",
			name: "Conversation",
			component: ConversationView,
		},
		{ path: "/:pathMatch(.*)*", redirect: "/Login" },
	],
});

export default router;
