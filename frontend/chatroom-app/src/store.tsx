import { createStore } from 'redux';
import { SET_MESSAGE } from './action_types';
import { GET_USER_ID } from './action_types';

// Define an interface for the action
interface Action {
  type: string;
  payload: any;
}

// Define an initial state value for the app
const initialState = {
  userID: null,
  message: ''
}

function reducer(state = initialState, action : Action) {
  switch (action.type) {
    case GET_USER_ID:
      return { ...state, userID: action.payload };
    case SET_MESSAGE:
      return { ...state, message: action.payload };
    default:
      return state;
  }
}

const store = createStore(reducer);

export default store;