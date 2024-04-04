import React, {useEffect, useState} from 'react';
import "@fortawesome/fontawesome-free/css/all.css";
import AppState from '../store'
import { useSelector, useDispatch } from 'react-redux';
import { setChatroomID, setMessage } from '../actions'
import { TimerExample } from '../SpamTimer'
import { unstable_useViewTransitionState, useNavigate } from 'react-router-dom';
import { useLocation } from 'react-router-dom';
import '../styles/chat.css'
import { createBrowserHistory, Update } from 'history';

export interface ChatProps { }

// Define the type for your store state
interface AppState {
    userID: string;
    message: string;
    chatID: string;
    username: string;
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
const usernameColors: { [username: string]: string } = {};

const getRandomColor = () => {
    const randomIndex = Math.floor(Math.random() * predefinedColors.length);
    return predefinedColors[randomIndex];
}

const ChatPage: React.FunctionComponent<ChatProps> = () => {
    const dispatch = useDispatch();
    const userID = useSelector((state: AppState) => state.userID);
    const message = useSelector((state: AppState) => state.message);
    let chatID = useSelector((state: AppState) => state.chatID);
    let customUsername = useSelector((state: AppState) => state.username);
    const [receivedMessages, setReceivedMessages] = useState<Message[]>([]);
    let usernameToSend = userID;
    const [messageLimit, _setMessageLimit] = useState<number>(32); // added underscore to remove warning
    const [isVisible, setVisibility] = useState<boolean>(false);
    const [isGlittering, setIsGlittering] = useState<boolean>(true);
    const [isOnce, SetIsOnce] = useState<boolean>(false);
    const [isReload, SetIsReload] = useState<boolean>(false);

    // when user comes here check if have userid, dont have, navigate to first page
    const navigate = useNavigate();


    // variables for spam feature
    const [enterKeyCount, setEnterKeyCount] = useState<number>(0);
    const maxKeyPress = 5; // 4 keypress as limit

    // extract host and port form config obj
    // const [data, setData] = useState<any | null>(null);
    
    // useEffect(() => {
    //     // Fetch and set the JSON data when the component mounts
    //     fetch(`/config.json`)
    //     .then(response => response.json())
    //     .then(jsonData => setData(jsonData.server))
    //     .catch(error => console.error('Error loading JSON:', error));
    // }, []);
    const host = "localhost"
    const port = 3333

    useEffect(() => {
        const handleBeforeUnload = (event: BeforeUnloadEvent) => {
            // indicate that exit to Server is called in local storage
            localStorage.setItem('isExitToServerCalled', JSON.stringify(true))
            const test = localStorage.getItem('isExitToServerCalled')
            console.log(`testing condition item: ${test}`);

            // save data that need to be retrieve in the event of refresh and not close
            localStorage.setItem('savedChatID', JSON.stringify(chatID))
            localStorage.setItem('savedCustomUsername', JSON.stringify(customUsername))
            localStorage.setItem('savedUsernameToSend', JSON.stringify(usernameToSend))
            
            exitToServer();
        };

        window.addEventListener('beforeunload', handleBeforeUnload);

        return () => {
            window.removeEventListener('beforeunload', handleBeforeUnload);
        };
    }, [usernameToSend]);

    const history = createBrowserHistory();
    const location = useLocation();

    useEffect(() => {
        let allowNavigation = true;
        const handlePop = (update: Update) => {
            if (allowNavigation && usernameToSend && update.action === 'POP') {
                if (window.confirm("Leaving so soon? Chat data will be lost.")) {
                    SetIsOnce(true);
                    exitRoom();
                }
                else {
                    SetIsOnce(false);
                    console.log("came");
                    if (!isReload) {
                        //history.replace('/chat');
                        history.forward();
                        SetIsReload(true);
                    }
                    //history.push('/chat');
                }
                allowNavigation = false;
            }
        };

        const handleWindowPopState = (event: PopStateEvent) => {
            if (!allowNavigation) {
                history.push('/chat')
            }
        };

        if (!isOnce) {
            history.listen(handlePop);
            window.onpopstate = handleWindowPopState;
        }

        return () => {
            window.onpopstate = null; // Cleanup event handler
        };

    }, [usernameToSend]);
    

    // useEffect(() => {
    //     console.log("Popstateevent useEffect called");

    //     let confirmedLeave = false;

    //     const handlePopState = (event: PopStateEvent) => {
    //         console.log("Popstateevent occurred");

    //         if (!confirmedLeave) {
    //             const exitConfirmation = 'Leaving so soon? Chat will be lost.';
    //             const shouldExit = window.confirm(exitConfirmation);

    //             if (!shouldExit) {
    //                 // Prevent leaving the page
    //                 event.preventDefault();
    //                 // Restore the URL to the current one
    //                 window.history.pushState(null, document.title, window.location.href);
    //             } else {
    //                 confirmedLeave = true; // Mark as confirmed
    //                 exitToServer(); // Perform exit action
    //             }
    //         }
    //     };
        
    //     window.addEventListener('popstate', handlePopState);
    //     console.log("Attached event");

    //     return () => {
    //         console.log("popstate event cleanup called");
    //         window.removeEventListener('popstate', handlePopState);
    //     };
    // }, [usernameToSend]);

    // if (window.history.length === 1 && window.location.pathname === '/') {
    //     console.log("Empty thats why cnnt work la cb")
    // }
    // else {
    //     console.log(window.history.length);
    // }

    // const location = useLocation();

    // useEffect(() => {
    //     const historyStack = [];
    //     for (let i = 0; i < window.history.length; i++) {
    //     const entry = window.history.state[i];
    //     if (entry && entry.location) {
    //         historyStack.push(entry.location.pathname);
    //     }
    //     }
    //     console.log('History stack:', historyStack);
    // }, [location]);

    useEffect(() => {
        var messages_container = document.getElementById("messages_container");
        if(messages_container != null) {
            messages_container.scrollTop = messages_container.scrollHeight;
        }

        // Fetch messages from the server and update receivedMessages state
        fetchMessagesFromServer();

        console.log(receivedMessages)

        //Make interval every 0.1 sec
        const intervalId = setInterval(fetchMessagesFromServer, 100);

        //Clear interval
        return () => {
            clearInterval(intervalId);
        };

    }, []);

    // placing usernameToSend and navigate under the [] meant that this useEffect() function will run whenever either usernameToSend or navigate changes
    useEffect(() => {
        if (localStorage.getItem('isExitToServerCalled')) {

            const storedChatID = localStorage.getItem('savedChatID');
            const storedCustomUsername = localStorage.getItem('savedCustomUsername');
            const storedUserID = localStorage.getItem('savedUsernameToSend');

            if (storedChatID !== null)
            {
                chatID = storedChatID;
            }
            if (storedCustomUsername !== null)
            {
                customUsername = storedCustomUsername;
            }
            if (storedUserID !== null)
            {
                usernameToSend = storedUserID;    
            }

            console.log(`AT ISEXITTOSERVER CALLED chatID: ${chatID}`);
            console.log(`AT ISEXITTOSERVER CALLED customUsername: ${customUsername}`);
            console.log(`AT ISEXITTOSERVER CALLED usernameToSend: ${usernameToSend}`);

            // remove all 3 stored info as no longer needed
            localStorage.removeItem('savedChatID');
            localStorage.removeItem('savedCustomUsername');
            localStorage.removeItem('savedUsernameToSend');
            // remove the boolean 
            localStorage.removeItem('isExitToServerCalled')
        }
        if (!usernameToSend) {
            navigate('/');
        }
    }, [usernameToSend, navigate]);
    
    useEffect(() => {
        if (typeof usernameToSend === 'string') {
            // assign a color to the username for display
            usernameColors[usernameToSend] = getRandomColor();
        }
    }, [usernameToSend]);

    // check for user spam
    useEffect(() => {
        // for message spam
        const timeFrame = 500 // 0.5 sec
        let timeOutID: NodeJS.Timeout | null = null; // initialized timeOutID and set to null

        const keyDown = (event: KeyboardEvent) => {
            if (event.key == 'Enter') {
                setEnterKeyCount((prevCount: number) => prevCount + 1);
                
                if (timeOutID) {
                    clearTimeout(timeOutID);
                }
                // setTimeout returns a unique identifier for the timeout
                timeOutID = setTimeout(() => {
                    // reset to 0 after every 0.5 sec
                    setEnterKeyCount(0);
                }, timeFrame);
            }
        };

        document.addEventListener('keydown', keyDown);

        return () => {
            document.removeEventListener('keydown', keyDown);
            // clear any exisiting timeout to prevent memory leak
            if (timeOutID) {
                clearTimeout(timeOutID);
            }
        };
    }, []);

    const send = async (event: React.FormEvent) => {
        event.preventDefault();

        const dataToSend = {
            chatroom_id: chatID,
            username: customUsername,
            message: message
        };

        console.log(`/chat sending chatroom: ${dataToSend.chatroom_id}`);
        console.log(`/chat sending username: ${dataToSend.username}`);
        console.log(`/chat sending message: ${dataToSend.message}`);

        try {
            if (!isWhitespace(dataToSend.message)) {

            fetch(`https://1bs9qf5xn1.execute-api.ap-southeast-1.amazonaws.com/chat`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(
                    dataToSend)
            });
                
            //response.status (CHECK IF NOT 200)
            // Clear the input field after sending message
            dispatch(setMessage(''));

            var messages_container = document.getElementById("messages_container");
            if(messages_container != null) {
                messages_container.scrollTop = messages_container.scrollHeight;
            }

            //console.log(enterKeyCount);

            if (enterKeyCount > maxKeyPress) {
                warnUsers();
                timer.setIsVisibleTrue(setVisibility);
            }

            // const data = await response.json();

            // if (!data)
            //     return;

            // // Update receivedMessages state with the messages received from the server
            // if (data.username)
            // {
            //     warnUsers();
            //     timer.setIsVisibleTrue(setVisibility);
            // }
            }

        } catch (error) {
            console.error('Error sending message:', error);
        }
    }

