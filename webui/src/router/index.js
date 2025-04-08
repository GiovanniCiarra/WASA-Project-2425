import { createRouter, createWebHashHistory } from "vue-router";
import HomeView from "../views/HomeView.vue";
import LoginView from "../views/LoginView.vue";
import ConversationView from "../views/ConversationView.vue";
import ProfileView from "../views/ProfileView.vue";

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{ path: "/Login", component: LoginView },
		{ path: "/Home", component: HomeView },
		{
			path: "/Conversations/id/:conversationId",
			component: ConversationView,
		},
		{ path: "/Conversations/:userId", component: ConversationView },
		{ path: "/Profile/:userId", component: ProfileView },
		{ path: "/:pathMatch(.*)*", redirect: "/Login" },
	],
});

export default router;
