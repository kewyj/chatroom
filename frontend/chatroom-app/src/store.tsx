import { createStore } from 'redux';
// Define an initial state value for the app
const initialState = {
  value: ''
}

function getUserID(state = initialState, action) {
  switch (action.type) {
    case 'get/newuserID':
      return { ...state, value: action.payload };
    default:
      return state;
  }
}

const store = createStore(getUserID);

export default store;