import * as APIs from '../utils/apis';
import { loadingMessage, errorMessage, successMessage } from '../components/message';
import { create, get } from '@github/webauthn-json';

// Base64 to ArrayBuffer
function bufferDecode(value: string) {
  return Uint8Array.from(atob(value), c => c.charCodeAt(0));
}

// ArrayBuffer to URLBase64
function bufferEncode(value: ArrayBuffer) {
  return btoa(String.fromCharCode(...new Uint8Array(value)))
    .replace(/\+/g, "-")
    .replace(/\//g, "_")
    .replace(/=/g, "");;
}

// Login
const webauthnLogin = async (username: string) => {
  loadingMessage("loginLoadingMsg");
  try {
    const beginResponse = await fetch(APIs.API_LOGIN_BEGIN.replace("{}", username));
    if (beginResponse.status == 400) {
      errorMessage("Bad username!");
      return;
    }
    if (beginResponse.status == 500) {
      errorMessage("Server error!");
      return;
    }

    const beginResponseBody = await beginResponse.json();
    const finishRequestBody = await get(beginResponseBody);
    const finishResponse = await fetch(APIs.API_LOGIN_FINISH.replace("{}", username), {
      method: "PATCH",
      body: JSON.stringify(finishRequestBody),
      headers: new Headers({
        'Content-Type': 'application/json'
      })
    });
    if (finishResponse.status == 400) {
      errorMessage("Bad credential!");
      return;
    }
    if (finishResponse.status == 500) {
      errorMessage("Server error!");
      return;
    }
    successMessage("loginLoadingMsg", "Login success!");
  } catch (e) {
    errorMessage(e.message);
  }
}

const webauthnReg = async (username: string) => {
  loadingMessage("regLoadingMsg");
  try {
    const beginResponse = await fetch(APIs.API_REGISTER_BEGIN, {
      method: "POST",
      body: JSON.stringify({username: username}),
      headers: new Headers({
        'Content-Type': 'application/json'
      })
    });
    if (beginResponse.status == 400) {
      errorMessage("Bad username!");
      return;
    }
    if (beginResponse.status == 500) {
      errorMessage("Server error!");
      return;
    }

    const beginResponseBody = await beginResponse.json();
    const finishRequestBody = await create(beginResponseBody);
    const finishResponse = await fetch(APIs.API_REGISTER_FINISH.replace("{}", username), {
      method: "PATCH",
      body: JSON.stringify(finishRequestBody),
      headers: new Headers({
        'Content-Type': 'application/json'
      })
    });
    if (finishResponse.status == 400) {
      errorMessage("Bad credential!");
      return;
    }
    if (finishResponse.status == 500) {
      errorMessage("Server error!");
      return;
    }
    successMessage("regLoadingMsg", "Registration success!");
  } catch (e) {
    errorMessage(e.message);
  }
}

export {
  webauthnLogin,
  webauthnReg
}
