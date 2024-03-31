import { createStore, Action } from 'redux';
import { SET_MESSAGE } from './action_types';
import { GET_USER_ID } from './action_types';
import { RESET_USER } from './action_types';
import { SET_USERNAME } from './action_types';
import { SET_CHATROOM_ID } from './action_types';

interface State {
  userID: string | null;
  message: string;
  username: string;
  chatID: string;
}

// Define an initial state value for the app
const initialState: State = {
  userID: null,
  message: '',
  username: '',
  chatID: '',
}

// Define an interface for the action
type ActionTypes =
  | { type: typeof GET_USER_ID; payload: string }
  | { type: typeof RESET_USER }
  | { type: typeof SET_MESSAGE; payload: string }
  | { type: typeof SET_USERNAME; payload: string }
  | { type: typeof SET_CHATROOM_ID; payload: string };

function reducer(state = initialState, action : ActionTypes) : State {
  switch (action.type) {
    case GET_USER_ID:
      console.log("New userID:", action.payload);
      return { ...state, userID: action.payload };
    case RESET_USER:
      return { ...state, userID: null, message: ''};
    case SET_MESSAGE:
      return { ...state, message: action.payload };
    case SET_USERNAME:
      return { ...state, username: action.payload };
    case SET_CHATROOM_ID:
      return { ...state, chatID: action.payload };
    default:
      return state;
  }
}

const store = createStore(reducer);

export default store;