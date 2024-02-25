import React from 'react';
import LoginButton from '.././components/LoginButton';
import './styles.css'
import chatroomLogo from '../assets/chatroomLogo.gif';

export interface LoginProps { }

const LoginPage: React.FunctionComponent<LoginProps> = () => {

    // make state to store the user
    //const [userID, setuserID] = useState("");

    return (
        <div className='LoginPage'>
            <div className='loginContainer'>
                <img src={chatroomLogo} alt = 'Logo' className = 'ChatroomGif'/>
                <LoginButton>Log In</LoginButton>
            </div>
        </div>
    );
};

export default LoginPage;