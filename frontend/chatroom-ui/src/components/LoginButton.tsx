import { useState, useRef } from 'react';
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
  const inputRef = useRef<HTMLInputElement>(null);
  
  const handleUsernameChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setUsernameState(event.target.value);
  }

    const handleClick = () => {
    try {
      if (!username.trim()) {
        alert("Please enter a username.");
        setUsernameState("")
        if (inputRef.current) {
          inputRef.current.focus();
        }
        return;
      }
      
      // Set username in the Redux store
      dispatch(setUsername(username.trim()));

      // Navigate to '/rooms' route
      navigate('/rooms');
  
    }
    catch (error) {
        console.error("Error:", error);
    }
  };

  return (
    <div className='input_container'>
      <input 
      ref= {inputRef}
      type="text"
      placeholder="Enter your username"
      className="input"
      value={username}
        onChange={handleUsernameChange}
        maxLength={28}
      />
      <div className='button_container'>
        <button className='button' onClick={handleClick}>
          {children}
        </button>
        <div></div>
      </div>
    </div>
  )
}

export default LoginButton