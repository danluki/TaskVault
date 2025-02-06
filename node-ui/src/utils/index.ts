import Cookies from "js-cookie";

export interface User {
    Username: string,
    Password: string,
}

const USER_COOKIE_KEY = "user_session";

export function getUser(): User | null {
    let userJson = Cookies.get(USER_COOKIE_KEY); // Get user data from cookie
    let user: User | null = null;

    if (userJson) {
        try {
            user = JSON.parse(userJson);
        } catch (e) {
            console.error("Invalid cookie format:", e);
        }
    }

    return user;
}

export function setUser(user: User) {
    Cookies.set(USER_COOKIE_KEY, JSON.stringify(user), { expires: 7 }); // Cookie expires in 7 days
}

export function logoutUser() {
    Cookies.remove(USER_COOKIE_KEY); // Remove cookie
}
