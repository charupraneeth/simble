import { computed, ref } from "vue";

export type User = {
    username: string
    email: string | null
    expiresAt: Date
}

const user = ref<User | null>(null)
const loading = ref(false)
const isLoaded = ref(false)

export function useAuth() {

    const isLoggedIn = computed(() => !!user.value)

    const login = async () => {
        if (isLoaded.value) return; // Don't fetch again if we already know state
        try {
            loading.value = true
            const response = await fetch('/api/me')
            if (response.ok) {

                const userJSON = await response.json()

                user.value = userJSON
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

    return {
        user,
        isLoggedIn,
        login,
        loading
    }
}