    const handleInputChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        const { value } = event.target;
        dispatch(setMessage(value));
    }

    const updateReceivedMessages = (data: { custom_username: string; message: string }[]) => {
        const newMessages = data.map(item => ({
            username: item.custom_username,
            content: item.message
        }));

        // Keep only the latest 'messageLimit' messages
        const limitedMessages = newMessages.slice(-messageLimit);

        setReceivedMessages(limitedMessages);
    }

    const warnUsers = () => {
        setVisibility(true);
    }

    const exitRoom = async () => {
        const dataToSend = {
            chatroom_id: chatID,
            user_uuid: usernameToSend
        };

        console.log(`/exitroom sending user_uuid: ${dataToSend.chatroom_id}`);
        console.log(`/exitroom sending user_uuid: ${dataToSend.user_uuid}`);

        try {

            await fetch(`https://1bs9qf5xn1.execute-api.ap-southeast-1.amazonaws.com/exitroom`, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(
                dataToSend)
            });

        } catch (error) {
            console.error('Error sending message:', error);
            throw error;
        }
    }

    const dataToSendExitRoom = {
        chatroom_id: chatID,
        user_uuid: usernameToSend
    };

    const exitToServer = async () => {
        try {
            await fetch(`https://1bs9qf5xn1.execute-api.ap-southeast-1.amazonaws.com/quit`, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(
                dataToSendExitRoom)
            });
        }
        catch (error) {
            console.error('Error sending message:', error);
            throw error;
        }
    }

    const fetchMessagesFromServer = async () => {
        try {

            const dataToSend = {
                chatroom_id: chatID
            }

            const chatResponse = await fetch(`https://1bs9qf5xn1.execute-api.ap-southeast-1.amazonaws.com/poll`, {
                method: 'PATCH',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(
                    dataToSend)
                });
                
                const chatData = await chatResponse.json();
                localStorage.setItem('userDetails', JSON.stringify(chatData))
                const users = localStorage.getItem('userDetails')
                const username = users ? JSON.parse(users) : '';
                
                // CHEck if server response was null before calling state change
                if (!chatData)
                return;
            
            // Update receivedMessages state with the messages received from the server
            updateReceivedMessages(chatData);
        }
        catch (error) {
            console.error('Error fetching messages:', error);
        }
    }

    // let substringResult = "";

    // if (usernameToSend !== undefined) {
    // substringResult = usernameToSend.substring(0, 4); // Safely accessing substring if myString is defined
    // console.log(substringResult);
    // } else {
    // console.log("myString is undefined");
    // }

    const users = localStorage.getItem('user');
    const username = users ? JSON.parse(users) : '';

    // to change when server updates
    const handleBackToRoomList = () => {
        exitRoom();
        navigate('/rooms')
    }

    return (
        <main className="chat_background">
            <section className="container d-flex">
                <div className="row p-3" id="chatroom_name">
                    <div className="col-lg-1" id="back_to_roomlist" onClick={handleBackToRoomList}>
                        <h3>&lt;</h3>
                    </div>
                    <div className="col-lg-11">
                        <h3>{chatID}</h3>
                        <span>✉</span>
                    </div>
                </div>
                <div className="row p-3" id="messages_container">
                    <div className="col-lg-12">
                        {receivedMessages && receivedMessages.length > 0 && (
                            <div>
                                {receivedMessages.map((msg, index) => {
                                    return (
                                        <div key={index} className="message">
                                            <strong style={{ color: (typeof msg.username === 'string' && msg.username.substring(0, 4) === (typeof usernameToSend === 'string' && usernameToSend.substring(0, 4))) ? (usernameColors[usernameToSend] || '#000000') : '#000000' }}>{msg.username}: </strong>{msg.content}
                                        </div>
                                    );
                                })}
                            </div>
                        )}
                    </div>
                </div>
                <div className="row p-3" id="textbox_container">
                    <div className="col-lg-3" id="username">
                        <div className={`input_container link-dark text-decoration-none`}>
                            <input
                                value={customUsername}
                                readOnly
                            />
                        </div>
                    </div>
                    <div className="col-lg-9 input_container">
                        <form onSubmit={send}>
                            <input
                                placeholder="Say something..."
                                className={`${isGlittering ? 'textbox_glitter' : ''}`}
                                value={message}
                                onChange={handleInputChange}
                                disabled={isVisible}
                            />
                            <span>↵</span>
                        </form>
                        {isVisible &&
                            <strong className='warning'>⚠ WARNING: You are spamming!</strong>
                        }
                    </div>
                </div>
            </section>
            <div className="animate_enter_room_background"></div>
        </main>
    );
}; 

export default ChatPage;