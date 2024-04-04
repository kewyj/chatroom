import React from "react";
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useSelector, useDispatch } from "react-redux";
import { GET_USER_ID, SET_CHATROOM_ID } from "../action_types";

import "bootstrap/dist/css/bootstrap.min.css";
import "../styles/components/roomButton.css";

interface AppState {
  userID: string;
  username: string;
}

interface Props {
  title: string;
  users: number;
}

const RoomButton: React.FC<Props> = ({ title, users }) => {
  const [username, setUsernameState] = useState("");
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const userID = useSelector((state: AppState) => state.userID);
  const customUsername = useSelector((state: AppState) => state.username);

  const handleUsernameChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setUsernameState(event.target.value);
  };

  const handleClick = async () => {
    try
    {
      const dataToSendNewUser = {
        custom_username: customUsername
      }

      if (userID == null)
      {
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
      
        const dataToSend = {
          chatroom_id: title,
          user_uuid: newUserID.user_uuid
        };
        
        await fetch(`https://1bs9qf5xn1.execute-api.ap-southeast-1.amazonaws.com/addtoroom`,
        {
          method: "PUT",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(
            dataToSend)
        });
        
        dispatch({ type: GET_USER_ID, payload: newUserID.user_uuid });
      }
      else
      {
        const dataToSend = {
          chatroom_id: title,
          user_uuid: userID
        };

        await fetch(`https://1bs9qf5xn1.execute-api.ap-southeast-1.amazonaws.com/addtoroom`,
        {
          method: "PUT",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(
            dataToSend)
        });
      }

      dispatch({ type: SET_CHATROOM_ID, payload: title });
      navigate("/chat");

    } catch (error) {
      console.error("Error fetching data:", error);
    }
  };

  return (
    <button className="enter_room" onClick={handleClick}>
      <div>
        <h4>Chatroom {title}</h4>
        <p>{users} active users</p>
      </div>
    </button>
  );
};

export default RoomButton;
