import { errorMessage } from './components/message';
import { supported } from '@github/webauthn-json';

(() => {
  if(!supported()) {
    errorMessage("This browser does not support WebAuthn!");
    setTimeout(() => {
      errorMessage("Please use compatible browers or operating system.");
    }, 3000);
  }
})();
