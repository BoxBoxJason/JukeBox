import Cookies from "js-cookie"
import { LOCAL_STORAGE_KEYS } from "@/constants/storage";
import { COOKIE_KEYS } from "@/constants/storage"

/**
 * Checks if the user is already logged in,
 * If they are not or some information is missing, they are automatically either
 * Logged in with refresh token or access token. Or they are completely logged out.
 *
 * @returns {Promise<boolean>} True if the user is logged in, false otherwise.
 */
export async function autoLogin(): Promise<boolean> {
  const localStorageUserId: string | null = localStorage.getItem(LOCAL_STORAGE_KEYS.USER_ID);
  const localStorageUsername: string | null = localStorage.getItem(LOCAL_STORAGE_KEYS.USERNAME);
  const accessToken: string | undefined = Cookies.get(COOKIE_KEYS.ACCESS_TOKEN);
  const refreshToken: string | undefined = Cookies.get(COOKIE_KEYS.REFRESH_TOKEN);

  let authSuccess = localStorageUserId != null && localStorageUsername != null && accessToken != undefined && refreshToken != undefined;

  if (!authSuccess) {
    if (accessToken != null) {
      authSuccess = await loginFromAccessToken();
    } else if (refreshToken != null) {
      authSuccess = await loginFromRefreshToken();
    }
  }
  return authSuccess;
}

/**
 * Logs in the user using the provided credentials.
 * Stores the user's identity in local storage if the login was successful.
 *
 * @return {Promise<boolean>} True if the login was successful, false otherwise.
 */
export async function loginFromRefreshToken(): Promise<boolean> {
  let loggedIn = false;
  const refreshToken: string | undefined = Cookies.get(COOKIE_KEYS.REFRESH_TOKEN);
  if (refreshToken != undefined && refreshToken.trim()) {
    const response = await fetch("/api/auth/refresh", {
      method: "POST",
      credentials: "include",
    });

    if (response.ok) {
      const data = await response.json();
      setIdentity(data.user_id, data.username);
      loggedIn = true;
    }
  }
  return loggedIn;
}

/**
 * Logs in the user using the provided credentials.
 * Stores the user's identity in local storage if the login was successful.
 *
 * @return {Promise<boolean>} True if the login was successful, false otherwise.
 */
export async function loginFromAccessToken(): Promise<boolean> {
  let loggedIn = false;
  const accessToken: string | undefined = Cookies.get(COOKIE_KEYS.ACCESS_TOKEN);
  if (accessToken != undefined) {
    const response = await fetch("/api/auth/login", {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${accessToken}`,
      },
    });

    if (response.ok) {
      const data = await response.json();
      setIdentity(data.user_id, data.username);
      loggedIn = true;
    }
  }
  return loggedIn;
}

/**
 * Sets the user's identity in local storage.
 *
 * @param {number} user_id The user's ID.
 * @param {string} username The user's username.
 *
 * @returns {void}
 */
export function setIdentity(user_id: number, username: string) {
  localStorage.setItem(LOCAL_STORAGE_KEYS.USERNAME, username);
  localStorage.setItem(LOCAL_STORAGE_KEYS.USER_ID, user_id.toString());
}

/**
 * Sends a request to the server to invalidate the access and refresh token.
 * Logs out the user by removing their identity from local storage.
 * Also deletes the access and refresh tokens.
 *
 * @returns {Promise<boolean>} True if the logout was successful, false otherwise.
 */
export async function fullLogout(): Promise<boolean> {
  const response = await fetch("/api/auth/logout", {
    method: "POST",
    credentials: "include",
  });

  if (response.ok) {
    clearIdentity();
  }

  return response.ok;
}

/**
 * Removes the user's identity from local storage.
 * Also removes the access and refresh tokens from cookies
 */
export function clearIdentity() {
  localStorage.removeItem(LOCAL_STORAGE_KEYS.USERNAME);
  localStorage.removeItem(LOCAL_STORAGE_KEYS.USER_ID);
  Cookies.remove(COOKIE_KEYS.REFRESH_TOKEN);
  Cookies.remove(COOKIE_KEYS.ACCESS_TOKEN);
}

/**
 * Checks if the user is fully signed in.
 *
 * @returns {boolean} True if the user is fully signed in, false otherwise.
 */
export function isUserFullySignedIn() {
  const localStorageUserId: string | null = localStorage.getItem(LOCAL_STORAGE_KEYS.USER_ID);
  const localStorageUsername: string | null = localStorage.getItem(LOCAL_STORAGE_KEYS.USERNAME);
  const accessToken: string | undefined = Cookies.get(COOKIE_KEYS.ACCESS_TOKEN);
  const refreshToken: string | undefined = Cookies.get(COOKIE_KEYS.REFRESH_TOKEN);

  return localStorageUserId != null && localStorageUsername != null && accessToken != undefined && refreshToken != undefined;
}