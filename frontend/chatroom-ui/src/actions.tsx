import { SET_MESSAGE } from './action_types';
import { GET_USER_ID } from './action_types';
import { SET_USERNAME } from './action_types';
import { RESET_USER } from './action_types';

export const setMessage = (message : string) => ({
  type: SET_MESSAGE,
  payload: message,
});

export const getUserID = (userID : string) => ({
  type: GET_USER_ID,
  payload: userID,
});

export const setUsername = (username: string) => ({
  type: SET_USERNAME,
  payload: username
});

export const resetUser = () => ({
  type: RESET_USER,
});
