import React, {useEffect, useState} from 'react';
import "@fortawesome/fontawesome-free/css/all.css";
import AppState from '../store.tsx'
import { useSelector, useDispatch } from 'react-redux';
import {setMessage} from '../actions'

export interface ChatProps { }

const ChatPage: React.FunctionComponent<ChatProps> = (props) => {

    const dispatch = useDispatch();
    const userID = useSelector((state: AppState) => state.userID);
    const message = useSelector((state: AppState) => state.message);
    const [receivedMessages, setReceivedMessages] = useState([]);
    const usernameToSend = AppState.getState().userID ? AppState.getState().userID.username : '';

    useEffect(() => {
        // Fetch messages from the server and update receivedMessages state
        fetchMessagesFromServer();

        // Make interval every 1 sec
        const intervalId = setInterval(fetchMessagesFromServer, 1000);

        // Clear interval
        return () => clearInterval(intervalId);
    }, []);

    const send = async (event: React.FormEvent) => {
        event.preventDefault();

        const dataToSend = {
            username: usernameToSend,
            content: message
        };

        console.log('Data to send:', dataToSend);

        try {
            // Send message to the server
            await fetch('http://localhost:3333/chat', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(
                    dataToSend)
            });

            // Clear the input field after sending message
            dispatch(setMessage(''));
        } catch (error) {
            console.error('Error sending message:', error);
        }
    }

    const handleInputChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        const { value } = event.target;
        dispatch(setMessage(value));
    }

const fetchMessagesFromServer = async () => {
    try {
        const url = 'http://localhost:3333/poll';
    
        const dataforPatch = {
            username: usernameToSend
        };

        const response = await fetch(url, {
            method: 'PATCH',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(dataforPatch)
        });
        if (!response.ok) {
            throw new Error('Failed to fetch messages');
        }
        const data = await response.json();

        console.log(data)

        // CHEck if server response was null before calling state change
        if (!data.messages || data.messages.length == 0)
            return;

        // Update receivedMessages state with the messages received from the server
        setReceivedMessages(data.messages);
    }
    catch (error) {
        console.error('Error fetching messages:', error);
    }
}

    return (
        <div className="container">
            <div className="d-flex flex-column align-items-stretch flex-shrink-8 bg-white">
                <div className="d-flex align-items-center flex-shrink-8 p-3 link-dark text-decoration-none border-bottom">
                    <input className="fs-5 fw-semibold" value={userID?.username?.substring(0, 4) || ''} readOnly />
                </div>
            </div>
            <form onSubmit={send}>
                <input className="form-control" placeholder="Say something..." value={message} onChange={handleInputChange} />
            </form>
            <div>
                {receivedMessages.map((msg, index) => (
                    <div key={index}>
                        <div>{msg.username}</div>
                        <div>{msg.content}</div>
                    </div>
                ))}
            </div>
        </div>
    )
}; 

export default ChatPage;