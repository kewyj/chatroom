import React, {useEffect, useState} from 'react';
import "@fortawesome/fontawesome-free/css/all.css";
import AppState from '../store'
import { useSelector, useDispatch } from 'react-redux';
import { setMessage } from '../actions'
import { TimerExample } from '../SpamTimer'
import { unstable_useViewTransitionState, useNavigate } from 'react-router-dom';
import ChatRoom from '../components/ChatRoom';
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
    const [isGlittering, setIsGlittering] = useState<boolean>(true);

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

    return (
        <main className="chat-background">
            <section>
                <div>
                    <p>Hello, <input className={`fs-5 fw-semibold ${isGlittering ? 'username-glitter' : ''}`} style={{ borderColor: 'mediumorchid', fontSize: '20px', color: usernameColors[usernameToSend] || '#000000', fontWeight: 'bold'}} value={userID?.username || ''} readOnly /></p>
                </div>
            </section>
            <section>
                <ChatRoom>{{title: "Chatroom Name 1", users: 120}}</ChatRoom>
                <ChatRoom>{{title: "Chatroom Name 2", users: 30}}</ChatRoom>
                <ChatRoom>{{title: "Chatroom Name 3", users: 12}}</ChatRoom>
                <ChatRoom>{{title: "Chatroom Name 4", users: 8}}</ChatRoom>
            </section>
            <div className="d-flex flex-column align-items-stretch flex-shrink-8">
                <div className={`d-flex align-items-center flex-shrink-8 p-3 link-dark text-decoration-none border-bottom `}>
                    <input className={`fs-5 fw-semibold ${isGlittering ? 'username-glitter' : ''}`} style={{ borderColor: 'mediumorchid', fontSize: '20px', color: usernameColors[usernameToSend] || '#000000', fontWeight: 'bold'}} value={userID?.username?.substring(0, 4) || ''} readOnly />
                </div>
            </div>
        </main>
    )
}; 

export default RoomListPage;