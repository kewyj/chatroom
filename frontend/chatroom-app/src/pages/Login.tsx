import React, { useEffect } from 'react';
import LoginButton from '.././components/LoginButton';

export interface LoginProps { }

const LoginPage: React.FunctionComponent<LoginProps> = (props) => {

    // make state to store the user
    //const [userID, setuserID] = useState("");

    return (
        <div className="LoginPage">
            <div>
                <LoginButton>Log In</LoginButton>
            </div>
        </div>
    );
};

export default LoginPage;