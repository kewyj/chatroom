import React, { useEffect } from 'react';
import LoginButton from '.././components/LoginButton';

export interface LoginProps { }

const LoginPage: React.FunctionComponent<LoginProps> = (props) => {

    // make state to store the user
    //const [userID, setuserID] = useState("");

    return (
        <div className="LoginPage">
            <div className="form-container">
                <form>
                    <div className="form-group">
                        <label htmlFor="userName" className="form-label">Username:</label>
                        <input type="text" className="form-control" id="userName" />
                    </div>
                </form>
            </div>
            <div>
                <LoginButton>Log In</LoginButton>
            </div>
        </div>
    );
};

export default LoginPage;