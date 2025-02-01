import { LOCAL_STORAGE_KEYS } from "@/constants/storage";

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

  let authSuccess = localStorageUserId != null && localStorageUsername != null;

  if (!authSuccess) {
    authSuccess = await loginFromAccessToken();
    if (!authSuccess) {
      authSuccess = await loginFromRefreshToken();
    }
  }

  if (!authSuccess) {
    clearIdentity();
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

  const response = await fetch("/api/auth/refresh", {
    method: "POST",
    credentials: "include",
  });

  if (response.ok) {
    const data = await response.json();
    setIdentity(data.user_id, data.username);
    loggedIn = true;
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

  const response = await fetch("/api/auth/login", {
    method: "POST",
    credentials: "include",
  });

  if (response.ok) {
    const data = await response.json();
    setIdentity(data.user_id, data.username);
    loggedIn = true;
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
  const event = new CustomEvent("localStorageChange");
  window.dispatchEvent(event);
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
}

export function getUsername(): string | null {
  return localStorage.getItem(LOCAL_STORAGE_KEYS.USERNAME);
}

export function getUserId(): number | null {
  const userId = localStorage.getItem(LOCAL_STORAGE_KEYS.USER_ID);
  return userId == null ? null : parseInt(userId);
}

export function isUserConnected(): boolean {
  return getUserId() != null && getUsername() != null;
}
