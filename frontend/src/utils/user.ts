import { API_USER_INFO } from './apis';
import { errorMessage } from '../components/message';
import { history } from 'umi';

const getUserInfo = async () => {

  // frontend debug
  // return "admin";

  const username = localStorage.getItem("username");
  if (username === null) {
    return invalidSessionJumpBack();
  }

  const userCredential = localStorage.getItem("jwt");
  if (userCredential === null) {
    return invalidSessionJumpBack();
  }

  const userResponse = await fetch(API_USER_INFO.replace("{}", username),
    {
      method: "GET",
      headers: new Headers({
        "authorization": `Bearer ${userCredential}`
      })
    });

  if (userResponse.status === 401) {
    errorMessage("Unauthorized!");
    return setTimeout(() => {
      invalidSessionJumpBack();
    }, 2500);
  }

  if (userResponse.status === 403) {
    errorMessage("Forbidden!");
    return setTimeout(() => {
      invalidSessionJumpBack();
    }, 2500);
  }

  const userResponseBody = await userResponse.json();
  return userResponseBody.username;
}

const invalidSessionJumpBack = () => {
  errorMessage("Invalid session!");
  setTimeout(() => {
    history.push("/login");
  }, 3500);
}

export {
  getUserInfo,
  invalidSessionJumpBack
}
