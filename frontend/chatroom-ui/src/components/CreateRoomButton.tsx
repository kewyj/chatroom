import Axios from "axios";
import React from "react";
import { useState } from 'react';
import { useNavigate } from "react-router-dom";
import { useSelector, useDispatch } from "react-redux";
import { GET_USER_ID, SET_CHATROOM_ID } from "../action_types";
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

      const response = await fetch(`https://1bs9qf5xn1.execute-api.ap-southeast-1.amazonaws.com/newroom`, {
          method: 'PUT',
          headers: {
              'Content-Type': 'application/json',
          }
      });

      const data = await response.json(); 
      console.log(data.chatroom_id)

      const dataToSendNewUser = {
        custom_username: customUsername
      }

      const newUserResponse = await fetch(`https://1bs9qf5xn1.execute-api.ap-southeast-1.amazonaws.com/newuser`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(
          dataToSendNewUser)
      });
      
      const newUserID = await newUserResponse.json();
      console.log(`${newUserID.user_uuid} is the user's uuid`);

      const dataToSendAddToRoom = {
        chatroom_id: data.chatroom_id,
        user_uuid: newUserID.user_uuid
      };

      console.log(`Sending ${dataToSendAddToRoom.chatroom_id} to addtoroom`);
      console.log(`Sending ${dataToSendAddToRoom.user_uuid} to addtoroom`);

      await fetch(`https://1bs9qf5xn1.execute-api.ap-southeast-1.amazonaws.com/addtoroom`,
        {
          method: "PUT",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify
            (dataToSendAddToRoom)
        });

      // Dispatch an action to update the store with new userID
      dispatch({ type: GET_USER_ID, payload: newUserID.user_uuid });
      dispatch({ type: SET_CHATROOM_ID, payload: data.chatroom_id });

      navigate('/chat');
    }
    catch (error) {
        console.error("Error fetching data:", error);
    }
  };

  return (
    <button id="create_new" className="button" onClick={handleClick}>
      {children}
    </button>
  )
}

export default CreateRoomButton