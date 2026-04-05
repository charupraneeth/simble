import { computed, ref } from "vue";
import router from "../router";

export type User = {
    username: string
    email: string | null
    expiresAt: Date
}

const user = ref<User | null>(null)
const loading = ref(false)
const isLoaded = ref(false)

console.log(user)

function transformUser(raw: { username: string; email: string; expires_at: string }): User {
    return {
        username: raw.username,
        email: raw.email || null,
        expiresAt: new Date(raw.expires_at)
    }
}


export function useAuth() {

    const isLoggedIn = computed(() => !!user.value)

    const login = async () => {
        if (isLoaded.value) return; // Don't fetch again if we already know state
        try {
            loading.value = true
            const response = await fetch('/api/me')
            if (response.ok) {

                const userJSON = await response.json()

                user.value = transformUser(userJSON)
            } else {
                user.value = null
            }
        } catch (error) {
            console.error("Failed to fetch user state: ", error)
            user.value = null
        } finally {
            loading.value = false
            isLoaded.value = true
        }
    }

    const logout = async () => {
        try {
            loading.value = true
            await fetch("/api/session", { method: 'DELETE' })
            user.value = null
            isLoaded.value = false
            router.push('/')
        } catch (error) {
            console.error("Failed to log out:", error)
        } finally {
            loading.value = false
        }
    }

    return {
        user,
        isLoggedIn,
        login,
        loading,
        logout
    }
}