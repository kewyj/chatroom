import { useState } from 'react';
import { useNavigate } from "react-router-dom";
import { useDispatch } from "react-redux";
import '../styles/components/loginButton.css'
import config from '../config.json';
import Axios from "axios";
import { GET_USER_ID } from "../action_types";
import 'bootstrap/dist/css/bootstrap.min.css';
import { setUsername } from '../actions';
import React from "react";
// Set action type
//const GET_NEW_USER_ID = "set_user_id";

interface Props {
    children: string;
    // '?' gives color a value is optional, similar to default param in c++
    // this allows for color to only represent either primary secondary or danger
    //color?: 'primary' | 'secondary' | 'danger';
    //color?: string;
    //onClick: () => void;
}

const LoginButton = ({ children }: Props) => {
  const [username, setUsernameState] = useState("");
  const navigate = useNavigate();
  const dispatch = useDispatch();
  
  const handleUsernameChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setUsernameState(event.target.value);
  }

    const handleClick = async () => {
      try {
        if (!username.trim()) {
          alert("Please enter a username.");
          return;
        }
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

        navigate('/rooms');
  
      }
      catch (error) {
          console.error("Error fetching data:", error);
      }
    };

  return (
    <div className='parent-container'>
      <input 
      type="text"
      placeholder="Enter your username"
      className="form-control-input"
      value={username}
      onChange={handleUsernameChange}
      />
    <button className='btn btn-primary btn-sx' onClick={handleClick}>
      {children}
    </button>
  </div>
  )
}

export default LoginButton