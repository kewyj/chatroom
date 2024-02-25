import { SET_MESSAGE } from './action_types';
import { GET_USER_ID } from './action_types';

export const setMessage = (message : string) => ({
  type: SET_MESSAGE,
  payload: message,
});

export const getUserID = (userID : string) => ({
  type: GET_USER_ID,
  payload: userID,
});
