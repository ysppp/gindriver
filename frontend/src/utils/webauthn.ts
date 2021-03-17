import * as APIs from '../utils/apis';
import { loadingMessage, errorMessage, successMessage } from '../components/message';
import { create, get, supported } from '@github/webauthn-json';
import { history } from 'umi';

// compatibility check
const browserSupport = () => {
  if (!supported()) {
    errorMessage("This browser does not support WebAuthn!");
    setTimeout(() => {
      errorMessage("Please use compatible browers or operating system.");
    }, 3000);
  }
}

// Login
const webauthnLogin = async (username: string) => {
  loadingMessage("loginLoadingMsg");
  try {
    if (!username.match(/^[a-zA-Z0-9]{4,16}$/)) {
      errorMessage("Bad username: ^[a-zA-Z0-9]{4,16}$");
      return;
    }

    const beginResponse = await fetch(APIs.API_LOGIN_BEGIN.replace("{}", username));
    if (beginResponse.status === 400) {
      errorMessage("Bad username!");
      return;
    }
    if (beginResponse.status === 500) {
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
    const finishResponseBody = await finishResponse.json();
    if (finishResponse.status === 400) {
      errorMessage("Bad credential!");
      return;
    }
    if (finishResponse.status === 500) {
      errorMessage("Server error!");
      return;
    }
    successMessage("loginLoadingMsg", "Login success!");

    localStorage.setItem("jwt", finishResponseBody.token);
    localStorage.setItem("username", username);
    setTimeout(() => {
      history.push("/uploadDocument");
    }, 2500);
  } catch (e) {
    if (e.message.indexOf("The operation either timed out or was not allowed.") > -1) {
      errorMessage("Operation failed or Canceled by user");
      return;
    }
    errorMessage(e.message);
  }
}

const webauthnReg = async (username: string) => {
  loadingMessage("regLoadingMsg");
  try {
    if (!username.match(/^[a-zA-Z0-9]{4,16}$/)) {
      errorMessage("Bad username: ^[a-zA-Z0-9]{4,16}$");
      return;
    }

    const beginResponse = await fetch(APIs.API_REGISTER_BEGIN, {
      method: "POST",
      body: JSON.stringify({username: username}),
      headers: new Headers({
        'Content-Type': 'application/json'
      })
    });

    const beginResponseBody = await beginResponse.json();
    if (beginResponse.status === 400) {
      if (beginResponseBody.message === "user exist") {
        errorMessage("Username aleady existed!");
        return;
      }
      errorMessage("Bad username!");
      return;
    }
    if (beginResponse.status === 500) {
      errorMessage("Server error!");
      return;
    }

    const finishRequestBody = await create(beginResponseBody);
    const finishResponse = await fetch(APIs.API_REGISTER_FINISH.replace("{}", username), {
      method: "PATCH",
      body: JSON.stringify(finishRequestBody),
      headers: new Headers({
        'Content-Type': 'application/json'
      })
    });
    if (finishResponse.status === 400) {
      errorMessage("Bad credential!");
      return;
    }
    if (finishResponse.status === 500) {
      errorMessage("Server error!");
      return;
    }
    successMessage("regLoadingMsg", "Registration success!");
  } catch (e) {
    if (e.message.indexOf("The operation either timed out or was not allowed.") > -1) {
      errorMessage("Operation failed or Canceled by user");
      return;
    }
    errorMessage(e.message);
  }
}

export {
  webauthnLogin,
  webauthnReg,
  browserSupport
}
