import {USER_LS_KEY} from "@/consts";

export interface User {
    Username: string,
    Password: string,
}

export function getUser(): User | null {
    let userJson = localStorage.getItem(USER_LS_KEY)
    let user: User | null = null

    if (userJson) {
        user = JSON.parse(userJson)
    }


    return user;
}

export function logoutUser() {
    localStorage.removeItem(USER_LS_KEY)
}