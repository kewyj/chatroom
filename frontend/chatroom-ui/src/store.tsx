import { createStore } from 'redux';
import { SET_MESSAGE } from './action_types';
import { GET_USER_ID } from './action_types';
import { RESET_USER } from './action_types';
import { SET_USERNAME } from './action_types';

// Define an interface for the action
interface Action {
  type: string;
  payload: any;
}

// Define an initial state value for the app
const initialState = {
  userID: null,
  message: '',
  username: ''
}

function reducer(state = initialState, action : Action) {
  switch (action.type) {
    case GET_USER_ID:
      return { ...state, userID: action.payload };
    case RESET_USER:
      return { userID: null, message: '', username: ''};
    case SET_MESSAGE:
      return { ...state, message: action.payload };
    case SET_USERNAME:
      return { ...state, username: action.payload };
    default:
      return state;
  }
}

const store = createStore(reducer);

export default store;