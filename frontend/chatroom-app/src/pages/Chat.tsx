import React, {useEffect, useState} from 'react';
import "@fortawesome/fontawesome-free/css/all.css";

export interface ChatProps { }

const ChatPage: React.FunctionComponent<ChatProps> = (props) => {

    const [username, setUsername] = useState('username');
    const [messages, setMessages] = useState([]);
    const [message, setMessage] = useState("");

    useEffect(() => {

    }, []);

    const send = event => {
        // prevents page refreshing
        event.preventDefault();
    }

    return (
        <div className="container">
            <div className="d-flex flex-column align-items-stretch flex-shrink-8 bg-white">
                <div
                    className="d-flex align-items-center flex-shrink-8 p-3 link-dark text-decoration-none border-bottom">
                    <input className="fs-5 fw-semibold" value={username} onChange={event => setUsername(event.target.value)} />
                </div>
                <div className="list-group list-geoup-flush border-bottom scrollarea">
                    {messages.map(message => {
                        return (
                        <div className="list-group-item list-group-item-acion py-3 lh-tight">
                            <div className="d-flex w-100 align-items-center justify-content-between">
                                <strong className="mb-1">{message.username}</strong>
                            </div>
                            <div className="col-10 mb-1 small">{message.message}</div>
                        </div>
                        )
                    })}
                </div>
            </div>
            <form onSubmit={event => send(event)}>
                <input className="form-control" placeholder="Say something..." value={ message } onChange={event => setMessage(event.target.value)} />
            </form>
        </div>
    )
}; 

export default ChatPage;