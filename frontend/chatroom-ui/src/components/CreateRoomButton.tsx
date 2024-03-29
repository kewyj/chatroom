import Axios from "axios";
import React from "react";
import { useState } from 'react';
import { useNavigate } from "react-router-dom";
import { useSelector, useDispatch } from "react-redux";
import { GET_USER_ID } from "../action_types";
import { setUsername } from '../actions';

import config from '../config.json';
import AppState from "../store";

import 'bootstrap/dist/css/bootstrap.min.css';
import '../styles/components/createRoomButton.css'

// Define the type for your store state
interface AppState {
  username: string;
}

interface Props {
  children: string;
}

const CreateRoomButton = ({ children }: Props) => {
  const [username, setUsernameState] = useState("");
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const customUsername = useSelector((state: AppState) => state.username);
  
  const handleUsernameChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setUsernameState(event.target.value);
  }

  const handleClick = async () => {
    try {
      // to get userID
      // const serverHost = config.server.host;
      // const serverPort = config.server.port;
      // const nrPath = "/newroom";
      // const nrUrl = `http://${serverHost}:${serverPort}${nrPath}`;

      // PUT /newroom to create room (implement when database is up)
      const response = await fetch(`https://1bs9qf5xn1.execute-api.ap-southeast-1.amazonaws.com/newroom`, {
          method: 'PUT',
          headers: {
              'Content-Type': 'application/json',
          }
      });

      const data = await response.json(); 
      console.log(data.chatroom_id)

      const dataToSend = {
        custom_username: customUsername,
        chatroom_id: data.chatroom_id
      }

      const newUserResponse = await fetch(`https://1bs9qf5xn1.execute-api.ap-southeast-1.amazonaws.com/newuser`, {
          method: 'PUT',
          headers: {
              'Content-Type': 'application/json',
          },
          body: JSON.stringify(
              dataToSend)
      });

      const newUserID = await newUserResponse.json();
      console.log(newUserID.user_uuid);

      // Dispatch an action to update the store with new userID
      //dispatch({ type: GET_USER_ID, payload: newUserID });

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