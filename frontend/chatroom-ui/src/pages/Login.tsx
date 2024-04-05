import React from 'react';
import LoginButton from '.././components/LoginButton';
import '../styles/login.css'
import chatroomLogo from '../assets/chatroomLogo.gif';

export interface LoginProps { }

const LoginPage: React.FunctionComponent<LoginProps> = () => {

    // make state to store the user
    //const [userID, setuserID] = useState("");

    //<img src={chatroomLogo} alt = 'Logo' className = 'ChatroomGif'/>

    localStorage.clear();

    return (
        <main className="login_background">
            <section className='login_page'>
                <div className='login_container'>
                    <p className='logo'>𝕱𝖎𝖑𝖙𝖊𝖗𝕹𝖔𝖙𝕱𝖔𝖚𝖓𝖉</p>
                    <LoginButton> Enter </LoginButton>
                </div>
            </section>
        </main>
    );
};

export default LoginPage;