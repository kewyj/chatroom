import React, { useEffect, useState } from "react";
import { useSelector, useDispatch } from "react-redux";
import { useNavigate } from "react-router-dom";

import AppState from "../store";

import RoomButton from "../components/RoomButton";
import CreateRoomButton from "../components/CreateRoomButton";

import "@fortawesome/fontawesome-free/css/all.css";
import "../styles/roomlist.css";

export interface ChatProps {}

// Define the type for your store state
interface AppState {
  username: string;
}

interface HomeClock {
  date: string;
  time: string;
}

// Define the type for Chatrooms
interface Chatrooms {
  chatroomID: string;
  users: number;
}

// to give a color to username
const predefinedColors = [
  "BlueViolet",
  "DeepPink",
  "Coral",
  "CornflowerBlue",
  "Crimson",
  "DarkOrange",
  "DodgerBlue",
  "Magenta",
  "MediumPurple",
  "RebeccaPurple",
  "DarkSeaGreen",
  "MediumSlateBlue",
  "OliveDrab",
];

const RoomListPage: React.FunctionComponent<ChatProps> = () => {
  const dispatch = useDispatch();

  const users = localStorage.getItem('user');
  const username = users ? JSON.parse(users) : '';
  
  //const username = useSelector((state: AppState) => state.username);
  const [currentClock, setClock] = useState<HomeClock>();
  // To update and edit when database is up
  const [receivedChatrooms, setReceivedChatrooms] = useState<Chatrooms[]>([]);
  const [zoomIn, setZoomIn] = useState(false);

  // when user comes here check if have userid, dont have, navigate to first page
  const navigate = useNavigate();

  const host = "localhost";
  const port = 3333;

  const updateClock = (data: HomeClock) => {
    setClock((prev) => {
      return data;
    });
  };

  const getTime = async () => {
    try {
      var today = new Date();
      var homeClock: HomeClock = {
        date:
          today.getFullYear() +
          "-" +
          (today.getMonth() + 1) +
          "-" +
          today.getDate(),
        time:
          today.getHours() +
          ":" +
          today.getMinutes() +
          ":" +
          today.getSeconds(),
      };
      updateClock(homeClock);
    } catch (error) {
      console.error("Error fetching messages:", error);
    }
  };

  useEffect(() => {
    const timer = setTimeout(() => {
      setZoomIn(true);
    }, 2000);

    return () => clearTimeout(timer);
    
  }, []);
  
  useEffect(() => {
    getTime();

    // Make interval every 1 sec
    const intervalId = setInterval(getTime, 1000);

    // Clear interval
    return () => {
      clearInterval(intervalId);
    };
  }, []);

  // To update and edit when database is up
  useEffect(() => {
    // Fetch messages from the server and update receivedMessages state
    fetchChatroomsFromServer();

    // Update chatrooms every 5 seconds
    const intervalId = setInterval(fetchChatroomsFromServer, 5000);

    // Clear interval
    return () => {
        clearInterval(intervalId);
    };
  }, []);

  //To update and edit when database is up
  const fetchChatroomsFromServer = async () => {
    try {
      const roomsResponse =  await fetch(`https://1bs9qf5xn1.execute-api.ap-southeast-1.amazonaws.com/rooms`, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
        }
      });

      //Check if server response was null before calling state change
      if (!roomsResponse.ok) {
          throw new Error('Failed to fetch response from /chat');
      }

      const roomsData = await roomsResponse.json();

      // UPDATE AND RENDER THE AVAILABLE CHATROOMS + USERS
      if (Array.isArray(roomsData)) {
        const constructChatrooms = roomsData.map(room => ({
          chatroomID: room.chatroom_id,
          users: room.num_users
        }))
        setReceivedChatrooms(constructChatrooms);
      }

    } catch (error) {
      console.error("Error fetching messages:", error);
    }
  };

  return (
    <div>
      <div className={`animate_roomlist_background ${zoomIn ? 'zoomInBackground' : ''}`}>
      </div>
        <main className="room_background">
          <section className="container">
            <div className="row p-3" id="clock">
              <div className="col-lg-4">
                <p>{currentClock?.date}</p>
                <p>{currentClock?.time}</p>
              </div>
            </div>
            <div className="row p-3" id="greetings">
              <div className="col-lg-4 d-flex">
                <p>Hello, {username}</p>
              </div>
            </div>
          </section>
          <section className="container">
            <div className="row p-3" id="create_room">
              <div className="col-lg-4 d-flex">
                <CreateRoomButton>New Chatroom</CreateRoomButton>
              </div>
            </div>
            <div className="row p-3" id="rooms_list">
              <div className="col-lg-12 d-flex flex-column align-items-stretch flex-shrink-8">
                {receivedChatrooms.map((chatroom) => (
                  <RoomButton key={chatroom.chatroomID}
                    title={chatroom.chatroomID}
                    users={chatroom.users}
                  />
              ))}
              </div>
            </div>
          </section>
          <div className="animate_roomlist_background"></div>
          </main>
    </div>
  );

  // // To update and edit when database is up
  // return (
  //   <main className="room_background">
  //     <section className="container">
  //       <div className="row p-3" id="clock">
  //         <div className="col-lg-4">
  //           <p>{currentClock?.date}</p>
  //           <p>{currentClock?.time}</p>
  //         </div>
  //       </div>
  //       <div className="row p-3" id="greetings">
  //         <div className="col-lg-4 d-flex">
  //           <p>Hello, {username}</p>
  //         </div>
  //       </div>
  //     </section>
  //     <section className="d-flex flex-column align-items-stretch flex-shrink-8">
  //       {receivedChatrooms.map((chatroom) => (
  //         <RoomButton
  //           key={chatroom.chatroomID}
  //           title={`Cabin ${chatroom.chatroomID}`}
  //           users={chatroom.users}
  //         />
  //       ))}
  //     </section>
  //   </main>
  // );
};

export default RoomListPage;
