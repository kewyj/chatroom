import Axios from "axios";
import React from "react";
import { useState } from 'react';
import { useNavigate } from "react-router-dom";
import { useDispatch } from "react-redux";
import { GET_USER_ID } from "../action_types";
import { setUsername } from '../actions';

import config from '../config.json';

import 'bootstrap/dist/css/bootstrap.min.css';
import '../styles/components/createRoomButton.css'

interface Props {
  children: string;
}

const CreateRoomButton = ({ children }: Props) => {
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
    <button id="create_new" onClick={handleClick}>
      {children}
    </button>
  )
}

export default CreateRoomButton