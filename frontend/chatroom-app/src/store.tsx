import { createStore } from 'redux';
import { SET_MESSAGE } from './action_types';

// Define an initial state value for the app
const initialState = {
  userID: '',
  message: ''
}

function reducer(state = initialState, action) {
  switch (action.type) {
    case 'set_user_id':
      return { ...state, userID: action.payload };
    case SET_MESSAGE:
      return { ...state, message: action.payload };
    default:
      return state;
  }
}

const store = createStore(reducer);

export default store;