import React from 'react';
import '@fortawesome/fontawesome-free/css/all.css';

export interface ChatProps { }

const ChatPage: React.FunctionComponent<ChatProps> = (props) => {

    return (
        <div className="container mt-4">
            <div className="card mx-auto" style={{ maxWidth:"400px"}}>
                <div className="card-header bg-transparent">
                    <div className="navbar navbar-expand p-0">
                        <ul className="navbar-nav ms-auto">
                            <li className="nav-item">
                                <a href="#!" className="nav-link">
                                    <i className="fas fa-video"></i>
                                </a>
                            </li>
                            <li className="nav-item">
                                <a href="#!" className="nav-link">
                                    <i className="fas fa-times"></i>
                                </a>
                            </li>
                        </ul>
                    </div>
                </div>
            </div>
            <p>This is the chat page.</p>
        </div>
    )
}; 

export default ChatPage;