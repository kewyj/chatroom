import Axios from "axios";
import React from "react";
import { useState } from 'react';
import { useNavigate } from "react-router-dom";
import { useDispatch } from "react-redux";
import { GET_USER_ID } from "../action_types";
import { setUsername } from '../actions';

import config from '../config.json';

import 'bootstrap/dist/css/bootstrap.min.css';
import '../styles/components/roomButton.css'

interface Props {
  children: {
    id: number;
    title: string;
    users: number;
  }
}

const RoomButton = ({ children : { id, title, users }}: Props) => {
  const [username, setUsernameState] = useState("");
  const navigate = useNavigate();
  const dispatch = useDispatch();
  
  const handleUsernameChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setUsernameState(event.target.value);
  }

  const handleClick = async () => {
    try {
      // to get userID
      const serverHost = config.server.host;
      const serverPort = config.server.port;
      const path = "/newuser";
      const url = `http://${serverHost}:${serverPort}${path}`;
      
      // PUT request to server
      const response = await Axios.put(url);
      const newUserID = response.data;

      // Dispatch an action to update the store with new userID
      dispatch({ type: GET_USER_ID, payload: newUserID });

      navigate('/chat');
    }
    catch (error) {
        console.error("Error fetching data:", error);
    }
  };

  return (
    <button className="enter_room" onClick={handleClick}>
      <div>
        <caption>Chapter {id}</caption>
        <h4>{title}</h4>
        <p>{users} active users</p>
      </div>
    </button>
  )
}

export default RoomButton

// import Axios from "axios";
// import React from "react";
// import { useState } from "react";
// import { useNavigate } from "react-router-dom";
// import { useDispatch } from "react-redux";
// import { GET_USER_ID } from "../action_types";
// import { setUsername } from "../actions";

// import config from "../config.json";

// import "bootstrap/dist/css/bootstrap.min.css";
// import "../styles/components/roomButton.css";

// interface RoomButtonProps {
//   title: string;
//   users: number;
// }

// const RoomButton: React.FC<RoomButtonProps> = ({ title, users }) => {
//   const [username, setUsernameState] = useState("");
//   const navigate = useNavigate();
//   const dispatch = useDispatch();

//   const handleUsernameChange = (event: React.ChangeEvent<HTMLInputElement>) => {
//     setUsernameState(event.target.value);
//   };

//   const handleClick = async () => {
//     try {
//       // to get userID
//       const serverHost = config.server.host;
//       const serverPort = config.server.port;
//       const path = "/newuser";
//       const url = `http://${serverHost}:${serverPort}${path}`;

//       // PUT request to server
//       const response = await Axios.put(url);
//       const newUserID = response.data;

//       // Dispatch an action to update the store with new userID
//       dispatch({ type: GET_USER_ID, payload: newUserID });

//       navigate("/chat");
//     } catch (error) {
//       console.error("Error fetching data:", error);
//     }
//   };

//   return (
//     <button className="room-button" onClick={handleClick}>
//       <div>
//         <h4>{title}</h4>
//         <p>{users} active users</p>
//       </div>
//     </button>
//   );
// };

// export default RoomButton;
