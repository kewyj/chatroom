import { createStore } from 'redux';
// Define an initial state value for the app
const initialState = {
  userID: '',
  msg: ''
}

function reducer(state = initialState, action) {
  switch (action.type) {
    case 'set_user_id':
      return { ...state, userID: action.payload };
    default:
      return state;
  }
}

const store = createStore(reducer);

export default store;