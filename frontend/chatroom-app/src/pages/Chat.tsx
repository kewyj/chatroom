import React, {useEffect, useState} from 'react';
import "@fortawesome/fontawesome-free/css/all.css";
import AppState from '../store'
import { useSelector, useDispatch } from 'react-redux';
import {setMessage} from '../actions'

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

const ChatPage: React.FunctionComponent<ChatProps> = () => {

    const dispatch = useDispatch();
    const userID = useSelector((state: AppState) => state.userID);
    const message = useSelector((state: AppState) => state.message);
    const [receivedMessages, setReceivedMessages] = useState<Message[]>([]);
    const usernameToSend = userID ? userID.username : '';
    const [messageLimit, setMessageLimit] = useState<number>(28);
    const [isVisible, setVisibility] = useState<boolean>(false);

    useEffect(() => {
        // Fetch messages from the server and update receivedMessages state
        fetchMessagesFromServer();
        //checkSpamFromServer();

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
            if (!isWhitespace(dataToSend.content)) {
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
            }

        } catch (error) {
            console.error('Error sending message:', error);
        }
    }

    const handleInputChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        const { value } = event.target;
        dispatch(setMessage(value));
    }

    const updateReceivedMessages = (data: Message[]) => {
        setReceivedMessages(prevMessages => {
            const combinedMessages = [...prevMessages, ...data];
            const newMessages = combinedMessages.slice(-messageLimit);
            return newMessages;
        })
    }

    // const warnUsers = () => {
    //     setVisibility(true);
    // }

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

            //console.log(data)

            // CHEck if server response was null before calling state change
            if (!data)
                return;

            //console.log(data)

            // Update receivedMessages state with the messages received from the server
            updateReceivedMessages(data);
        }
        catch (error) {
            console.error('Error fetching messages:', error);
        }
    }

    // const checkSpamFromServer = async () => {
    //     try {
    //         // change accordingly
    //         const url = 'http://localhost:3333/poll';

    //         const response = await fetch(url, {
    //             method: 'GET',
    //             headers: {
    //                 'Content-Type': 'application/json',
    //             }
    //         });
    //         if (!response.ok) {
    //             throw new Error('Failed to fetch messages');
    //         }
    //         const data = await response.json();

    //         // change accordingly
    //         if (data.content == "spam")
    //             warnUsers();
    //     }
    //     catch (error) {
    //         console.error('Error fetching messages:', error);
    //     }
    // }

    return (
        <div className="container">
            <div className="d-flex flex-column align-items-stretch flex-shrink-8 bg-white">
                <div className="d-flex align-items-center flex-shrink-8 p-3 link-dark text-decoration-none border-bottom">
                    <input className="fs-5 fw-semibold" value={userID?.username?.substring(0, 4) || ''} readOnly />
                </div>
                {receivedMessages && receivedMessages.length > 0 && (
                    <div className="messages-container">
                        {receivedMessages.map((msg, index) => {
                            console.log(msg.content);
                            return (
                                <div key={index} className="message">
                                    <strong>{msg.username}: </strong>{msg.content}
                                </div>
                            );
                        })}
                    </div>
                )}
            </div>
            <form onSubmit={send}>
                <input className="form-control" placeholder="Say something..." value={message} onChange={handleInputChange} />
            </form>
            {isVisible && <div className="spam-message"> 
                <strong>WARNING : You are spamming!</strong>
            </div>}
        </div>
    )
}; 

export default ChatPage;