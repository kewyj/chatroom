import { useState } from 'react';
import { useNavigate } from "react-router-dom";
import { useDispatch } from "react-redux";
import config from '../config.json';
import Axios from "axios";
import { GET_USER_ID } from "../action_types";
import 'bootstrap/dist/css/bootstrap.min.css';
import { setUsername } from '../actions';
import React from "react";

import '../styles/components/roomButton.css'

interface Props {
  children: {
    title: string;
    users: number;
  }
}

const RoomButton = ({ children : { title, users }}: Props) => {
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
      dispatch(setUsername(username));

      navigate('/chat');
    }
    catch (error) {
        console.error("Error fetching data:", error);
    }
  };

  return (
    <button className='' onClick={handleClick}>
      <div>
        <h4>{title}</h4>
        <p>{users} active users</p>
      </div>
    </button>
  )
}

export default RoomButton