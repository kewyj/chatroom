import React, {useEffect, useState} from 'react';
import "@fortawesome/fontawesome-free/css/all.css";
import AppState from '../store'
import { useSelector, useDispatch } from 'react-redux';
import { setMessage } from '../actions'
import { TimerExample } from '../SpamTimer'
import { useNavigate } from 'react-router-dom';
import './styles.css'

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

const ChatPage: React.FunctionComponent<ChatProps> = () => {
    const dispatch = useDispatch();
    const userID = useSelector((state: AppState) => state.userID);
    const message = useSelector((state: AppState) => state.message);
    const [receivedMessages, setReceivedMessages] = useState<Message[]>([]);
    const usernameToSend = userID ? userID.username : '';
    const [messageLimit, _setMessageLimit] = useState<number>(32); // added underscore to remove warning
    const [isVisible, setVisibility] = useState<boolean>(false);
    const [isGlittering, setIsGlittering] = useState<boolean>(false);

    // when user comes here check if have userid, dont have, navigate to first page
    const navigate = useNavigate();

    useEffect(() => {
        const handleBeforeUnload = (event: BeforeUnloadEvent) => {

            exitToServer();

            const exitConfirmation = 'Leaving so soon? Chat will be lost.';
            event.returnValue = exitConfirmation;
            return exitConfirmation;
        };

        window.addEventListener('beforeunload', handleBeforeUnload);

        return () => {
            window.removeEventListener('beforeunload', handleBeforeUnload);
        };
    }, [usernameToSend]);

    useEffect(() => {
        // Fetch messages from the server and update receivedMessages state
        fetchMessagesFromServer();

        // Make interval every 0.1 sec
        const intervalId = setInterval(fetchMessagesFromServer, 100);

        // Clear interval
        return () => {
            clearInterval(intervalId);
        };

    }, []);

    // placing usernameToSend and navigate under the [] meant that this useEffect() function will run whenever either usernameToSend or navigate changes
    useEffect(() => {
        if (!usernameToSend) {
            navigate('/');
        }
    }, [usernameToSend, navigate]);

    const send = async (event: React.FormEvent) => {
        event.preventDefault();

        const dataToSend = {
            username: usernameToSend,
            content: message
        };

        //console.log('Data to send:', dataToSend);

        try {
            if (!isWhitespace(dataToSend.content)) {
                // Send message to the server
                const response = await fetch('http://localhost:3333/chat', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(
                        dataToSend)
                });

                // Clear the input field after sending message
                dispatch(setMessage(''));

                if (!response.ok) {
                    throw new Error('Failed to fetch messages');
                }

                const data = await response.json();

                if (!data)
                    return;

                // Update receivedMessages state with the messages received from the server
                if (data.username)
                {
                    warnUsers();
                    timer.setIsVisibleTrue(setVisibility);
                    setIsGlittering(true); 
                }
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
    });
    }

    const warnUsers = () => {
        setVisibility(true);
    }

    const exitToServer = async () => {

        const dataToSend = {
        username: usernameToSend
        };

        console.log("sending exit to server")

        try {
            await fetch('http://localhost:3333/exit', {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(
                dataToSend)
            });

            console.log(dataToSend);

            const exitConfirmation = 'Leaving so soon?';
            return exitConfirmation;
        }
        catch (error) {
            console.error('Error sending message:', error);
            throw error;
        }
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

    return (
        <div className="chat-background">
            <div className="d-flex flex-column align-items-stretch flex-shrink-8">
                <div className={`d-flex align-items-center flex-shrink-8 p-3 link-dark text-decoration-none border-bottom ${isGlittering? 'username-glitter' : ''}`}>
                    <input className="fs-5 fw-semibold" style={{ backgroundColor: 'mediumorchid', fontSize: '20px', color: 'white' }} value={userID?.username?.substring(0, 4) || ''} readOnly />
                </div>
                {receivedMessages && receivedMessages.length > 0 && (
                    <div className="messages-container">
                        {receivedMessages.map((msg, index) => {
                            return (
                                <div key={index} className="message">
                                    <strong>{msg.username}: </strong>{msg.content}
                                </div>
                            );
                        })}
                    </div>
                )}
            </div>
            <div className='inputContainer'>
                <form onSubmit={send}>
                    <input className="form-control" placeholder="Say something..." value={message} onChange={handleInputChange} disabled={isVisible} />
                </form>
                {isVisible &&
                    <strong className='warning'>WARNING : You are spamming!</strong>
                }
            </div>
        </div>
    )
}; 

export default ChatPage;