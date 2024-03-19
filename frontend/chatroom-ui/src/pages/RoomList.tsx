import React, {useEffect, useState} from 'react';
import "@fortawesome/fontawesome-free/css/all.css";
import AppState from '../store'
import { useSelector, useDispatch } from 'react-redux';
import { setMessage } from '../actions'
import { TimerExample } from '../SpamTimer'
import { unstable_useViewTransitionState, useNavigate } from 'react-router-dom';
import RoomButton from '../components/RoomButton';
import { useLocation } from 'react-router-dom';
import '../styles/roomlist.css'
import { createBrowserHistory, Update } from 'history';

export interface ChatProps { }

// Define the type for your store state
interface AppState {
  userID: {
    username: string;
    } | null;
  message: string;
}

interface HomeClock {
  date: string;
  time: string;
}

// Define the type for Message
interface Message {
  username: string;
  content: string;
}

const isWhitespace = (str: string): boolean => {
    return /^\s*$/.test(str);
}

const timer = new TimerExample;

// to give a color to username
const predefinedColors = ['BlueViolet', 'DeepPink', 'Coral', 'CornflowerBlue', 'Crimson', 'DarkOrange', 'DodgerBlue', 'Magenta', 'MediumPurple', 'RebeccaPurple', 'DarkSeaGreen', 'MediumSlateBlue', 'OliveDrab'];

// map username to color
const usernameColors: { [key: string]: string } = {};

const getRandomColor = () => {
    const randomIndex = Math.floor(Math.random() * predefinedColors.length);
    return predefinedColors[randomIndex];
}

const RoomListPage: React.FunctionComponent<ChatProps> = () => {
    const dispatch = useDispatch();
    const userID = useSelector((state: AppState) => state.userID);
    const usernameToSend = userID ? userID.username : '';
    const [currentClock, setClock] = useState<HomeClock>();

    // when user comes here check if have userid, dont have, navigate to first page
    const navigate = useNavigate();

    const host = "localhost"
    const port = 3333

    // placing usernameToSend and navigate under the [] meant that this useEffect() function will run whenever either usernameToSend or navigate changes
    useEffect(() => {
        if (!usernameToSend) {
            navigate('/');
        }
    }, [usernameToSend, navigate]);
    
    useEffect(() => {
        if (usernameToSend) {
            // assign a color to the username for display
            usernameColors[usernameToSend] = getRandomColor();
        }
    }, [usernameToSend]);

    const updateClock = (data: HomeClock) => {
        setClock(prev => {
            return data;
        });
    }

    const getTime = async () => {
        try {
            var today = new Date();
            var homeClock: HomeClock = {
                date: today.getFullYear() + '-' + (today.getMonth() + 1) + '-' + today.getDate(),
                time: today.getHours() + ":" + today.getMinutes() + ":" + today.getSeconds()
            };
            updateClock(homeClock);
        }
        catch (error) {
            console.error('Error fetching messages:', error);
        }
    }

    useEffect(() => {
        getTime();

        // Make interval every 1 sec
        const intervalId = setInterval(getTime, 1000);

        // Clear interval
        return () => {
            clearInterval(intervalId);
        };
    }, []);

    return (
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
                        <p>Hello, {userID?.username?.substring(0, 10) || ''}{(userID?.username && userID?.username?.length > 10) ? "..." : ""}</p>
                    </div>
                </div>
            </section>
            <section className="d-flex flex-column align-items-stretch flex-shrink-8">
                <RoomButton>{{title: "Chatroom Name 1", users: 120}}</RoomButton>
                <RoomButton>{{title: "Chatroom Name 2", users: 30}}</RoomButton>
                <RoomButton>{{title: "Chatroom Name 3", users: 12}}</RoomButton>
                <RoomButton>{{title: "Chatroom Name 4", users: 8}}</RoomButton>
            </section>
        </main>
    )
}; 

export default RoomListPage;