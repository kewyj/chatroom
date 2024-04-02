import React from 'react';
import LoginButton from '.././components/LoginButton';
import '../styles/login.css'
import chatroomLogo from '../assets/chatroomLogo.gif';

export interface LoginProps { }

const LoginPage: React.FunctionComponent<LoginProps> = () => {

    // make state to store the user
    //const [userID, setuserID] = useState("");

    return (
        <main className="login_background">
            <section className='LoginPage'>
                <div className='loginContainer'>
                    <img src={chatroomLogo} alt = 'Logo' className = 'ChatroomGif'/>
                    <LoginButton> Enter </LoginButton>
                </div>
            </section>
        </main>
    );
};

export default LoginPage;