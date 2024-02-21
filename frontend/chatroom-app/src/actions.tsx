import { SET_MESSAGE } from './action_types';

export const setMessage = (message) => ({
  type: SET_MESSAGE,
  payload: message,
});