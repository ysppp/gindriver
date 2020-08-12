import { message } from 'antd';

const loadingMessage = (key: string) => {
  message.loading({ content: 'Loading...', key });
};

const successMessage = (key: string, msg: string) => {
  message.success({content: msg, key, duration: 2});
}

const errorMessage = (msg: string) => {
  message.destroy();
  message.error(msg);
}

const closeMessage = () => {
  message.destroy();
}

export {
  loadingMessage,
  successMessage,
  errorMessage,
  closeMessage
}
