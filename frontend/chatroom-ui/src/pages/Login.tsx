import React from 'react';
import LoginButton from '.././components/LoginButton';
import '../styles/login.css'
import chatroomLogo from '../assets/chatroomLogo.gif';
import { useSelector, useDispatch } from 'react-redux';

export interface LoginProps { }

interface AppState {
    userID: string;
}

const LoginPage: React.FunctionComponent<LoginProps> = () => {

    // make state to store the user
    //const [userID, setuserID] = useState("");

    const userID = useSelector((state: AppState) => state.userID);

    localStorage.removeItem('user_uuid');

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